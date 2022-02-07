package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"github.com/robomaze/bonfida_cli/api/common"
)

// CandlesService is supposed to provide candles, but it does not seem to work at the Bonfida side.
type CandlesService struct {
	c          *Client
	marketName string
	resolution int
	start      int64
	end        int64
	limit      int
}

func (s *CandlesService) SetMarketName(name string) *CandlesService {
	s.marketName = strings.Replace(name, "/", "", 1)

	return s
}

// SetResolution sets the window length.
//
// Allowed values: 60, 3600, 14400, 86400.
func (s *CandlesService) SetResolution(resolution int) *CandlesService {
	s.resolution = resolution

	return s
}

func (s *CandlesService) SetStartEnd(start int64, end int64) *CandlesService {
	s.start = start
	s.end = end

	return s
}

func (s *CandlesService) SetLimit(limit int) *CandlesService {
	s.limit = limit

	return s
}

func (s *CandlesService) Do(ctx context.Context) ([]*Candle, error) {
	if err := s.validateParams(); err != nil {
		return nil, err
	}

	r := &request{
		method:   http.MethodGet,
		endpoint: fmt.Sprintf("%s/candles/%s", s.c.BaseURL, s.marketName),
	}

	if s.resolution > 0 {
		r.setParam("resolution", s.resolution)
	}

	if s.start > 0 {
		r.setParam("startTime", s.start)
	}

	if s.end > 0 {
		r.setParam("endTime", s.end)
	}

	if s.limit > 0 {
		r.setParam("limit", s.limit)
	}

	bonfidaResponse, err := s.c.callAPI(ctx, r)
	if err != nil {
		return nil, errors.Wrap(err, "error getting candles")
	}

	var candles []*Candle

	err = json.Unmarshal(bonfidaResponse.Data, &candles)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to unmarshall bonfida response %s", bonfidaResponse.Data)
	}

	return candles, nil
}

func (s *CandlesService) validateParams() error {
	if s.marketName == "" {
		return &common.ErrReqParam{
			Name: "marketName",
			Msg:  "missing",
		}
	}

	if s.resolution == 0 {
		return &common.ErrReqParam{
			Name: "resolution",
			Msg:  "missing",
		}
	}

	return nil
}

type Candle struct {
	Close       float64 `json:"close"`
	Open        float64 `json:"open"`
	Low         float64 `json:"low"`
	High        float64 `json:"high"`
	StartTime   int64   `json:"startTime"`
	Market      string  `json:"market"`
	VolumeBase  int     `json:"volumeBase"`
	VolumeQuote int     `json:"volumeQuote"`
}
