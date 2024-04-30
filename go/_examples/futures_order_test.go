package main

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	gate "github.com/gateio/gatews/go"
	"github.com/gateio/gatews/go/model"
	"github.com/gateio/gatews/go/resp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

var (
	testFuturesOrderParam = &model.FuturesOrder{
		Contract: "BTC_USDT",
		Size:     50,
		Iceberg:  0,
		Price:    "30000",
		Text:     "t-my-custom-id",
	}

	testKeyVals = map[string]any{
		"X-Gate-Channel-Id": "T-xxx",
		"req_id":            "test_req_id",
	}
)

type futuresTester struct {
	svc *gate.WsService
}

func newFuturesTester() (*futuresTester, error) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	tester, err := gate.NewWsService(nil, nil, gate.NewConnConfFromOption(&gate.ConfOptions{
		// URL:          "",
		App:          "futures", // required
		Key:          "",        // required
		Secret:       "",        // required
		MaxRetryConn: 5,
	}))
	if err != nil {
		return nil, err
	}

	return &futuresTester{tester}, nil
}

func (s *futuresTester) loginCallback() gate.CallBack {
	return gate.NewCallBack(func(msg *gate.UpdateMsg) {
		if msg.Data.Errs != nil {
			log.Error().Msgf("[Login] label: %s, message: %s", msg.Data.Errs.Label, msg.Data.Errs.Message)
			return
		}
		log.Info().Msgf("[Login] result: %s", msg.Data.Result)
	})
}

func (s *futuresTester) createOrderCallback() gate.CallBack {
	return gate.NewCallBack(func(msg *gate.UpdateMsg) {
		if msg.Data.Errs != nil {
			log.Error().Msgf("[Create] label: %s, message: %s", msg.Data.Errs.Label, msg.Data.Errs.Message)
			return
		}

		var order resp.FutureOrder
		if err := json.Unmarshal(msg.Data.Result, &order); err != nil {
			log.Error().Msgf("[Create] failed to unmarshal response: %v, msg: %v", err, msg)
			return
		}

		if order.Id == 0 {
			return
		}

		log.Info().Msgf("[Create] order_id: %v, status: %s", order.Id, order.Status)
	})
}

func (s *futuresTester) orderAmendCallback() gate.CallBack {
	return gate.NewCallBack(func(msg *gate.UpdateMsg) {
		if msg.Data.Errs != nil {
			log.Error().Msgf("[Amend] label: %s, message: %s", msg.Data.Errs.Label, msg.Data.Errs.Message)
			return
		}

		var order resp.FutureOrder
		if err := json.Unmarshal(msg.Data.Result, &order); err != nil {
			log.Error().Msgf("[Amend] failed to unmarshal response: %v, msg: %v", err, msg)
			return
		}

		log.Info().Msgf("[Amend] order_id: %v, status: %s", order.Id, order.Status)
	})
}

func (s *futuresTester) orderStatusCallback() gate.CallBack {
	return gate.NewCallBack(func(msg *gate.UpdateMsg) {
		if msg.Data.Errs != nil {
			log.Info().Msgf("[Query] label: %s, message: %s", msg.Data.Errs.Label, msg.Data.Errs.Message)
			return
		}

		var order resp.FutureOrder
		if err := json.Unmarshal(msg.Data.Result, &order); err != nil {
			log.Info().Msgf("[Query] failed to unmarshal response: %v, msg: %v", err, msg)
			return
		}

		log.Info().Msgf("[Query] order: %#v", order)
	})
}

func (s *futuresTester) orderCancelCallback() gate.CallBack {
	return gate.NewCallBack(func(msg *gate.UpdateMsg) {
		if msg.Data.Errs != nil {
			log.Error().Msgf("[Cancel] failed to cancel order, label: %s, message: %s", msg.Data.Errs.Label, msg.Data.Errs.Message)
			return
		}

		var orders *resp.FutureOrder
		if err := json.Unmarshal(msg.Data.Result, &orders); err != nil {
			log.Error().Msgf("[Cancel] failed to unmarshal response: %v, msg: %v", err, msg)
			return
		}

		log.Info().Msgf("[Cancel] order_id: %v, status: %v", orders.Id, orders.Status)
	})
}

func TestFuturesCreateOrder(t *testing.T) {
	s, err := newFuturesTester()
	assert.NoError(t, err)

	s.svc.SetCallBack(gate.ChannelFutureLogin, s.loginCallback())
	s.svc.SetCallBack(gate.ChannelFutureOrderPlace, s.createOrderCallback())
	assert.NoError(t, s.svc.APIRequest(gate.ChannelFutureOrderPlace, testFuturesOrderParam, testKeyVals))

	time.Sleep(5 * time.Second)
}

func TestFuturesAmendOrder(t *testing.T) {
	s, err := newFuturesTester()
	assert.NoError(t, err)

	s.svc.SetCallBack(gate.ChannelFutureLogin, s.loginCallback())
	s.svc.SetCallBack(gate.ChannelFutureOrderAmend, s.orderAmendCallback())

	order := &model.AmendFuturesOrder{
		OrderId:   "order_id",
		Settle:    "USDT",
		Price:     "40000",
		AmendText: "",
		Size:      1,
	}
	assert.NoError(t, s.svc.APIRequest(gate.ChannelFutureOrderAmend, order, testKeyVals))
}

func TestFuturesQueryOrderStatus(t *testing.T) {
	s, err := newFuturesTester()
	assert.NoError(t, err)

	s.svc.SetCallBack(gate.ChannelFutureLogin, s.loginCallback())
	s.svc.SetCallBack(gate.ChannelFutureOrderStatus, s.orderStatusCallback())

	order := &model.StatusFuturesOrder{
		Settle:  "usdt",
		OrderId: "order_id",
	}

	assert.NoError(t, s.svc.APIRequest(gate.ChannelFutureOrderStatus, order, testKeyVals))

	time.Sleep(5 * time.Second)
}

func TestFuturesCancelOrder(t *testing.T) {
	s, err := newFuturesTester()
	assert.NoError(t, err)

	s.svc.SetCallBack(gate.ChannelFutureLogin, s.loginCallback())
	s.svc.SetCallBack(gate.ChannelFutureOrderCancel, s.orderCancelCallback())

	order := &model.CancelFuturesOrder{
		Settle:  "usdt",
		OrderId: "order_id",
	}
	assert.NoError(t, s.svc.APIRequest(gate.ChannelFutureOrderCancel, order, testKeyVals))

	time.Sleep(5 * time.Second)
}
