# !/usr/bin/env python
# coding: utf-8

from gate_ws.client import BaseChannel


class PilotTickerChannel(BaseChannel):
    name = "pilot.tickers"


class PilotPublicTradeChannel(BaseChannel):
    name = "pilot.trades"


class PilotCandlesticksChannel(BaseChannel):
    name = "pilot.candlesticks"


class PilotBookTickerChannel(BaseChannel):
    name = "pilot.book_ticker"


class PilotOrderBookUpdateChannel(BaseChannel):
    name = "pilot.order_book_update"


class PilotOrderBookChannel(BaseChannel):
    name = "pilot.order_book"


class PilotOrderChannel(BaseChannel):
    name = "pilot.orders"
    require_auth = True


class PilotUserTradesChannel(BaseChannel):
    name = "pilot.usertrades"
    require_auth = True


class PilotBalanceChannel(BaseChannel):
    name = "pilot.balances"
    require_auth = True


class PilotMarginBalanceChannel(BaseChannel):
    name = "pilot.margin_balances"
    require_auth = True


class PilotFundingBalanceChannel(BaseChannel):
    name = "pilot.funding_balances"
    require_auth = True


class PilotCrossMarginBalanceChannel(BaseChannel):
    name = "pilot.cross_balances"
    require_auth = True


class PilotLoginChannel(BaseChannel):
    name = "pilot.login"


class PilotOrderAmendChannel(BaseChannel):
    name = "pilot.order_amend"


class PilotOrderCancelChannel(BaseChannel):
    name = "pilot.order_cancel"


class PilotOrderCancelCpChannel(BaseChannel):
    name = "pilot.order_cancel_cp"


class PilotOrderCancelIdsChannel(BaseChannel):
    name = "pilot.order_cancel_ids"


class PilotOrderPlaceChannel(BaseChannel):
    name = "pilot.order_place"


class PilotOrderStatusChannel(BaseChannel):
    name = "pilot.order_status"
