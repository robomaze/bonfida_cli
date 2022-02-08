package api_test

import (
	"context"
	"github.com/robomaze/bonfida_cli/api"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	cli := &api.Client{
		BaseURL:    "https://bad.url",
		UserAgent:  "Bonfida/golang",
		HTTPClient: http.DefaultClient,
		Logger:     log.New(os.Stderr, "bonfida_cli ", log.LstdFlags),
	}

	service := cli.NewOrderBookService()
	service.SetMarketName("BTC/USDT")

	ctx := context.TODO()
	_, err := service.Do(ctx)

	assert.EqualError(t, err, "error getting order book: unable to talk to Bonfida: Get \"https://bad.url/orderbooks/BTCUSDT\": dial tcp: lookup bad.url: no such host")
}

func TestClient_RealRequest(t *testing.T) {
	cli := &api.Client{
		BaseURL:    api.BonfidaApiUrl,
		UserAgent:  "Bonfida/golang",
		HTTPClient: http.DefaultClient,
		Logger:     log.New(os.Stderr, "bonfida_cli ", log.LstdFlags),
	}

	service := cli.NewOrderBookService()
	service.SetMarketName("BTC/USDT")

	ctx := context.TODO()
	res, err := service.Do(ctx)

	assert.NoError(t, err)

	assert.Equal(t, "BTC/USDT", res.Market)
	assert.Equal(t, "C1EuT9VokAKLiW7i2ASnZUvxDoKuKkCpDDeNxAptuNe4", res.Address)
	assert.True(t, len(res.Asks) > 0)
	assert.True(t, len(res.Bids) > 0)
}
