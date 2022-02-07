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

func TestPairsService_Do(t *testing.T) {
	httpClient := http.Client{}
	httpmock.ActivateNonDefault(&httpClient)

	cli := &api.Client{
		BaseURL:    api.BonfidaApiUrl,
		UserAgent:  "Bonfida/golang",
		HTTPClient: &httpClient,
		Logger:     log.New(os.Stderr, "bonfida_cli ", log.LstdFlags),
	}

	httpmock.RegisterResponder("GET",
		fmt.Sprintf("%s/pairs", api.BonfidaApiUrl),
		httpmock.NewStringResponder(500, "{}"))

	service := cli.NewPairsService()

	ctx := context.TODO()
	_, err := service.Do(ctx)

	assert.EqualError(t, err, "error getting pairs: status 500 at https://serum-api.bonfida.com/pairs: Bonfida API error")

	httpmock.RegisterResponder("GET",
		fmt.Sprintf("%s/pairs", api.BonfidaApiUrl),
		httpmock.NewStringResponder(200, "{\"success\":true,\"data\":[\"BTC/USDT\",\"ETH/USDT\"]}"))

	res, err := service.Do(ctx)

	assert.NoError(t, err)

	assert.Equal(t, "BTC/USDT", res[0])
	assert.Equal(t, "ETH/USDT", res[1])
}
