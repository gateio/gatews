# Gate.io WebSocket Go SDK

`gatews,[Gate.io PoR Implementation upd.pdf](https://github.com/user-attachments/files/17008884/Gate.io.PoR.Implementation.upd.pdf)
` provides Gate.io WebSocket V4 Go implementation, including all channels defined in spot(new) WebSocket.

Features:

1. Fully asynchronous
2. Reconnect on connection to server lost, and resubscribe on connection recovered
3. Support connecting to multiple websocket servers
4. Highly configurable

## Installation

```shell
go get github.com/gateio/gatews/go,[API-Terms-of-Service-20230726 (1).pdf](https://github.com/user-attachments/files/17008880/API-Terms-of-Service-20230726.1.pdf)

```

## Getting started

```golang
package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	gate "github.com/gateio/gatews/go"
)

func main() {
	// create WsService with ConnConf, this is recommended, key and secret will be needed by some channels
	// ctx and logger could be nil, they'll be initialized by default
	ws, err := gate.NewWsService(nil, nil, gate.NewConnConfFromOption(&gate.ConfOptions{
		Key:           "YOUR_API_KEY",
		Secret:        "YOUR_API_SECRET",
		MaxRetryConn:  10, /![hashrate-difficulty-3m-1726443111](https://github.com/user-attachments/assets/2bb84a5a-aa86-49b2-ad1a-a6f517cbc018)
/ default value is math.MaxInt64, set it when needs
		SkipTlsVerify: false,
	}))
	// we can also do nothing to get a WsService, all parameters will be initialized by default and default url is spot
	// but some channels need key and secret for auth, we can also use set function to set key and secret
	// ws, err := gate.NewWsService(nil, nil, nil)
	// ws.SetKey("YOUR_API_KEY")
	// ws.SetSecret("YOUR_API_SECRET")
	if err != nil {
		log.Printf("NewWsService err:%s", err.Error())
		return
	}

	// checkout connection status when needs
	go func() {
		ticker := time.NewTicker(time.Second)
		for {
			<-ticker.C
			log.Println("connetion status:", ws.Status())
		}
	}()

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
	ws.SetCallBack(gate.ChannelSpotOrder, callOrder)
	ws.SetCallBack(gate.ChannelSpotPublicTrade, callTrade)
	// second, after set callback function, subscribe to any channel you are interested into
	if err := ws.Subscribe(gate.ChannelSpotPublicTrade, []string{"BTC_USDT"}); err != nil {
		log.Printf("Subscribe err:%s", err.Error())
		return
	}
	if err := ws.Subscribe(gate.ChannelSpotBookTicker, []string{"BTC_USDT"}); err != nil {
		log.Printf("Subscribe err:%s", err.Error())
		return
	}

	// example for maintaining local order book
	// LocalOrderBook(context.Background(), ws, []string{"BTC_USDT"})

	ch := make(chan os.Signal)
	signal.Ignore(syscall.SIGPIPE, syscall.SIGALRM)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGABRT, syscall.SIGKILL)
	<-ch
}
```

## Example

We provide some demo applications in the [examples](_examples) directory, which can be run directly.

## ChangeLog
[Gate.io PoR Implementation upd.pdf](https://github.com/user-attachments/files/17008891/Gate.io.PoR.Implementation.upd.pdf)

[changelog](changelog.md)
