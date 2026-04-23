# !/usr/bin/env python
# coding: utf-8

"""
Example of how to track WebSocket latency with local timestamps.
"""

import time
import logging
import asyncio
from gate_ws import Configuration, Connection, WebSocketResponse
from gate_ws.spot import SpotPublicTradeChannel

logger = logging.getLogger(__name__)


async def on_trade(conn: Connection, response: WebSocketResponse):
    # Local timestamp is injected into the 'result' dictionary
    data = response.result

    assert '_local_ts' in data, "Local timestamp not found in data"
    local_ts = data['_local_ts']
    
    # Immitation of local processing
    await asyncio.sleep(0.05) 

    now_ns = int(time.time() * 1_000_000_000)
    latency_ns = now_ns - local_ts
    latency_ms = latency_ns / 1_000_000
    logger.info(f"[TRADE] Price: {data.get('price')}, Latency: {latency_ms:.3f} ms")



if __name__ == "__main__":
    logging.basicConfig(level=logging.DEBUG, format="%(asctime)s: %(message)s")
    loop = asyncio.new_event_loop()
    asyncio.set_event_loop(loop)

    cfg = Configuration(
        event_loop=loop,
        add_local_ts=True # Enable local timestamp feature
    )

    conn = Connection(cfg)
    channel = SpotPublicTradeChannel(conn, on_trade)
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