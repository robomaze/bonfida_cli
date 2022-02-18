package api_test

import (
	"log"
	"net/http"
	"os"
	"sync"
	"testing"

	"github.com/robomaze/bonfida_cli/api"
	"github.com/stretchr/testify/assert"
)

func TestWSClient(t *testing.T) {
	var wsError error

	wg := &sync.WaitGroup{}
	wg.Add(1)

	ws := &api.WSClient{
		HTTPClient: http.DefaultClient,
		Logger:     log.New(os.Stderr, "bonfida_cli ", log.LstdFlags),
		OnError: func(err error, msg string) {
			wsError = err
		},
	}

	processor := testTradeProcessor{
		ws: ws,
		wg: wg,
	}

	ws.OnData = processor.process

	err := ws.Subscribe()

	// Wait until the test processor receives a trade, and we are good to test it.
	wg.Wait()

	assert.NoError(t, err)
	assert.NoError(t, wsError)

	assert.NotNil(t, processor.trade.Price)
	assert.NotNil(t, processor.trade.Market)
}

type testTradeProcessor struct {
	ws    *api.WSClient
	trade *api.Trade
	wg    *sync.WaitGroup
}

func (p *testTradeProcessor) process(trade *api.Trade) {
	p.trade = trade

	// We got our one trade for testing purposes so close the connection and notify above we are done here.
	p.ws.Close()
	p.wg.Done()
}
