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

	if err := ws.Subscribe(ChannelSpotPublicTrade, []string{"BCH_USDT"}); err != nil {
		log.Fatalf("Subscribe err:%s", err.Error())
		return
	}
	if err := ws.Subscribe(ChannelSpotPublicTrade, []string{"BTC_USDT"}); err != nil {
		log.Fatalf("Subscribe err:%s", err.Error())
		return
	}
	if err := ws.Subscribe(ChannelSpotOrderBookUpdate, []string{"BTC_USDT", "100ms"}); err != nil {
		log.Fatalf("Subscribe err:%s", err.Error())
		return
	}
	if err := ws.Subscribe(ChannelSpotOrderBookUpdate, []string{"ETH_USDT", "100ms"}); err != nil {
		log.Fatalf("Subscribe err:%s", err.Error())
		return
	}
	fmt.Println(ChannelSpotPublicTrade, " subscribed markets: ", ws.GetChannelMarkets(ChannelSpotPublicTrade))
	fmt.Println(ChannelSpotOrderBookUpdate, " subscribed markets: ", ws.GetChannelMarkets(ChannelSpotOrderBookUpdate))

	if err := ws.UnSubscribe(ChannelSpotPublicTrade, []string{"BTC_USDT"}); err != nil {
		log.Fatalf("Subscribe err:%s", err.Error())
		return
	}
	if err := ws.UnSubscribe(ChannelSpotOrderBookUpdate, []string{"BTC_USDT", "100ms"}); err != nil {
		log.Fatalf("Subscribe err:%s", err.Error())
		return
	}
	fmt.Println("after unsubscribe")
	fmt.Println(ChannelSpotPublicTrade, " subscribed markets: ", ws.GetChannelMarkets(ChannelSpotPublicTrade))
	fmt.Println(ChannelSpotOrderBookUpdate, " subscribed markets: ", ws.GetChannelMarkets(ChannelSpotOrderBookUpdate))
}

func TestGetChannels(t *testing.T) {
	ws, err := NewWsService(nil, nil, nil)
	if err != nil {
		log.Fatal(err)
	}

	call := NewCallBack(func(msg *UpdateMsg) {})
	ws.SetCallBack(ChannelSpotPublicTrade, call)
	if err := ws.Subscribe(ChannelSpotPublicTrade, []string{"BCH_USDT"}); err != nil {
		log.Fatalf("Subscribe err:%s", err.Error())
		return
	}
	if err := ws.Subscribe(ChannelSpotCandleStick, []string{"BTC_USDT", "10ms"}); err != nil {
		log.Fatalf("Subscribe err:%s", err.Error())
		return
	}

	fmt.Println(ws.GetChannels())
}

func TestGetConf(t *testing.T) {
	ws, err := NewWsService(nil, nil, NewConnConf(
		"", "KEY", "SECRET", 10))
	if err != nil {
		log.Fatal(err)
	}

	call := NewCallBack(func(msg *UpdateMsg) {})
	ws.SetCallBack(ChannelSpotPublicTrade, call)
	if err := ws.Subscribe(ChannelSpotPublicTrade, []string{"BCH_USDT"}); err != nil {
		log.Fatalf("Subscribe err:%s", err.Error())
		return
	}

	fmt.Println(ws.GetKey())
	fmt.Println(ws.GetSecret())
	fmt.Println(ws.GetMaxRetryConn())
}

func TestGetConfFromOption(t *testing.T) {
	ws, err := NewWsService(nil, nil, NewConnConfFromOption(&ConfOptions{
		"", "KEY", "SECRET", 10,
	}))
	if err != nil {
		log.Fatal(err)
	}

	call := NewCallBack(func(msg *UpdateMsg) {})
	ws.SetCallBack(ChannelSpotPublicTrade, call)
	if err := ws.Subscribe(ChannelSpotPublicTrade, []string{"BCH_USDT"}); err != nil {
		log.Fatalf("Subscribe err:%s", err.Error())
		return
	}
	fmt.Println(ws.GetKey())
	fmt.Println(ws.GetSecret())
	fmt.Println(ws.GetMaxRetryConn())
}
