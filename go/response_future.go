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
	// Lowest trading price in recent 24h
	Low24h string `json:"low_24h,omitempty"`
	// Highest trading price in recent 24h
	High24h string `json:"high_24h,omitempty"`
	// Trade size in recent 24h
	Volume24h string `json:"volume_24h,omitempty"`
	// Trade volumes in recent 24h in BTC(deprecated, use `volume_24h_base`, `volume_24h_quote`, `volume_24h_settle` instead)
	Volume24hBtc string `json:"volume_24h_btc,omitempty"`
	// Trade volumes in recent 24h in USD(deprecated, use `volume_24h_base`, `volume_24h_quote`, `volume_24h_settle` instead)
	Volume24hUsd string `json:"volume_24h_usd,omitempty"`
	// Trade volume in recent 24h, in base currency
	Volume24hBase string `json:"volume_24h_base,omitempty"`
	// Trade volume in recent 24h, in quote currency
	Volume24hQuote string `json:"volume_24h_quote,omitempty"`
	// Trade volume in recent 24h, in settle currency
	Volume24hSettle string `json:"volume_24h_settle,omitempty"`
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
}

type FuturesTrade struct {
	// Trade ID
	Id int64 `json:"id,omitempty"`
	// Trading time
	CreateTime float64 `json:"create_time,omitempty"`
	// Trading time, with milliseconds set to 3 decimal places.
	CreateTimeMs float64 `json:"create_time_ms,omitempty"`
	// Futures contract
	Contract string `json:"contract,omitempty"`
	// Trading size
	Size int64 `json:"size,omitempty"`
	// Trading price
	Price string `json:"price,omitempty"`
}

