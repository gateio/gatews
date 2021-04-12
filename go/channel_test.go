package gatews

import (
	"log"
	"testing"
	"time"
)

func TestNilCallBack(t *testing.T) {
	ws, err := NewWsService(nil, nil, nil)
	if err != nil {
		log.Fatal(err)
	}

	ws.SetCallBack(ChannelPublicTrade, nil)
	if err := ws.Subscribe(ChannelPublicTrade, []string{"BCH_USDT"}); err != nil {
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

	if err := ws.Subscribe(ChannelOrder, []string{"BCH_USDT"}); err != nil {
		log.Fatalf("Subscribe err:%s", err.Error())
		return
	}
}
