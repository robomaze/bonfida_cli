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

type VolumesService struct {
	c          *Client
	marketName string
}

func (s *VolumesService) SetMarketName(name string) *VolumesService {
	s.marketName = strings.Replace(name, "/", "", 1)

	return s
}

func (s *VolumesService) Do(ctx context.Context) ([]*Volume, error) {
	if err := s.validateParams(); err != nil {
		return nil, err
	}

	r := &request{
		method:   http.MethodGet,
		endpoint: fmt.Sprintf("%s/volumes/%s", s.c.BaseURL, s.marketName),
	}

	bonfidaResponse, err := s.c.callAPI(ctx, r)
	if err != nil {
		return nil, errors.Wrap(err, "error getting volumes")
	}

	var volumes []*Volume

	err = json.Unmarshal(bonfidaResponse.Data, &volumes)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to unmarshall bonfida response %s", bonfidaResponse.Data)
	}

	return volumes, nil
}

func (s *VolumesService) validateParams() error {
	if s.marketName == "" {
		return &common.ErrReqParam{
			Name: "marketName",
			Msg:  "missing",
		}
	}

	return nil
}

type Volume struct {
	VolumeUsd float64 `json:"volumeUsd"`
	Volume    float64 `json:"volume"`
}
