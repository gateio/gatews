# Gate.io WebSocket Go SDK

`gatews` provides Gate.io WebSocket V4 Go implementation, including  all channels defined in spot(new) WebSocket.

Features:

1. Fully asynchronous
2. Reconnect on connection to server lost, and resubscribe on connection recovered
3. Support connecting to multiple websocket servers
4. Highly configurable
## Installation
```shell
go get github.com/gateio/gatews/go
```
## Getting started
```golang
package main

import (
	"encoding/json"
	gate "github.com/gateio/gatews/go"
	"log"
	"time"
)

func main() {
	// create WsService with ConnConf, this is recommended, key and secret will be needed by some channels
	// ctx and logger could be nil, they'll be initialized by default
	ws, err := gate.NewWsService(nil, nil, gate.NewConnConf("",
		"YOUR_API_KEY", "YOUR_API_SECRET", 10))
	// we can also do nothing to get a WsService, all parameters will be initialized by default
	// but some channels need key and secret for auth, we can also use set function to set key and secret
	//ws, err := gate.NewWsService(nil, nil, nil)
	//ws.SetKey("YOUR_API_KEY")
	//ws.SetSecret("YOUR_API_SECRET")
	if err != nil {
		log.Printf("NewWsService err:%s", err.Error())
		return
	}

	// create callback functions for receive messages
	callOrder := gate.NewCallBack(func(msg *gate.UpdateMsg) {
		// parse the message to struct we need
		var order []gate.SpotOrderMsg
		if err := json.Unmarshal(msg.Result, &order); err != nil {
			log.Printf("order Unmarshal err:%s", err.Error())
		}
		log.Printf("%+v", order)
	})
	callTrade := gate.NewCallBack(func(msg *gate.UpdateMsg) {
		var trade gate.SpotTradeMsg
		if err := json.Unmarshal(msg.Result, &trade); err != nil {
			log.Printf("trade Unmarshal err:%s", err.Error())
		}
		log.Printf("%+v", trade)
	})
	// first, we need set callback function
	ws.SetCallBack(gate.ChannelOrder, callOrder)
	ws.SetCallBack(gate.ChannelPublicTrade, callTrade)
	// second, after set callback function, subscribe to any channel you are interested into
	if err := ws.Subscribe(gate.ChannelPublicTrade, []string{"BCH_USDT"}); err != nil {
		log.Printf("Subscribe err:%s", err.Error())
		return
	}
	if err := ws.Subscribe(gate.ChannelOrder, []string{"BCH_USDT"}); err != nil {
		log.Printf("Subscribe err:%s", err.Error())
		return
	}

	// example for maintaining local order book
	OrderBookExample(ws)

	ch := make(chan bool)
	defer close(ch)

	for {
		select {
		case <-ch:
			log.Printf("manual done")
		case <-time.After(time.Second * 1000):
			log.Printf("auto done")
			return
		}
	}
}
```
## Example
We provide some demo applications in the [examples](_examples) directory, which can be run directly.


