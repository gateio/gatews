package model

type ListFuturesOrders struct {
	Contract string `json:"contract,omitempty"`
	Status   string `json:"status,omitempty"`
	LastId   string `json:"last_id,omitempty"`
	Settle   string `json:"settle,omitempty"`
	Limit    int32  `json:"limit,omitempty"`
	Offset   int32  `json:"offset,omitempty"`
}

type CancelFuturesOrder struct {
	OrderId string `json:"order_id,omitempty"`
	Settle  string `json:"settle,omitempty"`
}

type CancelFuturesCpOrder struct {
	Contract string `json:"contract,omitempty"`
	Side     string `json:"side,omitempty"`
	Settle   string `json:"settle,omitempty"`
}

type StatusFuturesOrder struct {
	OrderId string `json:"order_id,omitempty"`
	Settle  string `json:"settle,omitempty"`
}

type AmendFuturesOrder struct {
	OrderId   string `json:"order_id,omitempty"`
	Settle    string `json:"settle,omitempty"`
	Price     string `json:"price,omitempty"`
	AmendText string `json:"amend_text"`
	Size      int64  `json:"size,omitempty"`
}

type CancelFuturesOrderIds struct {
	OrderIds []string `json:"order_ids,omitempty"`
	Settle   string   `json:"settle,omitempty"`
}
