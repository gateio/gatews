package gatews

import (
	"fmt"
	"log"
	"testing"
)

func TestGetChannelMarkets(t *testing.T) {
	ws, err := NewWsService(nil, nil, nil)
	if err != nil {
		log.Fatal(err)
	}

	if err := ws.Subscribe(ChannelPublicTrade, []string{"BCH_USDT"}); err != nil {
		log.Fatalf("Subscribe err:%s", err.Error())
		return
	}
	if err := ws.Subscribe(ChannelPublicTrade, []string{"BTC_USDT"}); err != nil {
		log.Fatalf("Subscribe err:%s", err.Error())
		return
	}
	fmt.Println(ws.GetChannelMarkets(ChannelPublicTrade))

	if err := ws.UnSubscribe(ChannelPublicTrade, []string{"BTC_USDT"}); err != nil {
		log.Fatalf("Subscribe err:%s", err.Error())
		return
	}
	fmt.Println(ws.GetChannelMarkets(ChannelPublicTrade))
}

func TestGetChannels(t *testing.T) {
	ws, err := NewWsService(nil, nil, nil)
	if err != nil {
		log.Fatal(err)
	}

	call := NewCallBack(func(msg *UpdateMsg) {})
	ws.SetCallBack(ChannelPublicTrade, call)
	if err := ws.Subscribe(ChannelPublicTrade, []string{"BCH_USDT"}); err != nil {
		log.Fatalf("Subscribe err:%s", err.Error())
		return
	}
	if err := ws.Subscribe(ChannelCandleStick, []string{"BTC_USDT", "10ms"}); err != nil {
		log.Fatalf("Subscribe err:%s", err.Error())
		return
	}

	fmt.Println(ws.GetChannels())
}

func TestGetConf(t *testing.T) {
	ws, err := NewWsService(nil, nil, NewConnConf("", "eqywieyqw", "sdsadsad", 10))
	if err != nil {
		log.Fatal(err)
	}

	call := NewCallBack(func(msg *UpdateMsg) {})
	ws.SetCallBack(ChannelPublicTrade, call)
	if err := ws.Subscribe(ChannelPublicTrade, []string{"BCH_USDT"}); err != nil {
		log.Fatalf("Subscribe err:%s", err.Error())
		return
	}

	fmt.Println(ws.GetKey())
	fmt.Println(ws.GetSecret())
	fmt.Println(ws.GetMaxRetryConn())
}
