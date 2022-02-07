package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/robomaze/bonfida_cli/api/common"

	"github.com/pkg/errors"
)

type OrderBookService struct {
	c          *Client
	marketName string
}

func (s *OrderBookService) SetMarketName(name string) *OrderBookService {
	s.marketName = strings.Replace(name, "/", "", 1)

	return s
}

func (s *OrderBookService) Do(ctx context.Context) (*OrderBook, error) {
	if err := s.validateParams(); err != nil {
		return nil, err
	}

	r := &request{
		method:   http.MethodGet,
		endpoint: fmt.Sprintf("%s/orderbooks/%s", s.c.BaseURL, s.marketName),
	}

	bonfidaResponse, err := s.c.callAPI(ctx, r)
	if err != nil {
		return nil, errors.Wrap(err, "error getting order book")
	}

	orderBook := &OrderBook{}

	err = json.Unmarshal(bonfidaResponse.Data, orderBook)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to unmarshall bonfida response %s", bonfidaResponse.Data)
	}

	return orderBook, nil
}

func (s *OrderBookService) validateParams() error {
	if s.marketName == "" {
		return &common.ErrReqParam{
			Name: "marketName",
			Msg:  "missing",
		}
	}

	return nil
}

type OrderBook struct {
	Market  string `json:"market"`
	Address string `json:"marketAddress"`
	Bids    []Quote
	Asks    []Quote
}

type Quote struct {
	Price float64 `json:"price"`
	Size  float64 `json:"size"`
}
