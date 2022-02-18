# Bonfida API Client

[![Coverage Status](https://coveralls.io/repos/github/robomaze/bonfida_cli/badge.svg)](https://coveralls.io/github/robomaze/bonfida_cli) ![Go](https://github.com/robomaze/bonfida_cli/actions/workflows/ci.yml/badge.svg)

Provides a client library written in Golang to interact with Bonfida API https://docs.bonfida.com/#introduction

Note that this is still work in progress.

## Installation

`go get github.com/robomaze/bonfida_cli`

## Usage

### REST API
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

### WS

```go
package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/robomaze/bonfida_cli/api"
)

func main() {
	logger := log.New(os.Stderr, "bonfida_cli ", log.LstdFlags)
	cli := api.WSClient{
		HTTPClient: http.DefaultClient,
		Logger:     logger,
		OnError: func(err error, msg string) {
			// Do stuff when things go feral.
		},
		OnData: func(trade *api.Trade) {
			// Do what you need to do when a trade happens.
			log.Printf("\ntrade %+v", *trade)
		},
	}

	err := cli.Subscribe()
	if err != nil {
		// Deal with error during the subscription process.
    }

	// Only to wait for some trades to appear in logs.
	time.Sleep(time.Second * 60)

	// Properly close the connection when done receiving trades.
	cli.Close()
}

```

## TODO

- [ ] Expand documentation
- [ ] Increase test coverage over 90%
