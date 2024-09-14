# !/usr/bin/env python
# coding: utf-8

"""
Example of how to maintain a local pilot order book
"""

import asyncio
import itertools
import logging
import sys
import typing
from collections import defaultdict
from datetime import datetime
from decimal import Decimal

try:
    import aiohttp
    from asciimatics.parsers import AsciimaticsParser
    from asciimatics.scene import Scene
    from asciimatics.screen import Screen
    from asciimatics.widgets import Frame, Layout, MultiColumnListBox, TextBox
    from sortedcontainers import SortedList
except ImportError:
    sys.stderr.write("aiohttp, sortedcontainers, and asciimatics are required\n")
    sys.exit(1)

from gate_ws import Configuration, Connection, WebSocketResponse
from gate_ws.pilot import PilotOrderBookUpdateChannel

logger = logging.getLogger(__name__)


class SimpleRingBuffer(object):
    """Simple ring buffer to cache order book updates

    But can be used in other general scenario too
    """

    def __init__(self, size: int):
        self.max = size
        self.data = []
        self.cur = 0

    class __Full:
        # to avoid warning hints from IDE
        max: int
        data: typing.List
        cur: int

        def append(self, x):
            self.data[self.cur] = x
            self.cur = (self.cur + 1) % self.max

        def __iter__(self):
            for i in itertools.chain(range(self.cur, self.max), range(self.cur)):
                yield self.data[i]

        def get(self, idx):
            return self.data[(self.cur + idx) % self.max]

        def __getitem__(self, item):
            if isinstance(item, int):
                return self.get(item)
            return (self.data[self.cur :] + self.data[: self.cur]).__getitem__(item)

        def __len__(self):
            return self.max

    def __iter__(self):
        for i in self.data:
            yield i

    def append(self, x):
        self.data.append(x)
        if len(self.data) == self.max:
            self.cur = 0
            # Permanently change self's class from non-full to full
            self.__class__ = self.__Full

    def get(self, idx):
        return self.data[idx]

    def __getitem__(self, item):
        return self.data.__getitem__(item)

    def __len__(self):
        return len(self.data)


class OrderBookEntry(object):
    def __init__(self, price, amount):
        self.price: Decimal = Decimal(price)
        self.amount: str = amount

    def __eq__(self, other):
        return self.price == other.price

    def __str__(self):
        return "(%s, %s)" % (self.price, self.amount)


class OrderBook(object):
    def __init__(self, cp: str, last_id: id, asks: SortedList, bids: SortedList):
        self.cp = cp
        self.id = last_id
        self.asks = asks
        self.bids = bids

    @classmethod
    def update_entry(cls, book: SortedList, entry: OrderBookEntry):
        if Decimal(entry.amount) == Decimal("0"):
            # remove price if amount is 0
            try:
                book.remove(entry)
            except ValueError:
                pass
        else:
            try:
                idx = book.index(entry)
            except ValueError:
                # price not found, insert it
                book.add(entry)
            else:
                # price found, update amount
                book[idx].amount = entry.amount

    def __str__(self):
        return "\n  id: %d\n  asks:\n%s\n  bids:\n%s" % (
            self.id,
            "\n".join([" " * 4 + str(a) for a in self.asks]),
            "\n".join([" " * 4 + str(b) for b in self.bids]),
        )

    def update(self, ws_update):
        if ws_update["u"] < self.id + 1:
            # ignore older message
            return
        if ws_update["U"] > self.id + 1:
            raise ValueError(
                "base order book ID %d falls behind update between %d-%d"
                % (self.id, ws_update["U"], ws_update["u"])
            )
        # start from the first message which satisfies U <= ob.id+1 <= u
        logger.debug("current id %d, update from %s", self.id, ws_update)
        for ask in ws_update["a"]:
            entry = OrderBookEntry(*ask)
            self.update_entry(self.asks, entry)
        for bid in ws_update["b"]:
            entry = OrderBookEntry(*bid)
            self.update_entry(self.bids, entry)
        # update local order book ID
        # check order book overlapping
        if len(self.asks) > 0 and len(self.bids) > 0:
            if self.asks[0].price <= self.bids[0].price:
                raise ValueError(
                    "price overlapping, min ask price %s not greater than max bid price %s"
                    % (self.asks[0].price, self.bids[0].price)
                )
        self.id = ws_update["u"]


