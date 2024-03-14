package resp

type SpotOrder struct {
	Left         string `json:"left,omitempty"`
	UpdateTime   string `json:"update_time,omitempty"`
	Amount       string `json:"amount"`
	CreateTime   string `json:"create_time,omitempty"`
	Price        string `json:"price,omitempty"`
	FinishAs     string `json:"finish_as,omitempty"`
	StpAct       string `json:"stp_act,omitempty"`
	TimeInForce  string `json:"time_in_force,omitempty"`
	CurrencyPair string `json:"currency_pair"`
	Type         string `json:"type,omitempty"`
	Account      string `json:"account,omitempty"`
	Side         string `json:"side"`
	AmendText    string `json:"amend_text,omitempty"`
	Text         string `json:"text,omitempty"`
	Status       string `json:"status,omitempty"`
	Iceberg      string `json:"iceberg,omitempty"`
	AvgDealPrice string `json:"avg_deal_price,omitempty"`
	FilledTotal  string `json:"filled_total,omitempty"`
	Id           string `json:"id,omitempty"`
	FillPrice    string `json:"fill_price,omitempty"`
	UpdateTimeMs int64  `json:"update_time_ms,omitempty"`
	CreateTimeMs int64  `json:"create_time_ms,omitempty"`
	StpId        int32  `json:"stp_id,omitempty"`
	AutoRepay    bool   `json:"auto_repay,omitempty"`
	AutoBorrow   bool   `json:"auto_borrow,omitempty"`
	Succeeded    bool   `json:"succeeded"`
}
