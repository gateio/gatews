package gatews

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"net"
	"time"
)

type SubscribeOptions struct {
	ID int64 `json:"id"`
}

func (ws *WsService) Subscribe(channel string, payload []string) error {
	if (ws.conf.Key == "" || ws.conf.Secret == "") && authChannel[channel] {
		return newAuthEmptyErr()
	}

	msgCh, ok := ws.msgChs.Load(channel)
	if !ok {
		msgCh = make(chan *UpdateMsg, 1)
		go ws.receiveCallMsg(channel, msgCh.(chan *UpdateMsg))
	}

	err := ws.newBaseChannel(channel, payload, msgCh.(chan *UpdateMsg), nil)
	if err != nil {
		return err
	}

	return nil
}

func (ws *WsService) SubscribeWithOption(channel string, payload []string, op *SubscribeOptions) error {
	if (ws.conf.Key == "" || ws.conf.Secret == "") && authChannel[channel] {
		return newAuthEmptyErr()
	}

	msgCh, ok := ws.msgChs.Load(channel)
	if !ok {
		msgCh = make(chan *UpdateMsg, 1)
		go ws.receiveCallMsg(channel, msgCh.(chan *UpdateMsg))
	}

	err := ws.newBaseChannel(channel, payload, msgCh.(chan *UpdateMsg), op)
	if err != nil {
		return err
	}

	return nil
}

func (ws *WsService) UnSubscribe(channel string, payload []string) error {
	return ws.baseSubscribe(UnSubscribe, channel, payload, nil)
}

func (ws *WsService) newBaseChannel(channel string, payload []string, bch chan *UpdateMsg, op *SubscribeOptions) error {
	err := ws.baseSubscribe(Subscribe, channel, payload, op)
	if err != nil {
		return err
	}

	if _, ok := ws.msgChs.Load(channel); !ok {
		ws.msgChs.Store(channel, bch)
	}

	ws.readMsg()

	return nil
}

func (ws *WsService) baseSubscribe(event string, channel string, payload []string, op *SubscribeOptions) error {
	ts := time.Now().Unix()
	hash := hmac.New(sha512.New, []byte(ws.conf.Secret))
	hash.Write([]byte(fmt.Sprintf("channel=%s&event=%s&time=%d", channel, Subscribe, ts)))
	req := Request{
		Time:    ts,
		Channel: channel,
		Event:   event,
		Payload: payload,
		Auth: Auth{
			Method: AuthMethodApiKey,
			Key:    ws.conf.Key,
			Secret: hex.EncodeToString(hash.Sum(nil)),
		},
	}
	// options
	if op != nil {
		req.Id = &op.ID
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

	if v, ok := ws.conf.subscribeMsg.Load(channel); ok {
		reqs := v.([]requestHistory)
		reqs = append(reqs, requestHistory{
			Channel: channel,
			Event:   event,
			Payload: payload,
		})
		ws.conf.subscribeMsg.Store(channel, reqs)
	} else {
		ws.conf.subscribeMsg.Store(channel, []requestHistory{{
			Channel: channel,
			Event:   event,
			Payload: payload,
		}})
	}

	return nil
}

// readMsg only run once to read message
func (ws *WsService) readMsg() {
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
						ne, hasNetErr := err.(net.Error)
						noe, hasNetOpErr := err.(*net.OpError)
						if websocket.IsUnexpectedCloseError(err) || (hasNetErr && ne.Timeout()) || (hasNetOpErr && noe != nil) ||
							websocket.IsCloseError(err) {
							ws.Logger.Printf("websocket err:%s", err.Error())
							if e := ws.reconnect(); e != nil {
								ws.Logger.Printf("reconnect err:%s", err.Error())
								return
							} else {
								ws.Logger.Printf("reconnect success, continue read message")
								continue
							}
						} else {
							ws.Logger.Printf("wsRead err:%s, type:%T", err.Error(), err)
							return
						}
					}
					var rawTrade UpdateMsg
					if err := json.Unmarshal(message, &rawTrade); err != nil {
						ws.Logger.Printf("Unmarshal err:%s, body:%s", err.Error(), string(message))
						continue
					}

					if bch, ok := ws.msgChs.Load(rawTrade.Channel); ok {
						bch.(chan *UpdateMsg) <- &rawTrade
					}
				}
			}
		}()
	})
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

func (ws *WsService) receiveCallMsg(channel string, msgCh chan *UpdateMsg) {
	defer close(msgCh)
	for {
		select {
		case <-ws.Ctx.Done():
			ws.Logger.Printf("received parent context exit")
			return
		case msg := <-msgCh:
			if call, ok := ws.calls.Load(channel); ok {
				call.(callBack)(msg)
			}
		}
	}
}
