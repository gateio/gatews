# Gate.io WebSocket Python SDK

`gate_ws` provides Gate.io WebSocket V4 Python implementation, including all channels defined in
spot(new) and futures WebSocket.

Features:

1. Fully asynchronous
2. Reconnect on connection to server lost, and resubscribe on connection recovered
3. Support connecting to multiple websocket servers
4. Highly configurable

## Installation

This package requires Python version 3.6+. Python 2 will NOT be supported.

```sh
pip install --user gate-ws
```

## Getting Started

```python
import asyncio

from gate_ws import Configuration, Connection, WebSocketResponse
from gate_ws.spot import SpotPublicTradeChannel


# define your callback function on message received
def print_message(conn: Connection, response: WebSocketResponse):
    if response.error:
        print('error returned: ', response.error)
        conn.close()
        return
    print(response.result)


async def main():
    # initialize default connection, which connects to spot WebSocket V4
    # it is recommended to use one conn to initialize multiple channels
    conn = Connection(Configuration())

    # subscribe to any channel you are interested into, with the callback function
    channel = SpotPublicTradeChannel(conn, print_message)
    channel.subscribe(["GT_USDT"])

    # start the client
    await conn.run()


if __name__ == '__main__':
    asyncio.run(main())
```

## Application Demos

We provide some demo applications in the [examples](examples) directory, which can be run directly.

## Advanced usage

1. Subscribe to private channels
   ```python
   from gate_ws import Configuration, Connection
   from gate_ws.spot import SpotOrderChannel


   async def main():
       conn = Connection(Configuration(api_key='YOUR_API_KEY', api_secret='YOUR_API_SECRET'))
       channel = SpotOrderChannel(conn, lambda c, r: print(r.result))
       channel.subscribe(["GT_USDT"])

       # start the client
       await conn.run()
   ```
2. Your callback function can also be a coroutine
   ```python
   import asyncio


   async def my_callback(conn, response):
       await asyncio.sleep(1)
       print(response.result)
   ```
3. You can provide a default callback function for all channels, so that when subscribing to new
   channels, no additional callback function are needed.
   ```python
   from gate_ws import Configuration, Connection
   from gate_ws.spot import SpotPublicTradeChannel


   async def main():
       # provide default callback for all channels
       conn = Connection(Configuration(default_callback=lambda c, r: print(r.result)))

       # default callback will be used if callback not provided when initializing channels
       channel = SpotPublicTradeChannel(conn)
       channel.subscribe(["GT_USDT"])

       # start the client
       await conn.run()
   ```
4. Subscribe to both spot and futures WebSockets
   ```python
   import asyncio

   from gate_ws import Configuration, Connection
   from gate_ws.spot import SpotPublicTradeChannel
   from gate_ws.futures import FuturesPublicTradeChannel


   async def main():
       # initialize a spot connection, which is the default if no parameters is provided
       spot_conn = Connection(Configuration(app='spot'))
       # initialize a futures connection
       futures_conn = Connection(Configuration(app='futures', settle='usdt', test_net=False))

       # subscribe to any channel you are interested into, with the callback function
       channel = SpotPublicTradeChannel(spot_conn, lambda c, r: print(r.result))
       channel.subscribe(["BTC_USDT"])

       channel = FuturesPublicTradeChannel(futures_conn, lambda c, r: print(r.result))
       channel.subscribe(["BTC_USDT"])

       # start both connection
       await asyncio.gather(spot_conn.run(), futures_conn.run())
   ```
5. You can use your own executor pool to run your callback function
   ```python
   import concurrent.futures

   from gate_ws import Configuration, Connection
   from gate_ws.spot import SpotPublicTradeChannel


   async def main():
       # use process pool to run your callback function
       with concurrent.futures.ProcessPoolExecutor() as pool:
           conn = Connection(Configuration(executor_pool=pool))

           # default callback will be used if callback not provided when subscribing
           channel = SpotPublicTradeChannel(conn, lambda c, r: print(r.result))
           channel.subscribe(["GT_USDT"])

           # start the client
           await conn.run()
   ```
