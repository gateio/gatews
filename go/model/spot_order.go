package model

type Order struct {
	// Order ID
	Id string `json:"id,omitempty"`
	// User defined information. If not empty, must follow the rules below:  1. prefixed with `t-` 2. no longer than 28 bytes without `t-` prefix 3. can only include 0-9, A-Z, a-z, underscore(_), hyphen(-) or dot(.)  Besides user defined information, reserved contents are listed below, denoting how the order is created:  - 101: from android - 102: from IOS - 103: from IPAD - 104: from webapp - 3: from web - 2: from apiv2 - apiv4: from apiv4
	Text string `json:"text,omitempty"`
	// 用户修改订单时备注的信息
	AmendText string `json:"amend_text,omitempty"`
	// Creation time of order
	CreateTime string `json:"create_time,omitempty"`
	// Last modification time of order
	UpdateTime string `json:"update_time,omitempty"`
	// Creation time of order (in milliseconds)
	CreateTimeMs int64 `json:"create_time_ms,omitempty"`
	// Last modification time of order (in milliseconds)
	UpdateTimeMs int64 `json:"update_time_ms,omitempty"`
	// Order status  - `open`: to be filled - `closed`: filled - `cancelled`: cancelled
	Status string `json:"status,omitempty"`
	// Currency pair
	CurrencyPair string `json:"currency_pair"`
	// Order Type    - limit : Limit Order - market : Market Order
	Type string `json:"type,omitempty"`
	// 账户类型，spot - 现货账户，margin - 杠杆账户，cross_margin - 全仓杠杆账户，unified - 统一账户 统一账户（旧）只能设置 `cross_margin`
	Account string `json:"account,omitempty"`
	// Order side
	Side string `json:"side"`
	// When `type` is limit, it refers to base currency.  For instance, `BTC_USDT` means `BTC`  When `type` is `market`, it refers to different currency according to `side`  - `side` : `buy` means quote currency, `BTC_USDT` means `USDT` - `side` : `sell` means base currency，`BTC_USDT` means `BTC`
	Amount string `json:"amount"`
	// Price can't be empty when `type`= `limit`
	Price string `json:"price,omitempty"`
	// Time in force  - gtc: GoodTillCancelled - ioc: ImmediateOrCancelled, taker only - poc: PendingOrCancelled, makes a post-only order that always enjoys a maker fee - fok: FillOrKill, fill either completely or none Only `ioc` and `fok` are supported when `type`=`market`
	TimeInForce string `json:"time_in_force,omitempty"`
	// Amount to display for the iceberg order. Null or 0 for normal orders.  Hiding all amount is not supported.
	Iceberg string `json:"iceberg,omitempty"`
	// Used in margin or cross margin trading to allow automatic loan of insufficient amount if balance is not enough.
	AutoBorrow bool `json:"auto_borrow,omitempty"`
	// Enable or disable automatic repayment for automatic borrow loan generated by cross margin order. Default is disabled. Note that:  1. This field is only effective for cross margin orders. Margin account does not support setting auto repayment for orders. 2. `auto_borrow` and `auto_repay` cannot be both set to true in one order.
	AutoRepay bool `json:"auto_repay,omitempty"`
	// Amount left to fill
	Left string `json:"left,omitempty"`
	// Total filled in quote currency. Deprecated in favor of `filled_total`
	FillPrice string `json:"fill_price,omitempty"`
	// Total filled in quote currency
	FilledTotal string `json:"filled_total,omitempty"`
	// Average fill price
	AvgDealPrice string `json:"avg_deal_price,omitempty"`
	// Fee deducted
	Fee string `json:"fee,omitempty"`
	// Fee currency unit
	FeeCurrency string `json:"fee_currency,omitempty"`
	// Points used to deduct fee
	PointFee string `json:"point_fee,omitempty"`
	// GT used to deduct fee
	GtFee string `json:"gt_fee,omitempty"`
	// GT used to deduct maker fee
	GtMakerFee string `json:"gt_maker_fee,omitempty"`
	// GT used to deduct taker fee
	GtTakerFee string `json:"gt_taker_fee,omitempty"`
	// Whether GT fee discount is used
	GtDiscount bool `json:"gt_discount,omitempty"`
	// Rebated fee
	RebatedFee string `json:"rebated_fee,omitempty"`
	// Rebated fee currency unit
	RebatedFeeCurrency string `json:"rebated_fee_currency,omitempty"`
	// 订单所属的`STP用户组`id，同一个`STP用户组`内用户之间的订单不允许发生自成交。  1. 如果撮合时两个订单的 `stp_id` 非 `0` 且相等，则不成交，而是根据 `taker` 的 `stp_act` 执行相应策略。 2. 没有设置`STP用户组`成交的订单，`stp_id` 默认返回 `0`。
	StpId int32 `json:"stp_id,omitempty"`
	// Self-Trading Prevention Action,用户可以用该字段设置自定义限制自成交策略。  1. 用户在设置加入`STP用户组`后，可以通过传递 `stp_act` 来限制用户发生自成交的策略，没有传递 `stp_act` 默认按照 `cn` 的策略。 2. 用户在没有设置加入`STP用户组`时，传递 `stp_act` 参数会报错。 3. 用户没有使用 `stp_act` 发生成交的订单，`stp_act` 返回 `-`。  - cn: Cancel newest,取消新订单，保留老订单 - co: Cancel oldest,取消⽼订单，保留新订单 - cb: Cancel both,新旧订单都取消
	StpAct string `json:"stp_act,omitempty"`
	// 订单结束方式，包括：  - open: 等待处理 - filled: 完全成交 - cancelled: 用户撤销 - ioc: 未立即完全成交，因为 tif 设置为 ioc - stp: 订单发生自成交限制而被撤销
	FinishAs string `json:"finish_as,omitempty"`
	// 费率折扣
	FeeDiscount string `json:"fee_discount,omitempty"`
	// 处理模式: 下单时根据action_mode返回不同的字段, 该字段只在请求时有效，响应结果中不包含该字段 `ACK`: 异步模式，只返回订单关键字段 `RESULT`: 无清算信息 `FULL`: 完整模式（默认）
	ActionMode string `json:"action_mode,omitempty"`
}