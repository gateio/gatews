# !/usr/bin/env python
# coding: utf-8

"""
Example of how to order spot
"""
import asyncio
import logging

from gate_ws import Configuration, Connection, WebSocketResponse
from gate_ws.spot import SpotOrderCancelChannel, SpotOrderPlaceChannel

logger = logging.getLogger(__name__)


async def callback(_: Connection, response: WebSocketResponse):
    if response.error:
        logger.error("failed to api_request: %s", response.error)

    if response.channel == "spot.login":
        return

    if response.channel == "spot.order_place":
        logger.info("order status: %s", response.result)

    if response.channel == "spot.order_cancel":
        logger.info("order cancel: %s", response.result)


order_place_param = {
    "text": "t-sssd",
    "currency_pair": "BCH_USDT",
    "type": "limit",
    "account": "spot",
    "side": "buy",
    "iceberg": "0",
    "price": "20",
    "amount": "0.05",
    "time_in_force": "gtc",
    "auto_borrow": False,
}

order_cancel_param = {"currency_pair": "BCH_USDT", "order_id": "1862000415"}

if __name__ == "__main__":
    logging.basicConfig(level=logging.DEBUG, format="%(asctime)s: %(message)s")
    cfg = Configuration(
        api_key="{API_KEY}", # required
        api_secret="{API_SECRET}", # required
    )

    conn = Connection(cfg)
    SpotOrderPlaceChannel(conn, callback).api_request(
        order_place_param, "header", "req_id"
    )

    SpotOrderCancelChannel(conn, callback).api_request(
        order_cancel_param, "header", "req_id"
    )

    loop = asyncio.get_event_loop()
    loop.create_task(conn.run())

    try:
        loop.run_forever()
    except KeyboardInterrupt:
        tasks = asyncio.Task.all_tasks(loop)
        for task in tasks:
            task.cancel()
        group = asyncio.gather(*tasks, return_exceptions=True)
        loop.run_until_complete(group)
        loop.close()
