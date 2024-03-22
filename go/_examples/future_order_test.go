package main

import (
	"encoding/json"
	"os"
	"strconv"
	"testing"

	gate "github.com/gateio/gatews/go"
	"github.com/gateio/gatews/go/model"
	"github.com/gateio/gatews/go/resp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"
)

type FutureOrderSuite struct {
	suite.Suite

	svc   *gate.WsService
	order *model.FuturesOrder

	finishCreate chan struct{}
	finishAmend  chan struct{}
	finishQuery  chan struct{}
	exit         chan struct{}
}

func (s *FutureOrderSuite) SetupSuite() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	var err error
	s.svc, err = gate.NewWsService(nil, nil, gate.NewConnConfFromOption(&gate.ConfOptions{
		App:          "futures",
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

	s.svc.SetCallBack(gate.ChannelFutureLogin, s.loginCallback())
	s.svc.SetCallBack(gate.ChannelFutureOrderPlace, s.orderCallback())
	s.svc.SetCallBack(gate.ChannelFutureOrderAmend, s.orderAmendCallback())
	s.svc.SetCallBack(gate.ChannelFutureOrderStatus, s.orderStatusCallback())
	s.svc.SetCallBack(gate.ChannelFutureOrderCancel, s.orderCancelCallback())

	s.order = &model.FuturesOrder{
		Contract: "BTC_USDT",
		Size:     100,
		Iceberg:  0,
		Price:    "30000",
		Text:     "t-my-custom-id",
	}
}

func (s *FutureOrderSuite) TestOrder() {
	s.createOrder()
	s.amendOrder()
	s.queryOrderStatus()
	s.cancelOrder()

	<-s.exit
}

func (s *FutureOrderSuite) createOrder() {
	s.NoError(s.svc.APIRequest(gate.ChannelFutureOrderPlace, s.order))
}

func (s *FutureOrderSuite) amendOrder() {
	<-s.finishCreate

	order := &model.AmendFuturesOrder{
		OrderId:   strconv.Itoa(int(s.order.Id)),
		Settle:    "USDT",
		Price:     "40000",
		AmendText: s.order.AmendText,
		Size:      s.order.Size,
	}

	s.NoError(s.svc.APIRequest(gate.ChannelFutureOrderAmend, order))
}

func (s *FutureOrderSuite) queryOrderStatus() {
	<-s.finishAmend

	orderStatus := &model.StatusFuturesOrder{
		Settle:  "usdt",
		OrderId: strconv.Itoa(int(s.order.Id)),
	}

	s.NoError(s.svc.APIRequest(gate.ChannelFutureOrderStatus, orderStatus))
}

func (s *FutureOrderSuite) cancelOrder() {
	<-s.finishQuery

	order := &model.CancelFuturesOrder{
		Settle:  "usdt",
		OrderId: strconv.Itoa(int(s.order.Id)),
	}
	s.NoError(s.svc.APIRequest(gate.ChannelFutureOrderCancel, order))
}

func (s *FutureOrderSuite) loginCallback() gate.CallBack {
	return gate.NewCallBack(func(msg *gate.UpdateMsg) {
		log.Error().Msgf("[Login] result: %s", msg.Data.Result)
	})
}

func (s *FutureOrderSuite) orderCallback() gate.CallBack {
	return gate.NewCallBack(func(msg *gate.UpdateMsg) {
		if msg.Data.Errs != nil {
			log.Error().Msgf("[Create] label: %s, message: %s", msg.Data.Errs.Label, msg.Data.Errs.Message)
			s.finishCreate <- struct{}{}
			return
		}

		var order resp.FutureOrder
		if err := json.Unmarshal(msg.Data.Result, &order); err != nil {
			log.Error().Msgf("[Create] failed to unmarshal response: %v, msg: %v", err, msg)
			s.finishCreate <- struct{}{}
			return
		}

		if order.Id == 0 {
			return
		}

		log.Info().Msgf("[Create] order_id: %v, status: %s", order.Id, order.Status)
		s.order.Id = order.Id
		s.finishCreate <- struct{}{}
	})
}

func (s *FutureOrderSuite) orderAmendCallback() gate.CallBack {
	return gate.NewCallBack(func(msg *gate.UpdateMsg) {
		defer func() {
			s.finishAmend <- struct{}{}
		}()

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

func (s *FutureOrderSuite) orderStatusCallback() gate.CallBack {
	return gate.NewCallBack(func(msg *gate.UpdateMsg) {
		defer func() {
			s.finishQuery <- struct{}{}
		}()

		if msg.Data.Errs != nil {
			log.Info().Msgf("[Query] label: %s, message: %s", msg.Data.Errs.Label, msg.Data.Errs.Message)
			return
		}

		var order resp.FutureOrder
		if err := json.Unmarshal(msg.Data.Result, &order); err != nil {
			log.Info().Msgf("[Query] failed to unmarshal response: %v, msg: %v", err, msg)
			return
		}

		log.Info().Msgf("[Query] order_id: %v, status: %v", order.Id, order.Status)
	})
}

func (s *FutureOrderSuite) orderCancelCallback() gate.CallBack {
	return gate.NewCallBack(func(msg *gate.UpdateMsg) {
		defer func() {
			s.exit <- struct{}{}
		}()

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

func TestFutureOrderSuite(t *testing.T) {
	suite.Run(t, new(FutureOrderSuite))
}
