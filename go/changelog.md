# Changelog

## v0.1.4

2020-06-02

- fix overlapping price for local order book example
- update README

## v0.1.3

2020-05-26

- Support futures websocket.
- Modify channels name to with flag `Spot` or `Future`.
- Add field `TimestampInMilli`  in models `SpotBalancesMsg`, `SpotFundingBalancesMsg`, `SpotMarginBalancesMsg`. Add
  field `TimeInMilli` in model `SpotBookTickerMsg`
- Add new method `NewConnConfFromOption` to get a ConnConf flexible.
- Add new method `SubscribeWithOption` to support futures subscribe with id.
- Add example for both spot and futures connection use.
- Fix reconnect websocket failed bug.
- Optimizing code structure.

## v0.1.2

2020-04-19

- Fix subscribe repeat bug.

## v0.1.1

2020-04-16

- Fix subscribe channel failed bug.

## v0.1.0

2020-04-12

- Support spot websocket function.