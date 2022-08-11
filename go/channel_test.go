package gatews

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func TestNilCallBack(t *testing.T) {
	ws, err := NewWsService(nil, nil, nil)
	if err != nil {
		log.Fatal(err)
	}

	ws.SetCallBack(ChannelSpotPublicTrade, nil)
	if err := ws.Subscribe(ChannelSpotPublicTrade, []string{"BCH_USDT"}); err != nil {
		log.Fatalf("Subscribe err:%s", err.Error())
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

func TestSubscribeFutures(t *testing.T) {
	ws, err := NewWsService(nil, nil, NewConnConfFromOption(&ConfOptions{
		URL: FuturesUsdtUrl, Key: "", Secret: "", MaxRetryConn: 10,
	}))
	if err != nil {
		log.Fatal(err)
	}

	call := NewCallBack(func(msg *UpdateMsg) {
		fmt.Println(string(msg.Result))
	})
	ws.SetCallBack(ChannelFutureCandleStick, call)
	if err := ws.Subscribe(ChannelFutureCandleStick, []string{"1m", "BTC_USDT"}); err != nil {
		log.Fatalf("Subscribe err:%s", err.Error())
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

func TestSubscribeFuturesWithOptions(t *testing.T) {
	ws, err := NewWsService(nil, nil, NewConnConfFromOption(&ConfOptions{
		URL: FuturesUsdtUrl,
	}))
	if err != nil {
		log.Fatal(err)
	}

	call := NewCallBack(func(msg *UpdateMsg) {
		fmt.Printf("%+v\n", msg)
	})
	ws.SetCallBack(ChannelFutureTrade, call)
	if err := ws.SubscribeWithOption(ChannelFutureTrade, []string{"BTC_USDT"}, &SubscribeOptions{
		ID: 123456,
	}); err != nil {
		log.Fatalf("Subscribe err:%s", err.Error())
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

func TestSubscribeAuthChannel(t *testing.T) {
	ws, err := NewWsService(nil, nil, nil)
	if err != nil {
		log.Fatal(err)
	}

	if err := ws.Subscribe(ChannelSpotOrder, []string{"BCH_USDT"}); err != nil {
		log.Fatalf("Subscribe err:%s", err.Error())
		return
	}
}
