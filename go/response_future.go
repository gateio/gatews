package gatews

type FuturesTicker struct {
	// Futures contract
	Contract string `json:"contract,omitempty"`
	// Last trading price
	Last string `json:"last,omitempty"`
	// Change percentage.
	ChangePercentage string `json:"change_percentage,omitempty"`
	// Contract total size
	TotalSize string `json:"total_size,omitempty"`
	// Trade size in recent 24h
	Volume24h string `json:"volume_24h,omitempty"`
	// Trade volume in recent 24h, in base currency
	Volume24hBase string `json:"volume_24h_base,omitempty"`
	// Trade volume in recent 24h, in quote currency
	Volume24hQuote string `json:"volume_24h_quote,omitempty"`
	// Trade volume in recent 24h, in settle currency
	Volume24hSettle string `json:"volume_24h_settle,omitempty"`
	Volume24Usd     string `json:"volume_24_usd"`
	Volume24Btc     string `json:"volume_24_btc"`
	// Recent mark price
	MarkPrice string `json:"mark_price,omitempty"`
	// Funding rate
	FundingRate string `json:"funding_rate,omitempty"`
	// Indicative Funding rate in next period
	FundingRateIndicative string `json:"funding_rate_indicative,omitempty"`
	// Index price
	IndexPrice string `json:"index_price,omitempty"`
	// Exchange rate of base currency and settlement currency in Quanto contract. Not existed in contract of other types
	QuantoBaseRate string `json:"quanto_base_rate,omitempty"`
	Low24h         string `json:"low_24h"`
	High24h        string `json:"high_24h"`
}

type FuturesTrade struct {
	// Trade ID
	Id int64 `json:"id,omitempty"`
	// Trading time
	CreateTime int64 `json:"create_time,omitempty"`
	// Trading time, with milliseconds set to 3 decimal places.
	CreateTimeMs int64 `json:"create_time_ms,omitempty"`
	// Futures contract
	Contract string `json:"contract,omitempty"`
	// Trading size
	Size int64 `json:"size,omitempty"`
	// Trading price
	Price string `json:"price,omitempty"`
}

type FuturesOrderBookItem struct {
	// Price
	P string `json:"p,omitempty"`
	// Size
	S int64 `json:"s,omitempty"`
}

type FuturesOrderBook struct {
	// Order Book ID. Increase by 1 on every order book change. Set `with_id=true` to include this field in response
	Id       int64  `json:"id,omitempty"`
	Contract string `json:"contract"`
	Time     int64  `json:"t"`
	// Asks order depth
	Asks []FuturesOrderBookItem `json:"asks"`
	// Bids order depth
	Bids []FuturesOrderBookItem `json:"bids"`
}

type FuturesOrderBookAll struct {
	Contract string `json:"c"`
	Price    string `json:"p"`
	Id       int64  `json:"id"`
	Size     int64  `json:"s"`
}

type FuturesBookTicker struct {
	TimeMillis   int64  `json:"t"`
	Contract     string `json:"s"`
	UpdateId     int64  `json:"u"`
	BestBidPrice string `json:"b"`
	BestBidSize  int64  `json:"B"`
	BestAskPrice string `json:"a"`
	BestAskSize  int64  `json:"A"`
}

type FuturesOrderBookUpdate struct {
	TimeMillis int64  `json:"t"`
	Contract   string `json:"s"`
	FirstId    int64  `json:"U"`
	LastId     int64  `json:"u"`
	// Asks order depth
	Asks []FuturesOrderBookItem `json:"a"`
	// Bids order depth
	Bids []FuturesOrderBookItem `json:"b"`
}

type FuturesCandlestick struct {
	// Unix timestamp in seconds
	T int64 `json:"t,omitempty"`
	// size volume. Only returned if `contract` is not prefixed
	V int64 `json:"v,omitempty"`
	// Close price
	C string `json:"c,omitempty"`
	// Highest price
	H string `json:"h,omitempty"`
	// Lowest price
	L string `json:"l,omitempty"`
	// Open price
	O string `json:"o,omitempty"`
	// futures contract name
	N      string `json:"n"`
	Amount string `json:"a"`
}

