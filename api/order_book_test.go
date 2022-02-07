package api_test

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/robomaze/bonfida_cli/api"
	"github.com/stretchr/testify/assert"
)

func TestOrderBookService_Do(t *testing.T) {
	httpClient := http.Client{}
	httpmock.ActivateNonDefault(&httpClient)

	cli := &api.Client{
		BaseURL:    api.BonfidaApiUrl,
		UserAgent:  "Bonfida/golang",
		HTTPClient: &httpClient,
		Logger:     log.New(os.Stderr, "bonfida_cli ", log.LstdFlags),
	}

	service := cli.NewOrderBookService()

	ctx := context.TODO()
	_, err := service.Do(ctx)

	assert.EqualError(t, err, "param name 'marketName' with error 'missing'")

	service = cli.NewOrderBookService()
	service.SetMarketName("BTC/USDT")

	httpmock.RegisterResponder("GET",
		fmt.Sprintf("%s/orderbooks/BTCUSDT", api.BonfidaApiUrl),
		httpmock.NewStringResponder(500, "{}"))

	_, err = service.Do(ctx)

	assert.EqualError(t, err, "error getting order book: status 500 at https://serum-api.bonfida.com/orderbooks/BTCUSDT: Bonfida API error")

	httpmock.RegisterResponder("GET",
		fmt.Sprintf("%s/orderbooks/BTCUSDT", api.BonfidaApiUrl),
		httpmock.NewStringResponder(200, "{\"success\":false,\"data\":{\"market\":\"\",\"bids\":null,\"asks\":null,\"marketAddress\":\"\"}}"))

	res, err := service.Do(ctx)

	assert.Nil(t, res)
	assert.EqualError(t, err, "error getting order book: Bonfida probably did not like the request data: Bonfida API error")

	httpmock.RegisterResponder("GET",
		fmt.Sprintf("%s/orderbooks/BTCUSDT", api.BonfidaApiUrl),
		httpmock.NewStringResponder(200, "{\n  \"success\": true,\n  \"data\": {\n    \"market\": \"ETH/USDT\",\n    \"bids\": [\n      { \"price\": 452.77, \"size\": 5 },\n      { \"price\": 452.71, \"size\": 0.5 },\n      { \"price\": 452.17, \"size\": 10 }\n    ],\n    \"asks\": [\n      { \"price\": 453.19, \"size\": 105.534 },\n      { \"price\": 453.41, \"size\": 10 },\n      { \"price\": 453.49, \"size\": 114.203 }\n    ],\n    \"marketAddress\": \"5abZGhrELnUnfM9ZUnvK6XJPoBU5eShZwfFPkdhAC7o\"\n  }\n}"))

	res, err = service.Do(ctx)

	assert.NoError(t, err)

	assert.Equal(t, "ETH/USDT", res.Market)
	assert.Equal(t, "5abZGhrELnUnfM9ZUnvK6XJPoBU5eShZwfFPkdhAC7o", res.Address)
	assert.Equal(t, 453.19, res.Asks[0].Price)
	assert.Equal(t, 105.534, res.Asks[0].Size)
	assert.Equal(t, 452.77, res.Bids[0].Price)
	assert.Equal(t, 5.0, res.Bids[0].Size)
}