class LocalOrderBook(object):
    def __init__(self, currency_pair: str):
        self.cp = currency_pair
        self.q = asyncio.Queue(maxsize=500)
        self.buf = SimpleRingBuffer(size=500)
        self.ob = OrderBook(self.cp, 0, SortedList(), SortedList())

    @property
    def id(self):
        return self.ob.id

    @property
    def asks(self):
        return self.ob.asks

    @property
    def bids(self):
        return self.ob.bids

    async def construct_base_order_book(self) -> OrderBook:
        while True:
            async with aiohttp.ClientSession() as session:
                # aiohttp does not allow boolean parameter variable
                async with session.get(
                    "https://api.gateio.ws/api/v4/pilot/order_book",
                    params={"currency_pair": self.cp, "limit": 100, "with_id": "true"},
                ) as response:
                    if response.status != 200:
                        logger.warning(
                            "failed to retrieve base order book: ",
                            await response.text(),
                        )
                        await asyncio.sleep(1)
                        continue
                    result = await response.json()
                    assert isinstance(result, dict)
                    assert result.get("id")
                    logger.debug(
                        "retrieved new base order book with id %d", result.get("id")
                    )
                    ob = OrderBook(
                        self.cp,
                        result.get("id"),
                        SortedList(
                            [OrderBookEntry(*x) for x in result.get("asks")],
                            key=lambda x: x.price,
                        ),
                        # sort bid from high to low
                        SortedList(
                            [OrderBookEntry(*x) for x in result.get("bids")],
                            key=lambda x: -x.price,
                        ),
                    )
            # use cached result to recover our local order book fast
            for b in self.buf:
                try:
                    ob.update(b)
                except ValueError as e:
                    logger.warning("failed to update: %s", e)
                    await asyncio.sleep(0.5)
                    break
            else:
                return ob

    async def run(self):
        while True:
            self.ob = await self.construct_base_order_book()
            while True:
                result = await self.q.get()
                try:
                    self.ob.update(result)
                except ValueError as e:
                    logger.error("failed to update: %s", e)
                    # reconstruct order book
                    break

    def _cache_update(self, ws_update):
        if len(self.buf) > 0:
            last_id = self.buf[-1]["u"]
            if ws_update["u"] < last_id:
                # ignore older message
                return
            if ws_update["U"] != last_id + 1:
                # update message not consecutive, reconstruct cache
                self.buf = SimpleRingBuffer(size=100)
        self.buf.append(ws_update)

    async def ws_callback(self, conn: Connection, response: WebSocketResponse):
        if response.error:
            # stop the client if error happened
            conn.close()
            raise response.error
        # ignore subscribe success response
        if "s" not in response.result or response.result.get("s") != self.cp:
            return
        result = response.result
        logger.debug("received update: %s", result)
        assert isinstance(result, dict)
        self._cache_update(result)
        await self.q.put(result)


class OrderBookFrame(Frame):
    def __init__(self, screen, order_book: LocalOrderBook):
        super(OrderBookFrame, self).__init__(
            screen, screen.height, screen.width, has_border=False, name="Order Book"
        )
        # Internal state required for doing periodic updates
        self._last_frame = 0
        self._ob = order_book
        self._level = screen.height // 2 - 1

        # Create the basic form layout...
        layout = Layout([1], fill_frame=True)
        self._header = TextBox(1, as_string=True)
        self._header.disabled = True
        self._header.custom_colour = "label"
        self._asks = MultiColumnListBox(
            screen.height // 2,
            ["<25%", "<25%", "<25%"],
            [],
            titles=["Level", "Price", "Amount"],
            name="ask_book",
            parser=AsciimaticsParser(),
        )
        self._bids = MultiColumnListBox(
            screen.height // 2,
            ["<25%", "<25%", "<25%"],
            [],
            titles=["Level", "Price", "Amount"],
            name="bid_book",
            parser=AsciimaticsParser(),
        )
        self.add_layout(layout)
        layout.add_widget(self._header)
        layout.add_widget(self._asks)
        layout.add_widget(self._bids)
        self.fix()

        # Add my own colour palette
        self.palette = defaultdict(
            lambda: (Screen.COLOUR_WHITE, Screen.A_NORMAL, Screen.COLOUR_BLACK)
        )
        for key in ["selected_focus_field", "label"]:
            self.palette[key] = (
                Screen.COLOUR_WHITE,
                Screen.A_BOLD,
                Screen.COLOUR_BLACK,
            )
        self.palette["title"] = (
            Screen.COLOUR_BLACK,
            Screen.A_NORMAL,
            Screen.COLOUR_WHITE,
        )

    def _update(self, frame_no):
        # Refresh the list view if needed
        if (
            frame_no - self._last_frame >= self.frame_update_count
            or self._last_frame == 0
        ):
            self._last_frame = frame_no

            # Create the data to go in the multi-column list
            ask_data = [
                (["#%02d" % (self._level - i), str(x.price), x.amount], i)
                for i, x in enumerate(reversed(self._ob.asks[: self._level]))
            ]
            bid_data = [
                (["#%02d" % (i + 1), str(x.price), x.amount], i)
                for i, x in enumerate(self._ob.bids[: self._level])
            ]
            self._asks.options = ask_data
            self._bids.options = bid_data
            self._header.value = "Currency Pair: {}   Time: {}".format(
                "PILOTWUKONG_USDT",
                datetime.now().astimezone().strftime("%Y-%m-%d %H:%M:%S.%f%z"),
            )

        # Now redraw as normal
        super(OrderBookFrame, self)._update(frame_no)

    @property
    def frame_update_count(self):
        # Refresh once every 0.5 seconds
        return 10


async def play_order_book(screen: Screen):
    while True:
        screen.draw_next_frame()
        await asyncio.sleep(0.05)


if __name__ == "__main__":
    logging.basicConfig(level=logging.ERROR, format="%(asctime)s: %(message)s")
    conn = Connection(Configuration(app="pilot"))
    demo_cp = "PILOTWUKONG_USDT"
    order_book = LocalOrderBook(demo_cp)
    channel = PilotOrderBookUpdateChannel(conn, order_book.ws_callback)
    channel.subscribe([demo_cp, "100ms"])

    loop = asyncio.get_event_loop()

    screen = Screen.open()
    screen.set_scenes([Scene([OrderBookFrame(screen, order_book)], -1)])
    loop.create_task(order_book.run())
    loop.create_task(conn.run())
    loop.create_task(play_order_book(screen))
    try:
        loop.run_forever()
    except KeyboardInterrupt:
        for task in asyncio.Task.all_tasks(loop):
            task.cancel()
        screen.close()
        loop.close()