type FuturesPriceTriggeredOrder struct {
	Initial FuturesInitialOrder `json:"initial"`
	Trigger FuturesPriceTrigger `json:"trigger"`
	// Auto order ID
	Id int64 `json:"id,omitempty"`
	// User ID
	User int32 `json:"user,omitempty"`
	// Creation time
	CreateTime float64 `json:"create_time,omitempty"`
	// Finished time
	FinishTime float64 `json:"finish_time,omitempty"`
	// ID of the newly created order on condition triggered
	TradeId int64 `json:"trade_id,omitempty"`
	// Order status.
	Status string `json:"status,omitempty"`
	// How order is finished
	FinishAs string `json:"finish_as,omitempty"`
	// Extra messages of how order is finished
	Reason string `json:"reason,omitempty"`
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
	FirstId      int64  `json:"U"`
	LastId       int64  `json:"u"`
	BestBidPrice string `json:"b"`
	BestBidSize  string `json:"B"`
	BestAskPrice string `json:"a"`
	BestAskSize  string `json:"A"`
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

type FuturesOrder struct {
	// Futures order ID
	Id int64 `json:"id,omitempty"`
	// User ID
	User int32 `json:"user,omitempty"`
	// Order creation time
	CreateTime float64 `json:"create_time,omitempty"`
	// Order finished time. Not returned if order is open
	FinishTime float64 `json:"finish_time,omitempty"`
	// How the order is finished.  - filled: all filled - cancelled: manually cancelled - liquidated: cancelled because of liquidation - ioc: time in force is `IOC`, finish immediately - auto_deleveraged: finished by ADL - reduce_only: cancelled because of increasing position while `reduce-only` set
	FinishAs string `json:"finish_as,omitempty"`
	// Order status  - `open`: waiting to be traded - `finished`: finished
	Status string `json:"status,omitempty"`
	// Futures contract
	Contract string `json:"contract"`
	// Order size. Specify positive number to make a bid, and negative number to ask
	Size int64 `json:"size"`
	// Display size for iceberg order. 0 for non-iceberg. Note that you would pay the taker fee for the hidden size
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
	// Time in force  - gtc: GoodTillCancelled - ioc: ImmediateOrCancelled, taker only - poc: PendingOrCancelled, reduce-only
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
}

type FuturesLiquidate struct {
	// Liquidation time
	Time int64 `json:"time,omitempty"`
	// Futures contract
	Contract string `json:"contract,omitempty"`
	// Position leverage. Not returned in public endpoints.
	Leverage string `json:"leverage,omitempty"`
	// Position size
	Size int64 `json:"size,omitempty"`
	// Position margin. Not returned in public endpoints.
	Margin string `json:"margin,omitempty"`
	// Average entry price. Not returned in public endpoints.
	EntryPrice string `json:"entry_price,omitempty"`
	// Liquidation price. Not returned in public endpoints.
	LiqPrice string `json:"liq_price,omitempty"`
	// Mark price. Not returned in public endpoints.
	MarkPrice string `json:"mark_price,omitempty"`
	// Liquidation order ID. Not returned in public endpoints.
	OrderId int64 `json:"order_id,omitempty"`
	// Liquidation order price
	OrderPrice string `json:"order_price,omitempty"`
	// Liquidation order average taker price
	FillPrice string `json:"fill_price,omitempty"`
	// Liquidation order maker size
	Left int64 `json:"left,omitempty"`
}
type FuturesInitialOrder struct {
	// Futures contract
	Contract string `json:"contract"`
	// Order size. Positive size means to buy, while negative one means to sell. Set to 0 to close the position
	Size int64 `json:"size,omitempty"`
	// Order price. Set to 0 to use market price
	Price string `json:"price"`
	// Set to true if trying to close the position
	Close bool `json:"close,omitempty"`
	// Time in force. If using market price, only `ioc` is supported.  - gtc: GoodTillCancelled - ioc: ImmediateOrCancelled
	Tif string `json:"tif,omitempty"`
	// How the order is created. Possible values are: web, api and app
	Text string `json:"text,omitempty"`
	// Set to true to create an reduce-only order
	ReduceOnly bool `json:"reduce_only,omitempty"`
	// Is the order reduce-only
	IsReduceOnly bool `json:"is_reduce_only,omitempty"`
	// Is the order to close position
	IsClose bool `json:"is_close,omitempty"`
}

type FuturesCandlestick struct {
	// Unix timestamp in seconds
	T float64 `json:"t,omitempty"`
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
	N string `json:"n"`
}

type FuturesAccountBook struct {
	// Change time
	Time float64 `json:"time,omitempty"`
	// Change amount
	Change string `json:"change,omitempty"`
	// Balance after change
	Balance string `json:"balance,omitempty"`
	// Changing Type: - dnw: Deposit & Withdraw - pnl: Profit & Loss by reducing position - fee: Trading fee - refr: Referrer rebate - fund: Funding - point_dnw: POINT Deposit & Withdraw - point_fee: POINT Trading fee - point_refr: POINT Referrer rebate
	Type string `json:"type,omitempty"`
	// Comment
	Text string `json:"text,omitempty"`
}

type FuturesAccount struct {
	// Total assets, total = position_margin + order_margin + available
	Total string `json:"total,omitempty"`
	// Unrealized PNL
	UnrealisedPnl string `json:"unrealised_pnl,omitempty"`
	// Position margin
	PositionMargin string `json:"position_margin,omitempty"`
	// Order margin of unfinished orders
	OrderMargin string `json:"order_margin,omitempty"`
	// Available balance to transfer out or trade
	Available string `json:"available,omitempty"`
	// POINT amount
	Point string `json:"point,omitempty"`
	// Settle currency
	Currency string `json:"currency,omitempty"`
	// Whether dual mode is enabled
	InDualMode bool `json:"in_dual_mode,omitempty"`
}

type FuturesUserTrade struct {
	Contract string `json:"contract"`
	// Trading time
	CreateTime float64 `json:"create_time,omitempty"`
	// Trading time, with milliseconds set to 3 decimal places.
	CreateTimeMs float64 `json:"create_time_ms,omitempty"`
	Id           string  `json:"id"`
	OrderId      string  `json:"order_id"`
	Price        string  `json:"price"`
	Size         int64   `json:"size"`
	Role         string  `json:"role"`
}
