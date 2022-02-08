# Bonfida API Client

[![Coverage Status](https://coveralls.io/repos/github/robomaze/bonfida_cli/badge.svg)](https://coveralls.io/github/robomaze/bonfida_cli) ![Go](https://github.com/robomaze/bonfida_cli/actions/workflows/ci.yml/badge.svg)

Provides a client library written in Golang to interact with Bonfida API https://docs.bonfida.com/#introduction

Note that this is still work in progress.

## Installation

`go get github.com/robomaze/bonfida_cli`

## Usage

```go
package main

import (
	"context"
	"log"

	"github.com/robomaze/bonfida_cli/api"
)

func main() {
	cli := api.NewClient()
	service := cli.NewOrderBookService()

	service.SetMarketName("ETH/USDT")
	ctx := context.TODO()
	res, err := service.Do(ctx)

	if err == nil {
		log.Printf("Order book market: %s", res.Address)
		log.Printf("First ask: %+v", res.Asks[0])
		log.Printf("First bid: %+v", res.Bids[0])
	} else {
		log.Printf("Error %s", err.Error())
	}
}

```

## TODO

- [ ] Implement WS https://docs.bonfida.com/#websocket
