## example for order book

### How to maintain local depth
1. Subscribe `spot.order_book_update` with specified level and update frequency, e.g. `["BTC_USDT", "1000ms"]` pushes the update in BTC_USDT order book every 1s
2. Cache WebSocket notifications. Every notification use `U` and `u` to tell the first and last update ID since last notification.
3. Retrieve base order book using REST API, and make sure the order book ID is recorded(referred as `baseID` below) e.g. `https://api.gateio.ws/api/v4/spot/order_book?currency_pair=BTC_USDT&limit=100&with_id=true` retrieves the base order book of BTC_USDT
4. Iterate the cached WebSocket notifications, and find the first one which contains the baseID, i.e. `U <= baseId+1` and `u >= baseId+1`, then start consuming from it. Note that sizes in notifications are all absolute values. Use them to replace original sizes in corresponding price. If size equals to 0, delete the price from the order book.
5. Dump all notifications which satisfy `u < baseID+1`. If `baseID+1 < first notification U`, it means current base order book falls behind notifications. Start from step 3 to retrieve newer base order book.
6. If any subsequent notification which satisfy `U > baseID+1` is found, it means some updates are lost. Reconstruct local order book from step 3.

