package resp

type FutureOrder struct {
	Text         string  `json:"text,omitempty"`
	Price        string  `json:"price,omitempty"`
	BizInfo      string  `json:"biz_info,omitempty"`
	Tif          string  `json:"tif,omitempty"`
	AmendText    string  `json:"amend_text,omitempty"`
	Status       string  `json:"status,omitempty"`
	Contract     string  `json:"contract"`
	StpAct       string  `json:"stp_act,omitempty"`
	FinishAs     string  `json:"finish_as,omitempty"`
	FillPrice    string  `json:"fill_price,omitempty"`
	AutoSize     string  `json:"auto_size,omitempty"`
	Id           int64   `json:"id,omitempty"`
	CreateTime   float64 `json:"create_time,omitempty"`
	Iceberg      int64   `json:"iceberg,omitempty"`
	Size         int64   `json:"size"`
	FinishTime   float64 `json:"finish_time,omitempty"`
	Left         int64   `json:"left"`
	Refu         int32   `json:"refu,omitempty"`
	User         int32   `json:"user,omitempty"`
	StpId        int32   `json:"stp_id,omitempty"`
	IsClose      bool    `json:"is_close,omitempty"`
	Close        bool    `json:"close,omitempty"`
	IsLiq        bool    `json:"is_liq,omitempty"`
	IsReduceOnly bool    `json:"is_reduce_only,omitempty"`
	ReduceOnly   bool    `json:"reduce_only,omitempty"`
}
