# !/usr/bin/env python
# coding: utf-8

from gate_ws.client import BaseChannel


class SpotTickerChannel(BaseChannel):
    name = 'spot.tickers'


class SpotPublicTradeChannel(BaseChannel):
    name = 'spot.trades'


class SpotCandlesticksChannel(BaseChannel):
    name = 'spot.candlesticks'


class SpotBookTickerChannel(BaseChannel):
    name = 'spot.book_ticker'


class SpotOrderBookUpdateChannel(BaseChannel):
    name = 'spot.order_book_update'


class SpotOrderBookChannel(BaseChannel):
    name = 'spot.order_book'


class SpotOrderChannel(BaseChannel):
    name = 'spot.orders'
    require_auth = True


class SpotUserTradesChannel(BaseChannel):
    name = 'spot.usertrades'
    require_auth = True


class SpotBalanceChannel(BaseChannel):
    name = 'spot.balances'
    require_auth = True


class SpotMarginBalanceChannel(BaseChannel):
    name = 'spot.margin_balances'
    require_auth = True


class SpotFundingBalanceChannel(BaseChannel):
    name = 'spot.funding_balances'
    require_auth = True
