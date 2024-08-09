package gatews

type SpotBalancesMsg struct {
	Timestamp        string `json:"timestamp"`
	TimestampInMilli string `json:"timestamp_ms"`
	User             string `json:"user"`
	Asset            string `json:"currency"`
	Change           string `json:"change"`
	Total            string `json:"total"`
	Available        string `json:"available"`
	Freeze           string `json:"freeze"`
	FreezeChange     string `json:"freeze_change"`
	ChangeType       string `json:"change_type"`
}

type SpotCandleUpdateMsg struct {
	Time        string `json:"t"`
	Volume      string `json:"v"`
	Close       string `json:"c"`
	High        string `json:"h"`
	Low         string `json:"l"`
	Open        string `json:"o"`
	Name        string `json:"n"`
	Amount      string `json:"a"`
	WindowClose bool   `json:"w"`
}

// SpotUpdateDepthMsg update order book
type SpotUpdateDepthMsg struct {
	TimeInMilli  int64      `json:"t"`
	Event        string     `json:"e"`
	ETime        int64      `json:"E"`
	CurrencyPair string     `json:"s"`
	FirstId      int64      `json:"U"`
	LastId       int64      `json:"u"`
	Bid          [][]string `json:"b"`
	Ask          [][]string `json:"a"`
}

// SpotUpdateAllDepthMsg all order book
type SpotUpdateAllDepthMsg struct {
	TimeInMilli  int64       `json:"t"`
	LastUpdateId int64       `json:"lastUpdateId"`
	CurrencyPair string      `json:"s"`
	Bid          [][2]string `json:"bids"`
	Ask          [][2]string `json:"asks"`
}

type SpotFundingBalancesMsg struct {
	Timestamp        string `json:"timestamp"`
	TimestampInMilli string `json:"timestamp_ms"`
	User             string `json:"user"`
	Asset            string `json:"currency"`
	Change           string `json:"change"`
	Freeze           string `json:"freeze"`
	Lent             string `json:"lent"`
}

type SpotMarginBalancesMsg struct {
	Timestamp        string `json:"timestamp"`
	TimestampInMilli string `json:"timestamp_ms"`
	User             string `json:"user"`
	Market           string `json:"currency_pair"`
	Asset            string `json:"currency"`
	Change           string `json:"change"`
	Available        string `json:"available"`
	Freeze           string `json:"freeze"`
	Borrowed         string `json:"borrowed"`
	Interest         string `json:"interest"`
}

type SpotBookTickerMsg struct {
	TimeInMilli  int64  `json:"t"`
	LastId       int64  `json:"u"`
	CurrencyPair string `json:"s"`
	Bid          string `json:"b"`
	BidSize      string `json:"B"`
	Ask          string `json:"a"`
	AskSize      string `json:"A"`
}

type SpotTickerMsg struct {
	// Currency pair
	CurrencyPair string `json:"currency_pair,omitempty"`
	// Last trading price
	Last string `json:"last,omitempty"`
	// Lowest ask
	LowestAsk string `json:"lowest_ask,omitempty"`
	// Highest bid
	HighestBid string `json:"highest_bid,omitempty"`
	// Change percentage.
	ChangePercentage string `json:"change_percentage,omitempty"`
	// Base currency trade volume
	BaseVolume string `json:"base_volume,omitempty"`
	// Quote currency trade volume
	QuoteVolume string `json:"quote_volume,omitempty"`
	// Highest price in 24h
	High24h string `json:"high_24h,omitempty"`
	// Lowest price in 24h
	Low24h string `json:"low_24h,omitempty"`
}

type SpotUserTradesMsg struct {
	Id           uint64 `json:"id"`
	UserId       uint64 `json:"user_id"`
	OrderId      string `json:"order_id"`
	CurrencyPair string `json:"currency_pair"`
	CreateTime   int64  `json:"create_time"`
	CreateTimeMs string `json:"create_time_ms"`
	Side         string `json:"side"`
	Amount       string `json:"amount"`
	Role         string `json:"role"`
	Price        string `json:"price"`
	Fee          string `json:"fee"`
	FeeCurrency  string `json:"fee_currency"`
	PointFee     string `json:"point_fee"`
	GtFee        string `json:"gt_fee"`
	Text         string `json:"text"`
	AmendText    string `json:"amend_text"`
	BizInfo      string `json:"biz_info"`
}

