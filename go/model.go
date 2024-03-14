package gatews

import (
	"encoding/json"
	"fmt"
)

type UpdateMsg struct {
	Header  ResponseHeader  `json:"header"`
	Time    int64           `json:"time"`
	TimeMs  int64           `json:"time_ms"`
	Id      *int64          `json:"id,omitempty"`
	Channel string          `json:"channel"`
	Event   string          `json:"event"`
	Error   *ServiceError   `json:"error,omitempty"`
	Result  json.RawMessage `json:"result"`
	Data    struct {
		Result json.RawMessage `json:"result"`
		Errs   *struct {
			Label   string `json:"label"`
			Message string `json:"message"`
		} `json:"errs"`
	} `json:"data"`
}

type ResponseHeader struct {
	ResponseTime string `json:"response_time"`
	Status       string `json:"status"`
	Channel      string `json:"channel"`
	Event        string `json:"event"`
	ClientID     string `json:"client_id"`
}

func (u *UpdateMsg) GetChannel() string {
	if u.Channel != "" {
		return u.Channel
	}

	return u.Header.Channel
}

type ServiceError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e ServiceError) Error() string {
	return e.Message
}

func newAuthEmptyErr() error {
	return fmt.Errorf("auth key or secret empty")
}

type WSEvent struct {
	UpdateMsg
}

type ChannelEvent struct {
	Event  string
	Market []string
}

type WebsocketRequest struct {
	Market []string
}

type Request struct {
	App     string `json:"app,omitempty"`
	Time    int64  `json:"time"`
	Id      *int64 `json:"id,omitempty"`
	Channel string `json:"channel"`
	Event   string `json:"event"`
	Auth    Auth   `json:"auth"`
	Payload any    `json:"payload"`
}

type Auth struct {
	Method string `json:"method"`
	Key    string `json:"KEY"`
	Secret string `json:"SIGN"`
}

type requestHistory struct {
	Channel string `json:"channel"`
	Event   string `json:"event"`
	Payload any    `json:"payload"`
	op      *SubscribeOptions
}

type APIReq struct {
	ApiKey    string          `json:"api_key"`
	Signature string          `json:"signature"`
	Timestamp string          `json:"timestamp"`
	ReqId     string          `json:"req_id"`
	ReqHeader json.RawMessage `json:"req_header"`
	ReqParam  json.RawMessage `json:"req_param"`
}

type APIResp struct {
	ClientID   string `json:"client_id"`
	ReqID      string `json:"req_id"`
	RespTimeMs int64  `json:"resp_time_ms"`
	Status     int    `json:"status"`
	ReqHeader  struct {
		XGateChannelID string `json:"x-gate-channel-id"`
	} `json:"req_header"`
	Data struct {
		Error  any `json:"error"`
		Result any `json:"result"`
	} `json:"data"`
}
