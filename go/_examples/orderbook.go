package main

import (
	"encoding/json"
	"fmt"
	"github.com/gansidui/skiplist"
	gate "github.com/gateio/gatews/go"
	"github.com/shopspring/decimal"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

const (
	MaxLimit = 100
	depthUrl = "https://api.gateio.ws/api/v4/spot/order_book?currency_pair=%s&limit=%d&with_id=true"
)

var (
	localOrderBook = sync.Map{}
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

func OrderBookExample(ws *gate.WsService) {
	callBack := gate.NewCallBack(func(msg *gate.UpdateMsg) {
		var depthMsg gate.SpotUpdateDepthMsg
		if err := json.Unmarshal(msg.Result, &depthMsg); err != nil {
			log.Printf("order Unmarshal err:%s", err.Error())
		}
		log.Printf("f:%d, l:%d", depthMsg.FirstId, depthMsg.LastId)
		if err := updateLocalOrderBook(depthMsg); err != nil {
			log.Printf("err:%s", err.Error())
		} else {
			localOrderBook.Range(func(key, value interface{}) bool {
				for e := value.(*OrderBook).Asks.Front(); e != nil; e = e.Next() {
					fmt.Println(e.Value.(*OrderBookEntry).Price, "-->", e.Value.(*OrderBookEntry).Size)
				}
				fmt.Println("<><><><><><><><>")
				for e := value.(*OrderBook).Bids.Front(); e != nil; e = e.Next() {
					fmt.Println(e.Value.(*OrderBookEntry).Price, "-->", e.Value.(*OrderBookEntry).Size)
				}
				fmt.Printf("%s >>> %+v\n", key.(string), value.(*OrderBook))
				return true
			})
		}
	})
	ws.SetCallBack(gate.ChannelOrderBookUpdate, callBack)
	if err := ws.Subscribe(gate.ChannelOrderBookUpdate, []string{"BCH_USDT", "1000ms"}); err != nil {
		log.Printf("Subscribe err:%s", err.Error())
	}
	if err := ws.Subscribe(gate.ChannelOrderBookUpdate, []string{"BTC_USDT", "1000ms"}); err != nil {
		log.Printf("Subscribe err:%s", err.Error())
	}
	if err := ws.Subscribe(gate.ChannelOrderBookUpdate, []string{"ETH_USDT", "1000ms"}); err != nil {
		log.Printf("Subscribe err:%s", err.Error())
	}
}

func updateLocalOrderBook(msg gate.SpotUpdateDepthMsg) error {
	if orderbook, ok := localOrderBook.Load(msg.CurrencyPair); ok {
		log.Printf("----------now id %d", orderbook.(*OrderBook).ID)

		if orderbook.(*OrderBook).ID+1 >= msg.FirstId && orderbook.(*OrderBook).ID+1 <= msg.LastId {
			if err := updateOrderBook(orderbook.(*OrderBook), msg); err != nil {
				return err
			}
		} else if msg.LastId < orderbook.(*OrderBook).ID+1 {
			return nil
		} else if orderbook.(*OrderBook).ID+1 < msg.FirstId {
			localOrderBook.Delete(msg.CurrencyPair)
			log.Printf(">>>>>>>>>>>>>>>>%s depth is fall behind, f:%d, l:%d", msg.CurrencyPair, msg.FirstId, msg.LastId)
			return nil
		}
	} else {
		log.Printf("init %s depth", msg.CurrencyPair)

		depth, err := getBaseDepth(msg.CurrencyPair, MaxLimit)
		if err != nil {
			return err
		}
		localOrderBook.Store(msg.CurrencyPair, depth)
	}
	return nil
}

func getBaseDepth(cp string, limit int) (*OrderBook, error) {
	c := http.DefaultClient
	url := fmt.Sprintf(depthUrl, cp, limit)
	resp, err := c.Get(url)
	if err != nil {
		return nil, err
	}
	body, _ := ioutil.ReadAll(resp.Body)

	var baseOrderBook HttpOrderBook
	err = json.Unmarshal(body, &baseOrderBook)
	if err != nil {
		return nil, err
	}

	asks := skiplist.New()
	bids := skiplist.New()
	if len(baseOrderBook.Asks) > 0 {
		for _, a := range baseOrderBook.Asks {
			askEntry, err := fromHttpOrderBook(a)
			if err != nil {
				return nil, err
			}
			asks.Insert(askEntry)
		}
		for _, b := range baseOrderBook.Bids {
			bidEntry, err := fromHttpOrderBook(b)
			if err != nil {
				return nil, err
			}
			bids.Insert(bidEntry)
		}
	}
	if asks.Len() > 0 && bids.Len() > 0 {
		// reject overlapping
		if asks.Front().Value.(*OrderBookEntry).Price.LessThanOrEqual(bids.Back().Value.(*OrderBookEntry).Price) {
			return nil, fmt.Errorf("overlapping ask and bid price")
		}
	}

	return &OrderBook{
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

func updateOrderBook(orderbook *OrderBook, update gate.SpotUpdateDepthMsg) error {
	orderbook.ID = update.LastId
	if len(update.Ask) > 0 {
		for _, ask := range update.Ask {
			askEntry, err := fromHttpOrderBook(ask)
			if err != nil {
				log.Printf("incorrect http ask entry %v: %v", ask, err)
				return err
			}
			if ask[1] == "0" {
				orderbook.Asks.Delete(askEntry)
			} else {
				if ele := orderbook.Asks.Find(askEntry); ele != nil {
					ele.Value.(*OrderBookEntry).Size = ask[1]
				} else {
					orderbook.Asks.Insert(askEntry)
				}
			}
		}
	} else if len(update.Bid) > 0 {
		for _, bid := range update.Bid {
			bidEntry, err := fromHttpOrderBook(bid)
			if err != nil {
				log.Printf("incorrect http bid entry %v: %v", bid, err)
				return err
			}
			if bid[1] == "0" {
				orderbook.Bids.Delete(bidEntry)
			} else {
				if ele := orderbook.Bids.Find(bidEntry); ele != nil {
					ele.Value.(*OrderBookEntry).Size = bid[1]
				} else {
					orderbook.Bids.Insert(bidEntry)
				}
			}
		}
	}
	return nil
}
