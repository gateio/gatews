# Demo applications

## Local order book

[local_order_book.py](local_order_book.py) provides a demo application showing how to maintain a
local spot BTC_USDT order book using Gate.io WebSocket `spot.order_book_update` channel and HTTP
API.

To run this demo, you need to install the following packages manually:

```sh
pip install --user gate-ws sortedcontainers aiohttp asciimatics 
```

Then run it directly `python local_order_book.py`

> Python3.6+ is required.

This application maintains the local order book through `LocalOrderBook` class which provides a
callback method `ws_callback` for WebSocket connection to call on order book update received. The
animation is shown through `OrderBookFrame` using `asciimatics`. You can use `LocalOrderBook`
instance in your own application, without `OrderBookFrame`, `screen` initialization
and `play_order_book` task.

Some notes:

1. When the demo application starts, you might see the order book did not change for several
   seconds. It is a normal case, as the HTTP result is trying to keep up with WebSocket updates'
   pace.
2. You can only run this application in a terminal. Order book maintained this way provides at most
   100 levels. The application will detect the terminal height to display with proper levels. The
   longer your terminal, the higher levels will be shown.
3. Resizing your terminal is not supported when running.

Here is a pre-recorded clip.

[![asciicast](https://asciinema.org/a/406646.svg)](https://asciinema.org/a/406646)
