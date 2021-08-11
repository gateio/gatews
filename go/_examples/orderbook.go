package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/gansidui/skiplist"
	"github.com/shopspring/decimal"
	queue2 "github.com/yireyun/go-queue"

	gate "github.com/gateio/gatews/go"
)

const (
	MaxLimit  = 100
	QueueSize = 3000
	depthUrl  = "https://api.gateio.ws/api/v4/spot/order_book?currency_pair=%s&limit=%d&with_id=true"
)

var (
	localOrderBook = sync.Map{}
	queue          = queue2.NewQueue(QueueSize)
	spotMsg        = new(sync.Map)
)

type OrderBookEntry struct {
	Price decimal.Decimal `json:"p"`
	Size  string          `json:"s"`
}

func (e *OrderBookEntry) Less(other interface{}) bool {
	return e.Price.LessThan(other.(*OrderBookEntry).Price)
}

type OrderBook struct {
	ID   int64
	Asks *skiplist.SkipList
	Bids *skiplist.SkipList
}

type HttpOrderBook struct {
	ID   int64      `json:"id"`
	Asks [][]string `json:"asks"`
	Bids [][]string `json:"bids"`
}

func LocalOrderBook(ctx context.Context, ws *gate.WsService, cps []string) {
	callBack := gate.NewCallBack(func(msg *gate.UpdateMsg) {
		queue.Put(msg)
		if queue.Quantity()+5 >= QueueSize {
			log.Printf("LocalOrderBook queue is almost full")
		}
	})

	channel := gate.ChannelSpotOrderBookUpdate

	ws.SetCallBack(channel, callBack)

	for _, cp := range cps {
		if err := ws.Subscribe(channel, []string{cp, "100ms"}); err != nil {
			log.Printf("Subscribe err:%s", err.Error())
		}
		if resp, depth, err := getBaseDepth(cp, MaxLimit); err == nil {
			localOrderBook.Store(cp, depth)
			log.Printf("spot init market %s order book asks:%v, bids:%v", cp, resp.Asks, resp.Bids)
		}
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
			msg, ok, _ := queue.Get()
			if !ok {
				continue
			}
			var depthMsg gate.SpotUpdateDepthMsg
			if err := json.Unmarshal(msg.(*gate.UpdateMsg).Result, &depthMsg); err != nil {
				log.Printf("order Unmarshal err:%s", err.Error())
			}
			if v, ok := spotMsg.Load(depthMsg.CurrencyPair); ok {
				if v.(gate.SpotUpdateDepthMsg).LastId+1 != depthMsg.FirstId {
					log.Printf("%s order book msg id discontinuous, old id %d-%d, new msg id %d-%d", depthMsg.CurrencyPair,
						v.(gate.SpotUpdateDepthMsg).FirstId, v.(gate.SpotUpdateDepthMsg).LastId, depthMsg.FirstId, depthMsg.LastId)
				}
			}
			spotMsg.Store(depthMsg.CurrencyPair, depthMsg)
			if err := updateLocalOrderBook(depthMsg); err != nil {
				log.Printf("err:%s", err.Error())
			}
		}
	}
}

func updateLocalOrderBook(msg gate.SpotUpdateDepthMsg) error {
	// log.Printf("updateLocalOrderBook msg:%+v", msg)

	if orderBook, ok := localOrderBook.Load(msg.CurrencyPair); ok {
		if orderBook.(*OrderBook).ID+1 >= msg.FirstId && orderBook.(*OrderBook).ID+1 <= msg.LastId {
			if err := updateOrderBook(orderBook.(*OrderBook), msg); err != nil {
				log.Printf("spot:%s", err.Error())
				if strings.Contains(err.Error(), "overlapping") {
					if err = reGetBaseDepth(0, msg); err != nil {
						return err
					}
				}
				return err
			}
		} else if msg.LastId < orderBook.(*OrderBook).ID+1 {
			return nil
		} else if orderBook.(*OrderBook).ID+1 < msg.FirstId {
			log.Printf(">>>>>>>>>>>>>>>>%s depth is fall behind, now:%d, f:%d, l:%d", msg.CurrencyPair, orderBook.(*OrderBook).ID, msg.FirstId, msg.LastId)
			log.Printf("reinit %s depth", msg.CurrencyPair)
			if err := reGetBaseDepth(orderBook.(*OrderBook).ID, msg); err != nil {
				return err
			} else {
				if orderBook, ok := localOrderBook.Load(msg.CurrencyPair); ok {
					if orderBook.(*OrderBook).ID+1 >= msg.FirstId && orderBook.(*OrderBook).ID+1 <= msg.LastId {
						if err := updateOrderBook(orderBook.(*OrderBook), msg); err != nil {
							log.Printf("after reGetBaseDepth reconsume msg failed, msg:%+v, err:%s", msg, err.Error())
							return err
						}
					}
				}
			}
		}
	} else if msg.CurrencyPair != "" {
		log.Printf("init %s depth", msg.CurrencyPair)
		if err := reGetBaseDepth(0, msg); err != nil {
			return err
		}
	}
	return nil
}

