package model

import (
	"encoding/json"
)

type AmendOrderParam struct {
	Amount       string `json:"amount,omitempty"`     // New order amount. `amount` and `price` must specify one of them
	Price        string `json:"price,omitempty"`      // New order price. `amount` and `Price` must specify one of them"
	AmendText    string `json:"amend_text,omitempty"` // Custom info during amending order
	OrderId      string `json:"order_id,omitempty" `
	CurrencyPair string `json:"currency_pair,omitempty" `
	Account      string `json:"account,omitempty"`
}

type CancelOrderParam struct {
	OrderId      string `json:"order_id,omitempty"`
	CurrencyPair string `json:"currency_pair,omitempty"`
	Account      string `json:"account,omitempty"`
}

type CancelOrderWithCpParam struct {
	CurrencyPair string `json:"currency_pair,omitempty"`
	Side         string `json:"side,omitempty"`
	Account      string `json:"account,omitempty"`
}

type StatusOrderParam struct {
	OrderId      string `json:"order_id,omitempty" `
	CurrencyPair string `json:"currency_pair,omitempty" `
	Account      string `json:"account,omitempty"`
}

type ApiRequestSummary struct {
	Event   string          `json:"event"`
	Channel string          `json:"channel"`
	Payload json.RawMessage `json:"payload"`
	Time    int             `json:"time,omitempty"`
}

type ApiRequestPayload struct {
	RequestId      string            `json:"req_id"`
	ApiKey         string            `json:"api_key"`
	Timestamp      string            `json:"timestamp"`
	Signature      string            `json:"signature"`
	Channel        string            `json:"-"`
	Event          string            `json:"-"`
	TraceID        string            `json:"trace_id"`
	ClientID       string            `json:"-"`
	RequestHeaders map[string]string `json:"req_header"`
	RequestParam   json.RawMessage   `json:"req_param"`
	CurrencyPairs  []string          `json:"-"`
}

type AmendOrderRequest struct {
	Channel string `json:"channel"`
	Event   string `json:"event"`
	Payload struct {
		RequestId      string            `json:"request_id"`
		ApiKey         string            `json:"api_key"`
		Timestamp      string            `json:"timestamp"`
		Signature      string            `json:"signature"`
		RequestHeaders map[string]string `json:"request_headers"`
		RequestParam   json.RawMessage   `json:"request_param"`
	} `json:"payload"`
	Time int `json:"time"`
}

type CancelOrderRequest struct {
	Channel string `json:"channel"`
	Event   string `json:"event"`
	Payload struct {
		RequestId      string            `json:"request_id"`
		RequestHeaders map[string]string `json:"request_headers"`
		RequestParam   CancelOrderParam  `json:"request_param"`
		ApiKey         string            `json:"api_key"`
		Timestamp      string            `json:"timestamp"`
		Signature      string            `json:"signature"`
	} `json:"payload"`
	Time int `json:"time"`
}

type CancelOrderWithCpRequest struct {
	Channel string `json:"channel"`
	Event   string `json:"event"`
	Payload struct {
		RequestId      string                 `json:"request_id"`
		RequestHeaders map[string]string      `json:"request_headers"`
		RequestParam   CancelOrderWithCpParam `json:"request_param"`
		ApiKey         string                 `json:"api_key"`
		Timestamp      string                 `json:"timestamp"`
		Signature      string                 `json:"signature"`
	} `json:"payload"`
	Time int `json:"time"`
}

type StatusOrderRequest struct {
	Channel string `json:"channel"`
	Event   string `json:"event"`
	Payload struct {
		RequestId      string            `json:"request_id"`
		RequestHeaders map[string]string `json:"request_headers"`
		RequestParam   StatusOrderParam  `json:"request_param"`
		ApiKey         string            `json:"api_key"`
		Timestamp      string            `json:"timestamp"`
		Signature      string            `json:"signature"`
	} `json:"payload"`
	Time int `json:"time"`
}

type SpotRestApiRequestSummary struct {
	Request  *ApiRequestPayload
	AuthInfo *AuthResp
}

type AuthResp struct {
	ModelApiKeyInfo *ModeAPIKey
	ApiKey          string
	UserId          int64
	ApiKeyType      int
	Mode            int
	StpGroupId      int64
	MarketWhitelist bool
}

type ModeAPIKey struct {
	// CreatedAt      time.Time             `json:"created_at"`
	// UpdatedAt      time.Time             `json:"updated_at"`
	// Perms          map[string]int        `json:"perms"`
	// DeleteMailType *constant.MailType    `json:"delete_mail_type"`
	// FreezeMailType *constant.MailType    `json:"freeze_mail_type"`
	// DeleteOccurAt  *time.Time            `json:"delete_occur_at"`
	// FreezeOccurAt  *time.Time            `json:"freeze_occur_at"`
	// BrokerID       *int64                `json:"broker_id"`
	// LastAccess     *time.Time            `json:"last_access"`
	// APIKey         string                `json:"api_key"`
	// CurrencyPairs  string                `json:"pairs"`
	// IPWhitelist    string                `json:"ip_whitelist"`
	// APISecret      string                `json:"secret"`
	// Name           string                `json:"name"`
	// StpGroupId     int64                 `json:"stp_group_id"`
	// Type           constant.KeyType      `json:"type"`
	// Mode           constant.ModeType     `json:"mode"`
	// ID             int64                 `json:"id"`
	// UserID         int64                 `json:"user_id"`
	// State          constant.StateType    `json:"state"`
	// Source         constant.APIKeySource `json:"source"`
}