type FuturesOrder struct {
	// Futures order ID
	Id int64 `json:"id,omitempty"`
	// User ID
	User string `json:"user,omitempty"`
	// Order creation time
	CreateTime   int64 `json:"create_time,omitempty"`
	CreateTimeMs int64 `json:"create_time_ms,omitempty"`
	// Order finished time. Not returned if order is open
	FinishTime   int64 `json:"finish_time,omitempty"`
	FinishTimeMs int64 `json:"finish_time_ms,omitempty"`
	// FinishAs indicates how the order was completed:
	// - filled: all filled
	// - cancelled: manually cancelled
	// - liquidated: cancelled due to liquidation
	// - ioc: time in force is IOC, finished immediately
	// - auto_deleveraged: finished by ADL
	// - reduce_only: cancelled due to increase in position while reduce-only set
	// - position_closed: cancelled due to position close
	// - stp: cancelled due to self trade prevention
	// - _new: order created
	// - _update: order filled, partially filled, or updated
	// - reduce_out: only reduce position, excluding pending orders hard to execute
	FinishAs string `json:"finish_as,omitempty"`
	// Futures contract
	Contract string `json:"contract"`
	// Order size. Specify positive number to make a bid, and negative number to ask
	Size int64 `json:"size"`
	// Display size for iceberg order. 0 for non-iceberg. Note that you would pay the taker fee for the hidden size
	Iceberg int64 `json:"iceberg,omitempty"`
	// Order price. 0 for market order with `tif` set as `ioc`
	Price float64 `json:"price,omitempty"`
	// Is the order to close position
	IsClose bool `json:"is_close,omitempty"`
	// Is the order reduce-only
	IsReduceOnly bool `json:"is_reduce_only,omitempty"`
	// Is the order for liquidation
	IsLiq bool `json:"is_liq,omitempty"`
	// Time in force  - gtc: GoodTillCancelled - ioc: ImmediateOrCancelled, taker only - poc: PendingOrCancelled, reduce-only
	Tif string `json:"tif,omitempty"`
	// Size left to be traded
	Left int64 `json:"left,omitempty"`
	// Fill price of the order
	FillPrice float64 `json:"fill_price,omitempty"`
	// User defined information. If not empty, must follow the rules below:  1. prefixed with `t-` 2. no longer than 28 bytes without `t-` prefix 3. can only include 0-9, A-Z, a-z, underscore(_), hyphen(-) or dot(.) Besides user defined information, reserved contents are listed below, denoting how the order is created:  - web: from web - api: from API - app: from mobile phones - auto_deleveraging: from ADL - liquidation: from liquidation - insurance: from insurance
	Text string `json:"text,omitempty"`
	// Taker fee
	Tkfr float64 `json:"tkfr,omitempty"`
	// Maker fee
	Mkfr float64 `json:"mkfr,omitempty"`
	// Reference user ID
	Refu int32   `json:"refu,omitempty"`
	Refr float64 `json:"refr"`

	StopProfitPrice string `json:"stop_profit_price"`
	StopLossPrice   string `json:"stop_loss_price"`
	// StpId represents the ID associated with the self-trade prevention mechanism.
	StpId int64 `json:"stp_id,omitempty"`
	// StpAct represents the self-trade prevention (STP) action:
	// - cn: Cancel newest (keep old orders)
	// - co: Cancel oldest (keep new orders)
	// - cb: Cancel both (cancel both old and new orders)
	// If not provided, defaults to 'cn'. Requires STP group membership; otherwise, an error is returned.
	StpAct string `json:"stp_act,omitempty"`
	// BizInfo represents business-specific information related to the order. The exact content and format can vary depending on the use case.
	BizInfo string `json:"biz_info,omitempty"`
	// AmendText provides the custom data that the user remarked when amending the order
	AmendText string `json:"amend_text,omitempty"`
}

