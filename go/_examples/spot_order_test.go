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

var testSpotOrderParam = &model.Order{
	CurrencyPair: "BTC_USDT",
	Amount:       "1",
	Account:      "spot",
	Iceberg:      "0",
	TimeInForce:  "gtc",
	Price:        "18000",
	Text:         "t-my-custom-id",
	Side:         "buy",
}

type spotOrderTester struct {
	svc *gate.WsService
}

func newSpotOrderTester() (*spotOrderTester, error) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	tester, err := gate.NewWsService(nil, nil, gate.NewConnConfFromOption(&gate.ConfOptions{
		// URL:          "",
		Key:          "", // required
		Secret:       "", // required
		MaxRetryConn: 5,
	}))
	if err != nil {
		return nil, err
	}

	return &spotOrderTester{tester}, nil
}

func (s *spotOrderTester) loginCallback() gate.CallBack {
	return gate.NewCallBack(func(msg *gate.UpdateMsg) {
		if msg.Data.Errs != nil {
			log.Error().Msgf("[Login] label: %s, message: %s", msg.Data.Errs.Label, msg.Data.Errs.Message)
			return
		}
		log.Info().Msgf("[Login] result: %s", msg.Data.Result)
	})
}

func (s *spotOrderTester) createOrderCallback() gate.CallBack {
	return gate.NewCallBack(func(msg *gate.UpdateMsg) {
		if msg.Data.Errs != nil {
			log.Error().Msgf("[Create] label: %s, message: %s", msg.Data.Errs.Label, msg.Data.Errs.Message)
			return
		}

		var order resp.SpotOrder
		if err := json.Unmarshal(msg.Data.Result, &order); err != nil {
			log.Error().Msgf("[Create] failed to unmarshal response: %v, msg: %v", err, msg)
			return
		}

		if order.Id == "" {
			return
		}

		log.Info().Msgf("[Create] order_id: %s, price: %v, amount: %v, status: %s", order.Id, order.Price, order.Amount, order.Status)
	})
}

func (s *spotOrderTester) orderAmendCallback() gate.CallBack {
	return gate.NewCallBack(func(msg *gate.UpdateMsg) {

		if msg.Data.Errs != nil {
			log.Error().Msgf("[Amend] label: %s, message: %s", msg.Data.Errs.Label, msg.Data.Errs.Message)
			return
		}

		var order resp.SpotOrder
		if err := json.Unmarshal(msg.Data.Result, &order); err != nil {
			log.Error().Msgf("[Amend] failed to unmarshal response: %v, msg: %v", err, msg)
			return
		}

		log.Info().Msgf("[Amend] order_id: %s, price: %v, amount: %v, status: %s", order.Id, order.Price, order.Amount, order.Status)
	})
}

func (s *spotOrderTester) orderStatusCallback() gate.CallBack {
	return gate.NewCallBack(func(msg *gate.UpdateMsg) {
		if msg.Data.Errs != nil {
			log.Error().Msgf("[Query] label: %s, message: %s", msg.Data.Errs.Label, msg.Data.Errs.Message)
			return
		}

		var order resp.SpotOrder
		if err := json.Unmarshal(msg.Data.Result, &order); err != nil {
			log.Error().Msgf("[Query] failed to unmarshal response: %v, msg: %v", err, msg)
			return
		}

		log.Info().Msgf("[Query] order_id: %s, status: %s", order.Id, order.Status)
	})
}

func (s *spotOrderTester) orderCancelCallback() gate.CallBack {
	return gate.NewCallBack(func(msg *gate.UpdateMsg) {
		if msg.Data.Errs != nil {
			log.Info().Msgf("[Cancel] failed to cancel order, label: %s, message: %s", msg.Data.Errs.Label, msg.Data.Errs.Message)
			return
		}

		orders := make([]*resp.SpotOrder, 0)
		if err := json.Unmarshal(msg.Data.Result, &orders); err != nil {
			log.Info().Msgf("[Cancel] failed to unmarshal response: %v, msg: %v", err, msg)
			return
		}

		log.Info().Msgf("[Cancel] order_id: %s, cancel succeeded: %v", orders[0].Id, orders[0].Succeeded)
	})
}

func TestSpotCreateOrder(t *testing.T) {
	s, err := newSpotOrderTester()
	assert.NoError(t, err)

	s.svc.SetCallBack(gate.ChannelSpotLogin, s.loginCallback())
	s.svc.SetCallBack(gate.ChannelSpotOrderPlace, s.createOrderCallback())

	assert.NoError(t, s.svc.APIRequest(gate.ChannelSpotOrderPlace, testSpotOrderParam, testKeyVals))

	time.Sleep(5 * time.Second)
}

func TestSpotAmendOrder(t *testing.T) {
	s, err := newSpotOrderTester()
	assert.NoError(t, err)

	s.svc.SetCallBack(gate.ChannelSpotLogin, s.loginCallback())
	s.svc.SetCallBack(gate.ChannelSpotOrderAmend, s.orderAmendCallback())

	// NOTE: Only can chose one of amount or price
	order := &model.AmendOrderParam{
		Price:        "19000",
		OrderId:      "order_id",
		CurrencyPair: "BTC_USDT",
		AmendText:    "",
	}
	assert.NoError(t, s.svc.APIRequest(gate.ChannelSpotOrderAmend, order, testKeyVals))

	time.Sleep(5 * time.Second)
}

func TestSpotQueryOrderStatus(t *testing.T) {
	s, err := newSpotOrderTester()
	assert.NoError(t, err)

	s.svc.SetCallBack(gate.ChannelSpotLogin, s.loginCallback())
	s.svc.SetCallBack(gate.ChannelSpotOrderStatus, s.orderStatusCallback())
	orderStatus := &model.StatusOrderParam{
		CurrencyPair: "BTC_USDT",
		OrderId:      "order_id",
	}

	assert.NoError(t, s.svc.APIRequest(gate.ChannelSpotOrderStatus, orderStatus, testKeyVals))

	time.Sleep(5 * time.Second)
}

func TestSpotCancelOrder(t *testing.T) {
	s, err := newSpotOrderTester()
	assert.NoError(t, err)

	s.svc.SetCallBack(gate.ChannelSpotLogin, s.loginCallback())
	s.svc.SetCallBack(gate.ChannelSpotOrderCancelIds, s.orderCancelCallback())

	testSpotOrderParam.Id = "order_id"
	assert.NoError(t, s.svc.APIRequest(gate.ChannelSpotOrderCancelIds, []*model.Order{testSpotOrderParam}, testKeyVals))

	time.Sleep(5 * time.Second)
}