type SpotTradeMsg struct {
	Id           uint64 `json:"id"`
	CreateTime   int64  `json:"create_time"`
	CreateTimeMs string `json:"create_time_ms"`
	Side         string `json:"side"`
	CurrencyPair string `json:"currency_pair"`
	Amount       string `json:"amount"`
	Price        string `json:"price"`
}

type OrderMsg struct {
	// SpotOrderMsg ID
	Id string `json:"id,omitempty"`
	// User defined information. If not empty, must follow the rules below:  1. prefixed with `t-` 2. no longer than 28 bytes without `t-` prefix 3. can only include 0-9, A-Z, a-z, underscore(_), hyphen(-) or dot(.)
	Text string `json:"text,omitempty"`
	// SpotOrderMsg creation time
	CreateTime string `json:"create_time,omitempty"`
	// SpotOrderMsg last modification time
	UpdateTime string `json:"update_time,omitempty"`
	// Currency pair
	CurrencyPair string `json:"currency_pair"`
	// SpotOrderMsg type. limit - limit order
	Type string `json:"type,omitempty"`
	// Account type. spot - use spot account; margin - use margin account
	Account string `json:"account,omitempty"`
	// SpotOrderMsg side
	Side string `json:"side"`
	// SpotTradeMsg amount
	Amount string `json:"amount"`
	// SpotOrderMsg price
	Price string `json:"price"`
	// Time in force  - gtc: GoodTillCancelled - ioc: ImmediateOrCancelled, taker only - poc: PendingOrCancelled, makes a post-only order that always enjoys a maker fee
	TimeInForce string `json:"time_in_force,omitempty"`
	// Amount to display for the iceberg order. Null or 0 for normal orders
	Iceberg string `json:"iceberg,omitempty"`
	// Used in margin trading(i.e. `account` is `margin`) to allow automatic loan of insufficient part if balance is not enough.
	AutoBorrow bool `json:"auto_borrow,omitempty"`
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
	// Point used to deduct fee
	PointFee string `json:"point_fee,omitempty"`
	// GT used to deduct fee
	GtFee string `json:"gt_fee,omitempty"`
	// Whether GT fee discount is used
	GtDiscount bool `json:"gt_discount,omitempty"`
	// Rebated fee
	RebatedFee string `json:"rebated_fee,omitempty"`
	// Rebated fee currency unit
	RebatedFeeCurrency string `json:"rebated_fee_currency,omitempty"`
	// StpId represents the ID associated with the self-trade prevention mechanism.
	StpId int64 `json:"stp_id,omitempty"`
	// StpAct represents the self-trade prevention (STP) action:
	// - cn: Cancel newest (keep old orders)
	// - co: Cancel oldest (keep new orders)
	// - cb: Cancel both (cancel both old and new orders)
	// If not provided, defaults to 'cn'. Requires STP group membership; otherwise, an error is returned.
	StpAct string `json:"stp_act,omitempty"`
	// FinishAs indicates how the order was finished:
	// - open: processing
	// - filled: fully filled
	// - cancelled: manually cancelled
	// - ioc: finished immediately (IOC)
	// - stp: cancelled due to self-trade prevention
	FinishAs string `json:"finish_as,omitempty"`
	// BizInfo represents business-specific information related to the order. The exact content and format can vary depending on the use case.
	BizInfo string `json:"biz_info,omitempty"`
	// AmendText provides the custom data that the user remarked when amending the order
	AmendText string `json:"amend_text,omitempty"`
}

type SpotOrderMsg struct {
	OrderMsg
	CreateTimeMs string `json:"create_time_ms,omitempty"`
	UpdateTimeMs string `json:"update_time_ms,omitempty"`
	User         int64  `json:"user"`
	Event        string `json:"event"`
}
