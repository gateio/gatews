# !/usr/bin/env python
# coding: utf-8

"""
Example of subscribe tickers
"""
import asyncio
import logging

from gate_ws import Configuration, Connection, WebSocketResponse
from gate_ws.spot import SpotTickerChannel

logger = logging.getLogger(__name__)


async def callback(conn: Connection, response: WebSocketResponse):
    if response.error:
        conn.close()
        raise response.error

    result = response.result
    logger.debug("received update: %s", result)
    assert isinstance(result, dict)


if __name__ == "__main__":
    logging.basicConfig(level=logging.DEBUG, format="%(asctime)s: %(message)s")
    loop = asyncio.new_event_loop()
    asyncio.set_event_loop(loop)

    cfg = Configuration(event_loop=loop)

    conn = Connection(cfg)
    channel = SpotTickerChannel(conn, callback)
    channel.subscribe(["BTC_USDT"])

    tasks: set[asyncio.Task] = {
        loop.create_task(conn.run()),
    }

    try:
        loop.run_forever()
    except KeyboardInterrupt:
        for task in tasks:
            task.cancel()

        group = asyncio.gather(*tasks, return_exceptions=True)
        loop.run_until_complete(group)
    finally:
        loop.close()
