package main

import (
	"encoding/json"
	"os"
	"testing"

	gate "github.com/gateio/gatews/go"
	"github.com/gateio/gatews/go/model"
	"github.com/gateio/gatews/go/resp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"
)

type SpotOrderSuite struct {
	suite.Suite

	svc   *gate.WsService
	order *model.Order

	finishCreate chan struct{}
	finishAmend  chan struct{}
	finishQuery  chan struct{}
	exit         chan struct{}
}

func (s *SpotOrderSuite) SetupSuite() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	var err error
	s.svc, err = gate.NewWsService(nil, nil, gate.NewConnConfFromOption(&gate.ConfOptions{
		URL:          "",
		Key:          "",
		Secret:       "",
		MaxRetryConn: 10,
	}))
	if err != nil {
		s.T().Fatal(err)
	}

	s.finishCreate = make(chan struct{}, 1)
	s.finishAmend = make(chan struct{}, 1)
	s.finishQuery = make(chan struct{}, 1)
	s.exit = make(chan struct{}, 1)

	s.svc.SetCallBack(gate.ChannelSpotLogin, s.loginCallback())
	s.svc.SetCallBack(gate.ChannelSpotOrderPlace, s.orderCallback())
	s.svc.SetCallBack(gate.ChannelSpotOrderAmend, s.orderAmendCallback())
	s.svc.SetCallBack(gate.ChannelSpotOrderStatus, s.orderStatusCallback())
	s.svc.SetCallBack(gate.ChannelSpotOrderCancelIds, s.orderCancelCallback())

	s.order = &model.Order{
		CurrencyPair: "BTC_USDT",
		Amount:       "1",
		Account:      "spot",
		Iceberg:      "0",
		TimeInForce:  "gtc",
		Price:        "18000",
		Text:         "t-my-custom-id",
		Side:         "buy",
	}
}

func (s *SpotOrderSuite) TestOrder() {
	s.createOrder()
	s.amendOrder()
	s.queryOrderStatus()
	s.cancelOrder()

	<-s.exit
}

func (s *SpotOrderSuite) createOrder() {
	s.NoError(s.svc.APIRequest(gate.ChannelSpotOrderPlace, s.order))
}

func (s *SpotOrderSuite) amendOrder() {
	<-s.finishCreate

	// NOTE: Only can chose one of amount or price
	order := &model.AmendOrderParam{
		Price:        "19000",
		OrderId:      s.order.Id,
		CurrencyPair: s.order.CurrencyPair,
		AmendText:    s.order.AmendText,
	}
	s.NoError(s.svc.APIRequest(gate.ChannelSpotOrderAmend, order))
}

func (s *SpotOrderSuite) queryOrderStatus() {
	<-s.finishAmend

	orderStatus := &model.StatusOrderParam{
		CurrencyPair: s.order.CurrencyPair,
		OrderId:      s.order.Id,
	}

	s.NoError(s.svc.APIRequest(gate.ChannelSpotOrderStatus, orderStatus))
}

func (s *SpotOrderSuite) cancelOrder() {
	<-s.finishQuery
	s.NoError(s.svc.APIRequest(gate.ChannelSpotOrderCancelIds, []*model.Order{s.order}))
}

func TestSpotOrderSuite(t *testing.T) {
	suite.Run(t, new(SpotOrderSuite))
}

func (s *SpotOrderSuite) loginCallback() gate.CallBack {
	return gate.NewCallBack(func(msg *gate.UpdateMsg) {
		log.Info().Msgf("[Login]: %s", msg.Data.Result)
	})
}

func (s *SpotOrderSuite) orderCallback() gate.CallBack {
	return gate.NewCallBack(func(msg *gate.UpdateMsg) {
		if msg.Data.Errs != nil {
			log.Error().Msgf("[Create] label: %s, message: %s", msg.Data.Errs.Label, msg.Data.Errs.Message)
			s.finishCreate <- struct{}{}
			return
		}

		var order resp.SpotOrder
		if err := json.Unmarshal(msg.Data.Result, &order); err != nil {
			log.Error().Msgf("[Create] failed to unmarshal response: %v, msg: %v", err, msg)
			s.finishCreate <- struct{}{}
			return
		}

		if order.Id == "" {
			return
		}

		log.Info().Msgf("[Create] order_id: %s, price: %v, amount: %v, status: %s", order.Id, order.Price, order.Amount, order.Status)
		s.order.Id = order.Id
		s.finishCreate <- struct{}{}
	})
}

func (s *SpotOrderSuite) orderAmendCallback() gate.CallBack {
	return gate.NewCallBack(func(msg *gate.UpdateMsg) {
		defer func() {
			s.finishAmend <- struct{}{}
		}()

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

func (s *SpotOrderSuite) orderStatusCallback() gate.CallBack {
	return gate.NewCallBack(func(msg *gate.UpdateMsg) {
		defer func() {
			s.finishQuery <- struct{}{}
		}()

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

func (s *SpotOrderSuite) orderCancelCallback() gate.CallBack {
	return gate.NewCallBack(func(msg *gate.UpdateMsg) {
		defer func() {
			s.exit <- struct{}{}
		}()

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
