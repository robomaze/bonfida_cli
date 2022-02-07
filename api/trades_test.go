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

func TestTradesService_Do(t *testing.T) {
	httpClient := http.Client{}
	httpmock.ActivateNonDefault(&httpClient)

	cli := &api.Client{
		BaseURL:    api.BonfidaApiUrl,
		UserAgent:  "Bonfida/golang",
		HTTPClient: &httpClient,
		Logger:     log.New(os.Stderr, "bonfida_cli ", log.LstdFlags),
	}

	httpmock.RegisterResponder("GET",
		fmt.Sprintf("%s/trades/all/recent", api.BonfidaApiUrl),
		httpmock.NewStringResponder(500, "{}"))

	service := cli.NewTradesService()

	ctx := context.TODO()
	_, err := service.Do(ctx)

	assert.EqualError(t, err, "error getting trades: status 500 at https://serum-api.bonfida.com/trades/all/recent: Bonfida API error")

	httpmock.RegisterResponder("GET",
		fmt.Sprintf("%s/trades/ETHUSDT", api.BonfidaApiUrl),
		httpmock.NewStringResponder(200, "{\n  \"success\": true,\n  \"data\": [\n    {\n      \"market\": \"ETH/USDT\",\n      \"price\": 451.51,\n      \"size\": 0.5,\n      \"side\": \"buy\",\n      \"time\": 1604767562476.2188,\n      \"orderId\": \"833220983065386731245551\",\n      \"feeCost\": 0.225755,\n      \"marketAddress\": \"5abZGhrELnUnfM9ZUnvK6XJPoBU5eShZwfFPkdhAC7o\"\n    }\n  ]\n}\n"))

	service.SetMarketName("ETH/USDT")

	res, err := service.Do(ctx)

	assert.NoError(t, err)
	assert.Equal(t, "ETH/USDT", res[0].Market)
	assert.Equal(t, "5abZGhrELnUnfM9ZUnvK6XJPoBU5eShZwfFPkdhAC7o", res[0].MarketAddress)
	assert.Equal(t, 0.5, res[0].Size)
	assert.Equal(t, 451.51, res[0].Price)
	assert.Equal(t, 0.225755, res[0].FeeCost)
	assert.Equal(t, "buy", res[0].Side)
	assert.Equal(t, 1604767562476.2188, res[0].Time)
	assert.Equal(t, "833220983065386731245551", res[0].OrderId)

	httpmock.RegisterResponder("GET",
		fmt.Sprintf("%s/trades/address/7dLVkUfBVfCGkFhSXDCq1ukM9usathSgS716t643iFGF", api.BonfidaApiUrl),
		httpmock.NewStringResponder(200, "{\n  \"success\": true,\n  \"data\": [\n    {\n      \"market\": \"ETH/USDT\",\n      \"price\": 451.51,\n      \"size\": 0.5,\n      \"side\": \"buy\",\n      \"time\": 1604767562476.2188,\n      \"orderId\": \"833220983065386731245551\",\n      \"feeCost\": 0.225755,\n      \"marketAddress\": \"5abZGhrELnUnfM9ZUnvK6XJPoBU5eShZwfFPkdhAC7o\"\n    }\n  ]\n}\n"))

	service.SetMarketName("")
	service.SetMarketAddress("7dLVkUfBVfCGkFhSXDCq1ukM9usathSgS716t643iFGF")

	res, err = service.Do(ctx)

	assert.NoError(t, err)
	assert.Equal(t, "ETH/USDT", res[0].Market)

	httpmock.RegisterResponder("GET",
		fmt.Sprintf("%s/trades/all/recent", api.BonfidaApiUrl),
		httpmock.NewStringResponder(200, "{\n  \"success\": true,\n  \"data\": [\n    {\n      \"market\": \"ETH/USDT\",\n      \"price\": 451.51,\n      \"size\": 0.5,\n      \"side\": \"buy\",\n      \"time\": 1604767562476.2188,\n      \"orderId\": \"833220983065386731245551\",\n      \"feeCost\": 0.225755,\n      \"marketAddress\": \"5abZGhrELnUnfM9ZUnvK6XJPoBU5eShZwfFPkdhAC7o\"\n    }\n  ]\n}\n"))

	service.SetMarketName("")
	service.SetMarketAddress("")

	res, err = service.Do(ctx)

	assert.NoError(t, err)
	assert.Equal(t, "ETH/USDT", res[0].Market)
}
