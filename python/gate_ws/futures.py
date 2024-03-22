# !/usr/bin/env python
# coding: utf-8

from gate_ws.client import BaseChannel


class FuturesTickerChannel(BaseChannel):
    name = "futures.tickers"


class FuturesPublicTradeChannel(BaseChannel):
    name = "futures.trades"


class FuturesCandlesticksChannel(BaseChannel):
    name = "futures.candlesticks"


class FuturesBookTickerChannel(BaseChannel):
    name = "futures.book_ticker"


class FuturesOrderBookUpdateChannel(BaseChannel):
    name = "futures.order_book_update"


class FuturesOrderBookChannel(BaseChannel):
    name = "futures.order_book"


class FuturesOrderChannel(BaseChannel):
    name = "futures.orders"
    require_auth = True


class FuturesUserTradesChannel(BaseChannel):
    name = "futures.usertrades"
    require_auth = True


class FuturesLiquidatesChannel(BaseChannel):
    name = "futures.liquidates"
    require_auth = True


class FuturesADLChannel(BaseChannel):
    name = "futures.auto_deleverages"
    require_auth = True


class FuturesPositionClosesChannel(BaseChannel):
    name = "futures.position_closes"
    require_auth = True


class FuturesBalanceChannel(BaseChannel):
    name = "futures.balances"
    require_auth = True


class FuturesReduceRiskLimitChannel(BaseChannel):
    name = "futures.reduce_risk_limits"
    require_auth = True


class FuturesPositionsChannel(BaseChannel):
    name = "futures.positions"
    require_auth = True


class FuturesAutoOrdersChannel(BaseChannel):
    name = "futures.autoorders"
    require_auth = True


class FuturesLoginChannel(BaseChannel):
    name = "futures.login"


class FuturesOrderAmendChannel(BaseChannel):
    name = "futures.order_amend"


class FuturesOrderCancelChannel(BaseChannel):
    name = "futures.order_cancel"


class FuturesOrderCancelCpChannel(BaseChannel):
    name = "futures.order_cancel_cp"


class FuturesOrderPlaceChannel(BaseChannel):
    name = "futures.order_place"


class FuturesOrderBatchPlaceChannel(BaseChannel):
    name = "futures.order_batch_place"


class FuturesOrderStatusChannel(BaseChannel):
    name = "futures.order_status"


class FuturesOrderListChannel(BaseChannel):
    name = "futures.order_list"
