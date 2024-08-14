# Changelog

## v0.5.1

2024-08-14

- added recent new fields

## v0.5.0

2024-04-30

- support order operations

## v0.4.2

2023-05-09

- no longer try websocket ping after disconnect
- support querying websocket connection status

## v0.4.1

2023-03-16

- fix subscribe msg concurrent write cause panic
- add `FuturesOrder` response field `stop_profit_price` and `stop_loss_price`
- add `FuturesAutoOrder` response field `me_order_id`, `order_type` and `initial.auto_size`

## v0.4.0

2022-12-07

- spot balance add fields `freeze`,`freeze_change` and `change_type`

## v0.3.0

2022-11-22

- add common response field `time_ms` for time of message created

## v0.2.8

2022-10-21

- futures model `FuturesUserTrade` add fields `fee` and `point_fee`
- avoid saving `ping` and `time` subscribe msg in request history

## v0.2.7

2022-08-11

- remove client method `NewConnConf`. Recommend to use `NewConnConfFromOption` instead
- add new config field `PingInterval` to send ping message
- add default ping message to avoid to be closed by server

## v0.2.6

2022-05-24

- fix reconnect panic
- update spot and futures models
- add `ShowReconnectMsg` config field to decide to show reconnect success msg
- update some test cases

## v0.2.5

- update future's models, fix fields wrong type

## v0.2.4

2021-08-12

- update futures models, fix fields wrong type
- add futures model `FuturesPositions`

## v0.2.3

2021-08-11

- Support websocket skip tls verify with `SkipTlsVerify` of ConfOptions
- Update futures models
- Update examples

## v0.2.2

2021-07-23

- Add constant `ChannelSpotCrossBalance` support `spot.cross_balances` channel
- SpotUserTradesMsg add field `Text` for orders' text
- Update local order book example

## v0.2.1

2021-06-24

- fix futures book ticker and order book update struct

## v0.2.0

2021-06-24

- fix futures order book update struct

## v0.1.9

2021-06-24

- add futures order book struct

## v0.1.8

2021-06-22

- add `WsService` method `GetConnection()` to get the connection
- fix `changelog` date error

## v0.1.7

2021-06-04

- fix reconnect msg nil `SubscribeOptions` caused reconnect msg lost

## v0.1.6

2021-06-04

- add `io.ErrUnexpectedEOF` error capture, it caused v0.1.5 can't reconnect
- fix reconnect msg repeat add

## v0.1.5

2021-06-02

- add `SpotUpdateAllDepthMsg` struct for parse all order book msg

## v0.1.4

2021-06-02

- fix overlapping price for local order book example
- update README

## v0.1.3

2021-05-26

- Support futures websocket.
- Modify channels name to with flag `Spot` or `Future`.
- Add field `TimestampInMilli` in models `SpotBalancesMsg`, `SpotFundingBalancesMsg`, `SpotMarginBalancesMsg`. Add
  field `TimeInMilli` in model `SpotBookTickerMsg`
- Add new method `NewConnConfFromOption` to get a ConnConf flexible.
- Add new method `SubscribeWithOption` to support futures subscribe with id.
- Add example for both spot and futures connection use.
- Fix reconnect websocket failed bug.
- Optimizing code structure.

## v0.1.2

2021-04-19

- Fix subscribe repeat bug.

## v0.1.1

2021-04-16

- Fix subscribe channel failed bug.

## v0.1.0

2021-04-12

- Support spot websocket function.
