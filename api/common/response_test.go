package common_test

import (
	"encoding/json"
	"testing"

	"github.com/robomaze/bonfida_cli/api/common"
	"github.com/stretchr/testify/assert"
)

func TestBonfidaResponse_UnmarshalJSON(t *testing.T) {
	jsonStr := "{\n  \"success\": true,\n  \"data\": {\n    \"market\": \"ETH/USDT\",\n    \"bids\": [\n      { \"price\": 452.77, \"size\": 5 },\n      { \"price\": 452.71, \"size\": 0.5 },\n      { \"price\": 452.17, \"size\": 10 }\n    ],\n    \"asks\": [\n      { \"price\": 453.19, \"size\": 105.534 },\n      { \"price\": 453.41, \"size\": 10 },\n      { \"price\": 453.49, \"size\": 114.203 }\n    ],\n    \"marketAddress\": \"5abZGhrELnUnfM9ZUnvK6XJPoBU5eShZwfFPkdhAC7o\"\n  }\n}"

	bonfidaRes := common.BonfidaResponse{}

	err := json.Unmarshal([]byte(jsonStr), &bonfidaRes)

	assert.NoError(t, err)

	assert.Equal(t, true, bonfidaRes.Success)
	assert.Equal(t, "{\"asks\":[{\"price\":453.19,\"size\":105.534},{\"price\":453.41,\"size\":10},{\"price\":453.49,\"size\":114.203}],\"bids\":[{\"price\":452.77,\"size\":5},{\"price\":452.71,\"size\":0.5},{\"price\":452.17,\"size\":10}],\"market\":\"ETH/USDT\",\"marketAddress\":\"5abZGhrELnUnfM9ZUnvK6XJPoBU5eShZwfFPkdhAC7o\"}", string(bonfidaRes.Data))
}
