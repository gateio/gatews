package gatews

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type SubscribeOptions struct {
	ID          int64 `json:"id"`
	IsReConnect bool  `json:"-"`
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

	return ws.newBaseChannel(channel, payload, msgCh.(chan *UpdateMsg), nil)
}

func (ws *WsService) SubscribeWithOption(channel string, payload any, op *SubscribeOptions) error {
	if (ws.conf.Key == "" || ws.conf.Secret == "") && authChannel[channel] {
		return newAuthEmptyErr()
	}

	msgCh, ok := ws.msgChs.Load(channel)
	if !ok {
		msgCh = make(chan *UpdateMsg, 1)
		go ws.receiveCallMsg(channel, msgCh.(chan *UpdateMsg))
	}

	return ws.newBaseChannel(channel, payload, msgCh.(chan *UpdateMsg), op)
}

func (ws *WsService) UnSubscribe(channel string, payload []string) error {
	return ws.baseSubscribe(UnSubscribe, channel, payload, nil)
}

func (ws *WsService) newBaseChannel(channel string, payload any, bch chan *UpdateMsg, op *SubscribeOptions) error {
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

func (ws *WsService) baseSubscribe(event, channel string, payload any, op *SubscribeOptions) error {
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
	ws.mu.Lock()
	defer ws.mu.Unlock()

	err = ws.Client.WriteMessage(websocket.TextMessage, byteReq)
	if err != nil {
		ws.Logger.Printf("wsWrite [%s] err:%s", channel, err.Error())
		return err
	}

	if strings.HasSuffix(channel, "ping") {
		return nil
	}

	if v, ok := ws.conf.subscribeMsg.Load(channel); ok {
		if op != nil && op.IsReConnect {
			return nil
		}
		reqs := v.([]requestHistory)
		reqs = append(reqs, requestHistory{
			Channel: channel,
			Event:   event,
			Payload: payload,
		})
		ws.conf.subscribeMsg.Store(channel, reqs)
	} else {
		// avoid saving invalid subscribe msg
		if strings.HasSuffix(channel, ".ping") || strings.HasSuffix(channel, ".time") {
			return nil
		}

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
					_, rawMsg, err := ws.Client.ReadMessage()
					if err != nil {
						ws.Logger.Printf("websocket err: %s", err.Error())
						if e := ws.reconnect(); e != nil {
							ws.Logger.Printf("reconnect err:%s", err.Error())
							return
						}
						ws.Logger.Println("reconnect success, continue read message")
						continue
					}

					var msg UpdateMsg
					if err := json.Unmarshal(rawMsg, &msg); err != nil {
						continue
					}

					channel := msg.GetChannel()
					if channel == "" {
						ws.Logger.Printf("channel is empty in message %v", msg)
						return
					}

					if bch, ok := ws.msgChs.Load(channel); ok {
						select {
						case <-ws.Ctx.Done():
							return
						default:
							if _, ok := ws.msgChs.Load(channel); ok {
								bch.(chan *UpdateMsg) <- &msg
							}
						}
					}
				}
			}
		}()
	})
}

type CallBack func(*UpdateMsg)

func NewCallBack(f func(*UpdateMsg)) func(*UpdateMsg) {
	return f
}

func (ws *WsService) SetCallBack(channel string, call CallBack) {
	if call == nil {
		return
	}
	ws.calls.Store(channel, call)
}

func (ws *WsService) receiveCallMsg(channel string, msgCh chan *UpdateMsg) {
	// avoid send closed channel error
	// defer close(msgCh)
	for {
		select {
		case <-ws.Ctx.Done():
			ws.Logger.Printf("received parent context exit")
			return
		case msg := <-msgCh:
			if call, ok := ws.calls.Load(channel); ok {
				call.(CallBack)(msg)
			}
		}
	}
}

func (ws *WsService) APIRequest(channel string, payload any, keyVals map[string]any) error {
	var err error
	ws.loginOnce.Do(func() {
		err = ws.login()
	})

	if err != nil {
		return err
	}

	if (ws.conf.Key == "" || ws.conf.Secret == "") && authChannel[channel] {
		return newAuthEmptyErr()
	}

	msgCh, ok := ws.msgChs.Load(channel)
	if !ok {
		msgCh = make(chan *UpdateMsg, 1)
		go ws.receiveCallMsg(channel, msgCh.(chan *UpdateMsg))
	}

	if _, ok := ws.msgChs.Load(channel); !ok {
		ws.msgChs.Store(channel, msgCh)
	}

	ws.readMsg()

	return ws.apiRequest(channel, payload, keyVals)
}

func (ws *WsService) login() error {
	if ws.conf.Key == "" || ws.conf.Secret == "" {
		return newAuthEmptyErr()
	}
	channel := ChannelSpotLogin
	if ws.conf.App == "futures" {
		channel = ChannelFutureLogin
	}
	msgCh, ok := ws.msgChs.Load(channel)
	if !ok {
		msgCh = make(chan *UpdateMsg, 1)
		go ws.receiveCallMsg(channel, msgCh.(chan *UpdateMsg))
	}

	if _, ok := ws.msgChs.Load(channel); !ok {
		ws.msgChs.Store(channel, msgCh)
	}

	ws.readMsg()

	return ws.apiRequest(channel, nil, nil)
}

func (ws *WsService) apiRequest(channel string, payload any, keyVals map[string]any) error {
	req := Request{
		Time:    time.Now().Unix(),
		Channel: channel,
		Event:   API,
		Payload: ws.generateAPIRequest(channel, payload, keyVals),
	}

	byteReq, err := json.Marshal(req)
	if err != nil {
		ws.Logger.Printf("req Marshal err:%s", err.Error())
		return err
	}
	ws.mu.Lock()
	defer ws.mu.Unlock()

	return ws.Client.WriteMessage(websocket.TextMessage, byteReq)
}

func (ws *WsService) generateAPIRequest(channel string, placeParam any, keyVals map[string]any) any {
	reqID := "req_id"
	gateChannelID := "T_channel_id"

	if v, ok := keyVals["req_id"]; ok {
		reqID, _ = v.(string)
	}

	if v, ok := keyVals["X-Gate-Channel-Id"]; ok {
		gateChannelID, _ = v.(string)
	}

	now := time.Now().Unix()

	reqParam, _ := json.Marshal(placeParam)

	message := fmt.Sprintf("api\n%s\n%s\n%d", channel, reqParam, now)

	return APIReq{
		ApiKey:    ws.conf.Key,
		Signature: calculateSignature(ws.conf.Secret, message),
		Timestamp: strconv.Itoa(int(now)),
		ReqId:     reqID,
		ReqHeader: json.RawMessage(fmt.Sprintf(`{"X-Gate-Channel-Id":"%s"}`, gateChannelID)),
		ReqParam:  reqParam,
	}
}

func calculateSignature(secret string, message string) string {
	h := hmac.New(sha512.New, []byte(secret))
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}
