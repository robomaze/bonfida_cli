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

func TestVolumesService_Do(t *testing.T) {
	httpClient := http.Client{}
	httpmock.ActivateNonDefault(&httpClient)

	cli := &api.Client{
		BaseURL:    api.BonfidaApiUrl,
		UserAgent:  "Bonfida/golang",
		HTTPClient: &httpClient,
		Logger:     log.New(os.Stderr, "bonfida_cli ", log.LstdFlags),
	}

	service := cli.NewVolumesService()

	ctx := context.TODO()
	_, err := service.Do(ctx)

	assert.EqualError(t, err, "param name 'marketName' with error 'missing'")

	httpmock.RegisterResponder("GET",
		fmt.Sprintf("%s/volumes/ETHUSDT", api.BonfidaApiUrl),
		httpmock.NewStringResponder(500, "{}"))

	service.SetMarketName("ETH/USDT")

	_, err = service.Do(ctx)

	assert.EqualError(t, err, "error getting volumes: status 500 at https://serum-api.bonfida.com/volumes/ETHUSDT: Bonfida API error")

	httpmock.RegisterResponder("GET",
		fmt.Sprintf("%s/volumes/ETHUSDT", api.BonfidaApiUrl),
		httpmock.NewStringResponder(200, "{\n  \"success\": true,\n  \"data\": [\n    {\n      \"volumeUsd\": 377446,\n      \"volume\": 835\n    }\n  ]\n}"))

	res, err := service.Do(ctx)

	assert.NoError(t, err)

	assert.Equal(t, 835.0, res[0].Volume)
	assert.Equal(t, 377446.0, res[0].VolumeUsd)
}
