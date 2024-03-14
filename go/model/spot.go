package model

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
