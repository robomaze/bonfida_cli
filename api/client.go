package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/pkg/errors"
	"github.com/robomaze/bonfida_cli/api/common"
)

const BonfidaApiUrl = "https://serum-api.bonfida.com"

var errBonfidaAPI = errors.New("Bonfida API error")

// NewClient will instantiate a ready to use Bonfida API Client.
func NewClient() *Client {
	return &Client{
		BaseURL:    BonfidaApiUrl,
		UserAgent:  "Bonfida/golang",
		HTTPClient: http.DefaultClient,
		Logger:     log.New(os.Stderr, "bonfida_cli ", log.LstdFlags),
	}
}

type Client struct {
	BaseURL    string
	UserAgent  string
	HTTPClient *http.Client
	Logger     *log.Logger
}

func (c *Client) NewOrderBookService() *OrderBookService {
	return &OrderBookService{c: c}
}

func (c *Client) NewPairsService() *PairsService {
	return &PairsService{c: c}
}

func (c *Client) NewTradesService() *TradesService {
	return &TradesService{c: c}
}

func (c *Client) NewVolumesService() *VolumesService {
	return &VolumesService{c: c}
}

func (c *Client) callAPI(ctx context.Context, r *request) (bonfidaResponse *common.BonfidaResponse, err error) {
	req, err := http.NewRequest(r.method, r.getFullUrl(), r.body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create request object")
	}

	req = req.WithContext(ctx)
	req.Header = r.header

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "unable to talk to Bonfida")
	}

	if res.StatusCode >= http.StatusBadRequest {
		return nil, errors.Wrapf(errBonfidaAPI, "status %d at %s", res.StatusCode, req.URL)
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read Bonfida response body")
	}

	bonfidaResponse = &common.BonfidaResponse{}

	err = json.Unmarshal(data, bonfidaResponse)
	if err != nil {
		return nil, errors.Wrap(err, "unable to unmarshal Bonfida response body")
	}

	defer func() {
		closeErr := res.Body.Close()
		if err == nil && closeErr != nil {
			err = closeErr
		}
	}()

	if !bonfidaResponse.Success {
		return nil, errors.Wrap(errBonfidaAPI, "Bonfida probably did not like the request data")
	}

	return bonfidaResponse, nil
}

// request defines an API request
type request struct {
	apiUrl   string
	method   string
	endpoint string
	header   http.Header
	body     io.Reader
}

func (r *request) getFullUrl() string {
	return fmt.Sprintf("%s/%s", r.apiUrl, r.endpoint)
}
