package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

type PairsService struct {
	c *Client
}

func (s *PairsService) Do(ctx context.Context) ([]string, error) {
	r := &request{
		apiUrl:   s.c.BaseURL,
		method:   http.MethodGet,
		endpoint: "pairs",
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
