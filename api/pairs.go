package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type PairsService struct {
	c *Client
}

func (s *PairsService) Do(ctx context.Context) ([]string, error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: fmt.Sprintf("%s/pairs", s.c.BaseURL),
	}

	bonfidaResponse, err := s.c.callAPI(ctx, r)
	if err != nil {
		return nil, errors.Wrap(err, "error getting pairs")
	}

	var pairs []string

	err = json.Unmarshal(bonfidaResponse.Data, &pairs)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to unmarshall bonfida response %s", bonfidaResponse.Data)
	}

	return pairs, err
}