type FuturesUserTrade struct {
	Contract string `json:"contract"`
	// Trading time
	CreateTime int64 `json:"create_time,omitempty"`
	// Trading time, with milliseconds set to 3 decimal places.
	CreateTimeMs int64   `json:"create_time_ms,omitempty"`
	Id           string  `json:"id"`
	OrderId      string  `json:"order_id"`
	Price        string  `json:"price"`
	Size         int64   `json:"size"`
	Role         string  `json:"role"`
	Text         string  `json:"text"`
	Fee          float64 `json:"fee"`
	PointFee     float64 `json:"point_fee"`
}

type FuturesLiquidate struct {
	// Liquidation time
	Time int64 `json:"time,omitempty"`
	// time in milliseconds
	TimeMs int64 `json:"time_ms"`
	// Futures contract
	Contract string `json:"contract,omitempty"`
	// Position leverage. Not returned in public endpoints.
	Leverage float64 `json:"leverage,omitempty"`
	// Position size
	Size int64 `json:"size,omitempty"`
	// Position margin. Not returned in public endpoints.
	Margin float64 `json:"margin,omitempty"`
	// Average entry price. Not returned in public endpoints.
	EntryPrice float64 `json:"entry_price,omitempty"`
	// Liquidation price. Not returned in public endpoints.
	LiqPrice float64 `json:"liq_price,omitempty"`
	// Mark price. Not returned in public endpoints.
	MarkPrice float64 `json:"mark_price,omitempty"`
	// Liquidation order ID. Not returned in public endpoints.
	OrderId int64 `json:"order_id,omitempty"`
	// Liquidation order price
	OrderPrice float64 `json:"order_price,omitempty"`
	// Liquidation order average taker price
	FillPrice float64 `json:"fill_price,omitempty"`
	// Liquidation order maker size
	Left int64 `json:"left,omitempty"`
	// user id
	User string `json:"user"`
}

type FuturesAutoDeleverages struct {
	EntryPrice   float64 `json:"entry_price"`
	FillPrice    float64 `json:"fill_price"`
	PositionSize int64   `json:"position_size"`
	TradeSize    int64   `json:"trade_size"`
	Time         int64   `json:"time"`
	TimeMs       int64   `json:"time_ms"`
	Contract     string  `json:"contract"`
	User         string  `json:"user"`
}

type FuturesPositionCloses struct {
	Contract string  `json:"contract"`
	Pnl      float64 `json:"pnl"`
	Side     string  `json:"side"`
	Text     string  `json:"text"`
	Time     int64   `json:"time"`
	TimeMs   int64   `json:"time_ms"`
	User     string  `json:"user"`
}

type FuturesBalance struct {
	Balance  float64 `json:"balance"`
	Change   float64 `json:"change"`
	Text     string  `json:"text"`
	Time     int64   `json:"time"`
	TimeMs   int64   `json:"time_ms"`
	User     string  `json:"user"`
	Type     string  `json:"type"`
	Currency string  `json:"currency"`
}

type FuturesReduceRiskLimits struct {
	CancelOrders    int64   `json:"cancel_orders"`
	Contract        string  `json:"contract"`
	LeverageMax     float64 `json:"leverage_max"`
	LiqPrice        float64 `json:"liq_price"`
	MaintenanceRate float64 `json:"maintenance_rate"`
	RiskLimit       float64 `json:"risk_limit"`
	Time            int64   `json:"time"`
	TimeMs          int64   `json:"time_ms"`
	User            string  `json:"user"`
}

