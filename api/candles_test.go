package api_test

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/robomaze/bonfida_cli/api"
)

func TestCandlesService_Do(t *testing.T) {
	httpClient := http.Client{}

	cli := &api.Client{
		BaseURL:    api.BonfidaApiUrl,
		UserAgent:  "Bonfida/golang",
		HTTPClient: &httpClient,
		Logger:     log.New(os.Stderr, "bonfida_cli ", log.LstdFlags),
	}

	service := cli.NewCandlesService()
	service.SetMarketName("BTC/USDC")
	service.SetResolution(14400)

	ctx := context.TODO()

	_, err := service.Do(ctx)
	if err != nil {
		fmt.Printf("%s", err.Error())
	}
}