func reGetBaseDepth(nowID int64, msg gate.SpotUpdateDepthMsg) (err error) {
	var depth *OrderBook
	var resp HttpOrderBook
	for nowID < msg.FirstId {
		resp, depth, err = getBaseDepth(msg.CurrencyPair, MaxLimit)
		if err != nil {
			return err
		}
		nowID = depth.ID
	}

	if depth != nil && depth.ID > 0 {
		log.Printf("after reGetBaseDepth resp %+v, ask min %s, bid max %s", resp,
			depth.Asks.Front().Value.(*OrderBookEntry).Price.String(), depth.Bids.Back().Value.(*OrderBookEntry).Price.String())
		localOrderBook.Store(msg.CurrencyPair, depth)
	}

	return nil
}

func getBaseDepth(cp string, limit int) (HttpOrderBook, *OrderBook, error) {
	url := fmt.Sprintf(depthUrl, cp, limit)
	resp, err := http.DefaultClient.Get(url)
	if err != nil {
		return HttpOrderBook{}, nil, err
	}
	body, _ := ioutil.ReadAll(resp.Body)

	var baseOrderBook HttpOrderBook
	err = json.Unmarshal(body, &baseOrderBook)
	if err != nil {
		return baseOrderBook, nil, err
	}

	asks := skiplist.New()
	bids := skiplist.New()
	if len(baseOrderBook.Asks) > 0 {
		for _, a := range baseOrderBook.Asks {
			askEntry, err := fromHttpOrderBook(a)
			if err != nil {
				return baseOrderBook, nil, err
			}
			asks.Insert(askEntry)
		}
		for _, b := range baseOrderBook.Bids {
			bidEntry, err := fromHttpOrderBook(b)
			if err != nil {
				return baseOrderBook, nil, err
			}
			bids.Insert(bidEntry)
		}
	}
	if asks.Len() > 0 && bids.Len() > 0 {
		// reject overlapping
		if asks.Front().Value.(*OrderBookEntry).Price.LessThanOrEqual(bids.Back().Value.(*OrderBookEntry).Price) {
			return baseOrderBook, nil, fmt.Errorf("overlapping price ask[%s] and bid[%s]",
				asks.Front().Value.(*OrderBookEntry).Price.String(), bids.Back().Value.(*OrderBookEntry).Price.String())
		}
	}

	return baseOrderBook, &OrderBook{
		ID:   baseOrderBook.ID,
		Asks: asks,
		Bids: bids,
	}, nil
}

func fromHttpOrderBook(apiEntry []string) (*OrderBookEntry, error) {
	if len(apiEntry) != 2 {
		return nil, fmt.Errorf("invalid http order book entry")
	}
	price, err := decimal.NewFromString(apiEntry[0])
	if err != nil {
		return nil, fmt.Errorf("invalid price %s: %v", apiEntry[0], err)
	}
	return &OrderBookEntry{Price: price, Size: apiEntry[1]}, nil
}

func updateOrderBook(orderBook *OrderBook, update gate.SpotUpdateDepthMsg) error {
	orderBook.ID = update.LastId
	if len(update.Ask) > 0 {
		for _, ask := range update.Ask {
			askEntry, err := fromHttpOrderBook(ask)
			if err != nil {
				log.Printf("incorrect http ask entry %v: %v", ask, err)
				return err
			}
			if ask[1] == "0" {
				for e := orderBook.Asks.Front(); e != nil; e = e.Next() {
					if e.Value.(*OrderBookEntry).Price.String() == ask[0] {
						orderBook.Asks.Delete(e.Value)
						break
					}
				}
			} else {
				if e := orderBook.Asks.Find(askEntry); e != nil {
					e.Value.(*OrderBookEntry).Size = ask[1]
				} else {
					orderBook.Asks.Insert(askEntry)
				}
			}
		}
	}

	if len(update.Bid) > 0 {
		for _, bid := range update.Bid {
			bidEntry, err := fromHttpOrderBook(bid)
			if err != nil {
				log.Printf("incorrect http bid entry %v: %v", bid, err)
				return err
			}
			if bid[1] == "0" {
				for e := orderBook.Bids.Back(); e != nil; e = e.Prev() {
					if e.Value.(*OrderBookEntry).Price.String() == bid[0] {
						orderBook.Bids.Delete(e.Value)
						break
					}
				}
			} else {
				if e := orderBook.Bids.Find(bidEntry); e != nil {
					e.Value.(*OrderBookEntry).Size = bid[1]
				} else {
					orderBook.Bids.Insert(bidEntry)
				}
			}
		}
	}

	// judge overlapping
	if orderBook.Asks.Len() > 0 && orderBook.Bids.Len() > 0 {
		// reject overlapping
		if orderBook.Asks.Front().Value.(*OrderBookEntry).Price.LessThanOrEqual(orderBook.Bids.Back().Value.(*OrderBookEntry).Price) {
			return fmt.Errorf("overlapping price ask[%s] and bid[%s]",
				orderBook.Asks.Front().Value.(*OrderBookEntry).Price.String(), orderBook.Bids.Back().Value.(*OrderBookEntry).Price.String())
		}
	}
	return nil
}
