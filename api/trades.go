package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

type TradesService struct {
	c             *Client
	marketName    string
	marketAddress string
}

func (s *TradesService) SetMarketName(name string) *TradesService {
	s.marketName = strings.Replace(name, "/", "", 1)

	return s
}

func (s *TradesService) SetMarketAddress(address string) *TradesService {
	s.marketAddress = address

	return s
}

func (s *TradesService) Do(ctx context.Context) ([]*Trade, error) {
	r := &request{
		apiUrl:   s.c.BaseURL,
		method:   http.MethodGet,
		endpoint: s.getEndpoint(),
	}

	bonfidaResponse, err := s.c.callAPI(ctx, r)
	if err != nil {
		return nil, errors.Wrap(err, "error getting trades")
	}

	var trades []*Trade

	err = json.Unmarshal(bonfidaResponse.Data, &trades)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to unmarshall bonfida response %s", bonfidaResponse.Data)
	}

	return trades, err
}

func (s *TradesService) getEndpoint() string {
	if len(s.marketName) > 0 {
		return fmt.Sprintf("trades/%s", s.marketName)
	}

	if len(s.marketAddress) > 0 {
		return fmt.Sprintf("trades/address/%s", s.marketAddress)
	}

	return "trades/all/recent"
}

type Trade struct {
	Market        string  `json:"market"`
	Price         float64 `json:"price"`
	Size          float64 `json:"size"`
	Side          string  `json:"side"`
	Time          float64 `json:"time"`
	OrderId       string  `json:"orderId"`
	FeeCost       float64 `json:"feeCost"`
	MarketAddress string  `json:"marketAddress"`
}
