package model

// Futures order details
type FuturesOrder struct {
	// Futures order ID
	Id int64 `json:"id,omitempty"`
	// User ID
	User int32 `json:"user,omitempty"`
	// Creation time of order
	CreateTime float64 `json:"create_time,omitempty"`
	// Order finished time. Not returned if order is open
	FinishTime float64 `json:"finish_time,omitempty"`
	// 结束方式，包括：  - filled: 完全成交 - cancelled: 用户撤销 - liquidated: 强制平仓撤销 - ioc: 未立即完全成交，因为tif设置为ioc - auto_deleveraged: 自动减仓撤销 - reduce_only: 增持仓位撤销，因为设置reduce_only或平仓 - position_closed: 因为仓位平掉了，所以挂单被撤掉 - reduce_out: 只减仓被排除的不容易成交的挂单 - stp: 订单发生自成交限制而被撤销
	FinishAs string `json:"finish_as,omitempty"`
	// Order status  - `open`: waiting to be traded - `finished`: finished
	Status string `json:"status,omitempty"`
	// Futures contract
	Contract string `json:"contract"`
	// Order size. Specify positive number to make a bid, and negative number to ask
	Size int64 `json:"size"`
	// Display size for iceberg order. 0 for non-iceberg. Note that you will have to pay the taker fee for the hidden size
	Iceberg int64 `json:"iceberg,omitempty"`
	// Order price. 0 for market order with `tif` set as `ioc`
	Price string `json:"price,omitempty"`
	// Set as `true` to close the position, with `size` set to 0
	Close bool `json:"close,omitempty"`
	// Is the order to close position
	IsClose bool `json:"is_close,omitempty"`
	// Set as `true` to be reduce-only order
	ReduceOnly bool `json:"reduce_only,omitempty"`
	// Is the order reduce-only
	IsReduceOnly bool `json:"is_reduce_only,omitempty"`
	// Is the order for liquidation
	IsLiq bool `json:"is_liq,omitempty"`
	// Time in force  - gtc: GoodTillCancelled - ioc: ImmediateOrCancelled, taker only - poc: PendingOrCancelled, makes a post-only order that always enjoys a maker fee - fok: FillOrKill, fill either completely or none
	Tif string `json:"tif,omitempty"`
	// Size left to be traded
	Left int64 `json:"left,omitempty"`
	// Fill price of the order
	FillPrice string `json:"fill_price,omitempty"`
	// User defined information. If not empty, must follow the rules below:  1. prefixed with `t-` 2. no longer than 28 bytes without `t-` prefix 3. can only include 0-9, A-Z, a-z, underscore(_), hyphen(-) or dot(.) Besides user defined information, reserved contents are listed below, denoting how the order is created:  - web: from web - api: from API - app: from mobile phones - auto_deleveraging: from ADL - liquidation: from liquidation - insurance: from insurance
	Text string `json:"text,omitempty"`
	// Taker fee
	Tkfr string `json:"tkfr,omitempty"`
	// Maker fee
	Mkfr string `json:"mkfr,omitempty"`
	// Reference user ID
	Refu int32 `json:"refu,omitempty"`
	// Set side to close dual-mode position. `close_long` closes the long side; while `close_short` the short one. Note `size` also needs to be set to 0
	AutoSize string `json:"auto_size,omitempty"`
	// 订单所属的`STP用户组`id，同一个`STP用户组`内用户之间的订单不允许发生自成交。  1. 如果撮合时两个订单的 `stp_id` 非 `0` 且相等，则不成交，而是根据 `taker` 的 `stp_act` 执行相应策略。 2. 没有设置`STP用户组`成交的订单，`stp_id` 默认返回 `0`。
	StpId int32 `json:"stp_id,omitempty"`
	// Self-Trading Prevention Action,用户可以用该字段设置自定义限制自成交策略。  1. 用户在设置加入`STP用户组`后，可以通过传递 `stp_act` 来限制用户发生自成交的策略，没有传递 `stp_act` 默认按照 `cn` 的策略。 2. 用户在没有设置加入`STP用户组`时，传递 `stp_act` 参数会报错。 3. 用户没有使用 `stp_act` 发生成交的订单，`stp_act` 返回 `-`。  - cn: Cancel newest,取消新订单，保留老订单 - co: Cancel oldest,取消⽼订单，保留新订单 - cb: Cancel both,新旧订单都取消
	StpAct string `json:"stp_act,omitempty"`
	// 用户修改订单时备注的信息
	AmendText string `json:"amend_text,omitempty"`
	// 附加信息
	BizInfo string `json:"biz_info,omitempty"`
}
