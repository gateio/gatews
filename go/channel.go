package gatews

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	mapset "github.com/deckarep/golang-set"
	"github.com/gorilla/websocket"
	"net"
	"strings"
	"time"
)

type Auth struct {
	Method string `json:"method"`
	Key    string `json:"KEY"`
	Secret string `json:"SIGN"`
}

func (ws *WsService) newBaseChannel(payload []string, channel string, bch chan *UpdateMsg) error {
	err := ws.baseSubscribe(Subscribe, payload, channel)
	if err != nil {
		return err
	}

	ws.buChs.Store(channel, bch)

	ws.once.Do(func() {
		go func() {
			defer ws.Client.Close()

			for {
				select {
				case <-ws.Ctx.Done():
					ws.Logger.Printf("closing reader")
					return
				default:
					_, message, err := ws.Client.ReadMessage()
					if err != nil {
						ne, ok := err.(net.Error)
						noe, ok2 := err.(*net.OpError)
						if websocket.IsUnexpectedCloseError(err) || (ok && ne.Timeout()) || (ok2 && noe.Err != nil) {
							if e := ws.reconnect(); e != nil {
								ws.Logger.Printf("reconnect err:%s", err.Error())
								return
							}
						}
						ws.Logger.Printf("wsRead err:%s, type:%T", err.Error(), err)
						continue
					}
					var rawTrade UpdateMsg
					if err := json.Unmarshal(message, &rawTrade); err != nil {
						ws.Logger.Printf("Unmarshal err:%s, body:%s", err.Error(), string(message))
						continue
					}

					if bch, ok := ws.buChs.Load(rawTrade.Channel); ok {
						bch.(chan *UpdateMsg) <- &rawTrade
					}
				}
			}
		}()
	})

	return nil
}

func (ws *WsService) baseSubscribe(event string, market []string, channel string) error {
	ts := time.Now().Unix()
	hash := hmac.New(sha512.New, []byte(ws.conf.Secret))
	hash.Write([]byte(fmt.Sprintf("channel=%s&event=%s&time=%d", channel, Subscribe, ts)))
	req := Request{
		Time:    ts,
		Channel: channel,
		Event:   event,
		Payload: market,
		Auth: Auth{
			Method: AuthMethodApiKey,
			Key:    ws.conf.Key,
			Secret: hex.EncodeToString(hash.Sum(nil)),
		},
	}
	byteReq, err := json.Marshal(req)
	if err != nil {
		ws.Logger.Printf("req Marshal err:%s", err.Error())
		return err
	}

	err = ws.Client.WriteMessage(websocket.TextMessage, byteReq)
	if err != nil {
		ws.Logger.Printf("wsWrite err:%s", err.Error())
		return err
	}

	if event == Subscribe {
		if v, ok := ws.conf.markets.Load(channel); ok {
			set := mapset.NewSet()
			for _, payload := range market {
				set.Add(payload)
			}
			for _, payload := range v.([]string) {
				set.Add(payload)
			}
			market = []string{}
			for _, payload := range set.ToSlice() {
				market = append(market, payload.(string))
			}
			ws.conf.markets.Store(channel, market)
		} else {
			ws.conf.markets.Store(channel, market)
		}
	} else if event == UnSubscribe {
		if v, ok := ws.conf.markets.Load(channel); ok {
			set := mapset.NewSet()
			for _, payload := range v.([]string) {
				set.Add(payload)
			}
			for _, payload := range market {
				set.Remove(payload)
			}
			market = []string{}
			for _, payload := range set.ToSlice() {
				market = append(market, payload.(string))
			}
			ws.conf.markets.Store(channel, market)
		}
	}

	return nil
}

type callBack func(*UpdateMsg)

func NewCallBack(f func(*UpdateMsg)) func(*UpdateMsg) {
	return f
}

func (ws *WsService) SetCallBack(channel string, call callBack) {
	if call == nil {
		return
	}
	ws.calls.Store(channel, call)
}

func (ws *WsService) Subscribe(channel string, payload []string) error {
	if (ws.conf.Key == "" || ws.conf.Secret == "") && authChannel[channel] {
		return newAuthEmptyErr()
	}

	msgCh, ok := ws.buChs.Load(channel)
	if !ok {
		msgCh = make(chan *UpdateMsg)
	}

	err := ws.newBaseChannel(payload, channel, msgCh.(chan *UpdateMsg))
	if err != nil {
		return err
	}

	go func() {
		defer close(msgCh.(chan *UpdateMsg))
		for {
			select {
			case <-ws.Ctx.Done():
				ws.Logger.Printf("received parent context exit")
				return
			case msg := <-msgCh.(chan *UpdateMsg):
				if msg.Event == Subscribe && strings.Contains(string(msg.Result), "success") {
					continue
				}

				go func() {
					if call, ok := ws.calls.Load(channel); ok {
						call.(callBack)(msg)
					}
				}()
			}
		}
	}()
	return nil
}

func (ws *WsService) UnSubscribe(channel string, payload []string) error {
	return ws.baseSubscribe(UnSubscribe, payload, channel)
}
