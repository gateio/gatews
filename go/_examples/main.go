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
		MaxRetryConn:  10, // default value is math.MaxInt64, set it when needs
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
		if msg.Event != "update" {
			return
		}
		// parse the message to struct we need
		var order []gate.SpotOrderMsg
		if err := json.Unmarshal(msg.Result, &order); err != nil {
			log.Printf("order %s unmarshal err: %v", msg.Result, err)
		}
		log.Printf("order: %+v", order)
	})

	callTrade := gate.NewCallBack(func(msg *gate.UpdateMsg) {
		var trade gate.SpotTradeMsg
		if err := json.Unmarshal(msg.Result, &trade); err != nil {
			log.Printf("trade %s unmarshal err: %v", msg.Result, err)
		}
		log.Printf("trade: %+v", trade)
	})

	// first, we need set callback function
	ws.SetCallBack(gate.ChannelSpotOrder, callOrder)
	ws.SetCallBack(gate.ChannelSpotPublicTrade, callTrade)
	// second, after set callback function, subscribe to any channel you are interested into
	if err := ws.Subscribe(gate.ChannelSpotOrder, []string{"BTC_USDT"}); err != nil {
		log.Printf("Subscribe err:%s", err.Error())
		return
	}
	if err := ws.Subscribe(gate.ChannelSpotPublicTrade, []string{"BTC_USDT"}); err != nil {
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
