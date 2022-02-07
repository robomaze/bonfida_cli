package api_test

import (
	"context"
	"testing"

	"github.com/robomaze/bonfida_cli/api"
	"github.com/stretchr/testify/assert"
)

func TestClient_NewOrderBookService(t *testing.T) {
	cli := api.NewClient()

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
