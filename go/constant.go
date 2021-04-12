package gatews

const (
	BaseUrl          = "wss://api.gateio.ws/ws/v4/"
	AuthMethodApiKey = "api_key"
	MaxRetryConn     = 10
)

// channels
const (
	ChannelBalance         = "spot.balances"
	ChannelCandleStick     = "spot.candlesticks"
	ChannelOrder           = "spot.orders"
	ChannelOrderBook       = "spot.order_book"
	ChannelBookTicker      = "spot.book_ticker"
	ChannelOrderBookUpdate = "spot.order_book_update"
	ChannelTicker          = "spot.tickers"
	ChannelUserTrade       = "spot.usertrades"
	ChannelPublicTrade     = "spot.trades"
	ChannelFundingBalance  = "spot.funding_balances"
	ChannelMarginBalance   = "spot.margin_balances"
)

var (
	authChannel = map[string]bool{
		ChannelBalance:        true,
		ChannelFundingBalance: true,
		ChannelMarginBalance:  true,
		ChannelOrder:          true,
		ChannelUserTrade:      true,
	}
)

const (
	Subscribe   = "subscribe"
	UnSubscribe = "unsubscribe"
)