type FuturesPositions struct {
	Contract           string  `json:"contract"`
	CrossLeverageLimit float64 `json:"cross_leverage_limit"`
	EntryPrice         float64 `json:"entry_price"`
	HistoryPnl         float64 `json:"history_pnl"`
	HistoryPoint       float64 `json:"history_point"`
	LastClosePnl       float64 `json:"last_close_pnl"`
	Leverage           float64 `json:"leverage"`
	LeverageMax        float64 `json:"leverage_max"`
	LiqPrice           float64 `json:"liq_price"`
	MaintenanceRate    float64 `json:"maintenance_rate"`
	Margin             float64 `json:"margin"`
	Mode               string  `json:"mode"`
	RealisedPnl        float64 `json:"realised_pnl"`
	RealisedPoint      float64 `json:"realised_point"`
	RiskLimit          float64 `json:"risk_limit"`
	Size               int64   `json:"size"`
	Time               int64   `json:"time"`
	TimeMs             int64   `json:"time_ms"`
	User               string  `json:"user"`
}

type FuturesAutoOrder struct {
	Initial     FuturesInitialOrder `json:"initial"`
	Trigger     FuturesPriceTrigger `json:"trigger"`
	StopTrigger FutureStopTrigger   `json:"stop_trigger"`
	// Auto order ID
	Id int64 `json:"id,omitempty"`
	// User ID
	User int64 `json:"user,omitempty"`
	// Creation time
	CreateTime int64 `json:"create_time,omitempty"`
	// Finished time
	FinishTime int64 `json:"finish_time,omitempty"`
	// ID of the newly created order on condition triggered
	TradeId int64 `json:"trade_id,omitempty"`
	// Order status.
	Status string `json:"status,omitempty"`
	// Extra messages of how order is finished
	Reason      string `json:"reason,omitempty"`
	Name        string `json:"name"`
	IsStopOrder bool   `json:"is_stop_order"`
	FinishAs    string `json:"finish_as"`
	MeOrderId   int64  `json:"me_order_id"`
	OrderType   string `json:"order_type"`
}

type FutureStopTrigger struct {
	Rule         int32  `json:"rule"`
	TriggerPrice string `json:"trigger_price"`
	OrderPrice   string `json:"order_price"`
}

type FuturesPriceTrigger struct {
	// How the order will be triggered   - `0`: by price, which means order will be triggered on price condition satisfied  - `1`: by price gap, which means order will be triggered on gap of recent two prices of specified `price_type` satisfied.  Only `0` is supported currently
	StrategyType int32 `json:"strategy_type,omitempty"`
	// Price type. 0 - latest deal price, 1 - mark price, 2 - index price
	PriceType int32 `json:"price_type,omitempty"`
	// Value of price on price triggered, or price gap on price gap triggered
	Price string `json:"price,omitempty"`
	// Trigger condition type  - `1`: calculated price based on `strategy_type` and `price_type` >= `price` - `2`: calculated price based on `strategy_type` and `price_type` <= `price`
	Rule int32 `json:"rule,omitempty"`
	// How many seconds will the order wait for the condition being triggered. Order will be cancelled on timed out
	Expiration int32 `json:"expiration,omitempty"`
}

type FuturesInitialOrder struct {
	// Futures contract
	Contract string `json:"contract"`
	// Order size. Positive size means to buy, while negative one means to sell. Set to 0 to close the position
	Size int64 `json:"size,omitempty"`
	// Order price. Set to 0 to use market price
	Price string `json:"price"`
	// Time in force. If using market price, only `ioc` is supported.  - gtc: GoodTillCancelled - ioc: ImmediateOrCancelled
	Tif string `json:"tif,omitempty"`
	// How the order is created. Possible values are: web, api and app
	Text    string `json:"text,omitempty"`
	Iceberg int64  `json:"iceberg"`
	// Is the order reduce-only
	IsReduceOnly bool `json:"is_reduce_only,omitempty"`
	// Is the order to close position
	IsClose  bool   `json:"is_close,omitempty"`
	AutoSize string `json:"auto_size"`
}
