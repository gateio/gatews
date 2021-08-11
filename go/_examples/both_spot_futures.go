package main

import (
	"encoding/json"
	"log"
	"time"

	gate "github.com/gateio/gatews/go"
)

func main2() {
	spotWs, err := gate.NewWsService(nil, nil, gate.NewConnConfFromOption(&gate.ConfOptions{
		URL:          gate.BaseUrl,
		Key:          "YOUR_API_KEY",
		Secret:       "YOUR_API_SECRET",
		MaxRetryConn: 10,
	}))
	if err != nil {
		log.Printf("new spot wsService err:%s", err.Error())
		return
	}

	futureWs, err := gate.NewWsService(nil, nil, gate.NewConnConfFromOption(&gate.ConfOptions{
		URL:          gate.FuturesUsdtUrl,
		Key:          "YOUR_API_KEY",
		Secret:       "YOUR_API_SECRET",
		MaxRetryConn: 10,
	}))
	if err != nil {
		log.Printf("new futures wsService err:%s", err.Error())
		return
	}

	// create callback functions for receive messages
	// spot order book update
	callOrderBookUpdate := gate.NewCallBack(func(msg *gate.UpdateMsg) {
		// parse the message to struct we need
		var update gate.SpotUpdateDepthMsg
		if err := json.Unmarshal(msg.Result, &update); err != nil {
			log.Printf("order book update Unmarshal err:%s", err.Error())
		}
		log.Printf("%+v", update)
	})

	// futures trade
	callTrade := gate.NewCallBack(func(msg *gate.UpdateMsg) {
		var trades []gate.FuturesTrade
		if err := json.Unmarshal(msg.Result, &trades); err != nil {
			log.Printf("trade Unmarshal err:%s", err.Error())
		}
		log.Printf("%+v", trades)
	})

	// first, set callback
	spotWs.SetCallBack(gate.ChannelSpotOrderBookUpdate, callOrderBookUpdate)
	futureWs.SetCallBack(gate.ChannelFutureTrade, callTrade)
	if err := spotWs.Subscribe(gate.ChannelSpotOrderBookUpdate, []string{"BTC_USDT", "100ms"}); err != nil {
		log.Printf("spotWs Subscribe err:%s", err.Error())
		return
	}

	if err := futureWs.Subscribe(gate.ChannelFutureTrade, []string{"BTC_USDT"}); err != nil {
		log.Printf("futureWs Subscribe err:%s", err.Error())
		return
	}

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
