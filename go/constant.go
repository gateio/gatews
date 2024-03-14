package gatews

import "math"

const (
	BaseUrl        = "wss://api.gateio.ws/ws/v4/"
	FuturesBtcUrl  = "wss://fx-ws.gateio.ws/v4/ws/btc"
	FuturesUsdtUrl = "wss://fx-ws.gateio.ws/v4/ws/usdt"

	AuthMethodApiKey = "api_key"
	MaxRetryConn     = math.MaxInt64
)

// spot channels
const (
	ChannelSpotBalance         = "spot.balances"
	ChannelSpotCandleStick     = "spot.candlesticks"
	ChannelSpotOrder           = "spot.orders"
	ChannelSpotOrderBook       = "spot.order_book"
	ChannelSpotBookTicker      = "spot.book_ticker"
	ChannelSpotOrderBookUpdate = "spot.order_book_update"
	ChannelSpotTicker          = "spot.tickers"
	ChannelSpotUserTrade       = "spot.usertrades"
	ChannelSpotPublicTrade     = "spot.trades"
	ChannelSpotFundingBalance  = "spot.funding_balances"
	ChannelSpotMarginBalance   = "spot.margin_balances"
	ChannelSpotCrossBalance    = "spot.cross_balances"

	// order
	ChannelSpotLogin          = "spot.login"
	ChannelSpotOrderAmend     = "spot.order_amend"
	ChannelSpotOrderCancel    = "spot.order_cancel"
	ChannelSpotOrderCancelCp  = "spot.order_cancel_cp"
	ChannelSpotOrderCancelIds = "spot.order_cancel_ids"
	ChannelSpotOrderPlace     = "spot.order_place"
	ChannelSpotOrderStatus    = "spot.order_status"
)

// future channels
const (
	ChannelFutureTicker           = "futures.tickers"
	ChannelFutureTrade            = "futures.trades"
	ChannelFutureOrderBook        = "futures.order_book"
	ChannelFutureBookTicker       = "futures.book_ticker"
	ChannelFutureOrderBookUpdate  = "futures.order_book_update"
	ChannelFutureCandleStick      = "futures.candlesticks"
	ChannelFutureOrder            = "futures.orders"
	ChannelFutureUserTrade        = "futures.usertrades"
	ChannelFutureLiquidates       = "futures.liquidates"
	ChannelFutureAutoDeleverages  = "futures.auto_deleverages"
	ChannelFuturePositionCloses   = "futures.position_closes"
	ChannelFutureBalance          = "futures.balances"
	ChannelFutureReduceRiskLimits = "futures.reduce_risk_limits"
	ChannelFuturePositions        = "futures.positions"
	ChannelFutureAutoOrders       = "futures.autoorders"

	// order
	ChannelFutureLogin           = "futures.login"
	ChannelFutureOrderAmend      = "futures.order_amend"
	ChannelFutureOrderCancel     = "futures.order_cancel"
	ChannelFutureOrderCancelCp   = "futures.order_cancel_cp"
	ChannelFutureOrderPlace      = "futures.order_place"
	ChannelFutureOrderBatchPlace = "futures.order_batch_place"
	ChannelFutureOrderStatus     = "futures.order_status"
	ChannelFutureOrderList       = "futures.order_list"
)

var authChannel = map[string]bool{
	// spot
	ChannelSpotBalance:        true,
	ChannelSpotFundingBalance: true,
	ChannelSpotMarginBalance:  true,
	ChannelSpotOrder:          true,
	ChannelSpotUserTrade:      true,

	// future
	ChannelFutureOrder:            true,
	ChannelFutureUserTrade:        true,
	ChannelFutureLiquidates:       true,
	ChannelFutureAutoDeleverages:  true,
	ChannelFuturePositionCloses:   true,
	ChannelFutureReduceRiskLimits: true,
	ChannelFuturePositions:        true,
	ChannelFutureAutoOrders:       true,
	ChannelFutureBalance:          true,
}

const (
	Subscribe   = "subscribe"
	UnSubscribe = "unsubscribe"
	API         = "api"

	ServiceTypeSpot    = 1
	ServiceTypeFutures = 2

	DefaultPingInterval = "10s"
)
