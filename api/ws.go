package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

const (
	subscribeURL   = "https://serum-ws.bonfida.com/subscribe"
	maxRsbAttempts = 3
)

type WSClient struct {
	Logger     *log.Logger
	HTTPClient *http.Client
	OnError    func(err error, msg string)
	OnData     func(trade *Trade)

	wsCloseContext context.Context
	wsCloseFunc    context.CancelFunc
	wsConn         *websocket.Conn
	wsURL          *subscribeRes
}

func (c *WSClient) Subscribe() error {
	err := c.loadSubscribeURL()
	if err != nil {
		return errors.Wrap(err, "unable to load Bonfida subscribe URL")
	}

	c.wsConn, _, err = websocket.DefaultDialer.Dial(c.wsURL.URL, nil)
	if err != nil {
		return errors.Wrapf(err, "error dialing to %s", c.wsURL.URL)
	}

	c.wsCloseContext, c.wsCloseFunc = context.WithCancel(context.Background())

	c.Logger.Printf("Subscribed to Bonfida WS")

	go c.websocketWorker(c.wsCloseContext)

	return nil
}

func (c *WSClient) Close() {
	c.wsCloseFunc()
}

func (c *WSClient) websocketWorker(closeCtx context.Context) {
	for {
		select {
		case <-closeCtx.Done():
			c.closeWSConnection()

			return
		default:
			_, message, err := c.wsConn.ReadMessage()
			if err != nil {
				if resubscribed := c.resubscribeRetry(); !resubscribed {
					c.OnError(err, "failed to resubscribe")

					return
				}

				continue
			}

			var apiResponse *Trade

			err = json.Unmarshal(message, &apiResponse)
			if err != nil {
				c.Logger.Printf("Error unmarshalling response: %s", err.Error())

				continue
			}

			c.OnData(apiResponse)
		}
	}
}

func (c *WSClient) resubscribeRetry() bool {
	c.Logger.Printf("resubscribing to Bonfida WS")

	for attempt := 0; attempt < maxRsbAttempts; attempt++ {
		err := c.Subscribe()
		if err == nil {
			return true
		}

		time.Sleep(time.Second)
	}

	c.Logger.Printf("failed to resubscribe to Bonfida WS")

	return false
}

func (c *WSClient) loadSubscribeURL() error {
	body := bytes.NewReader([]byte("{ \"channel\": \"DEX\" }"))

	req, err := http.NewRequest(http.MethodPost, subscribeURL, body)
	if err != nil {
		return errors.Wrap(err, "cannot create request object")
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "unable to get WS URL")
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.Wrap(err, "unable to read Bonfida response body")
	}

	c.wsURL = &subscribeRes{}

	err = json.Unmarshal(data, c.wsURL)
	if err != nil {
		return errors.Wrapf(err, "unable to unmarshal Bonfida response %s", string(data))
	}

	_ = res.Body.Close()

	return nil
}

func (c *WSClient) getUnsubscribeURL() string {
	parts := strings.Split(c.wsURL.URL, "/ws/")

	return fmt.Sprintf("https://serum-ws.bonfida.com/unsubscribe/%s", parts[1])
}

func (c *WSClient) closeWSConnection() {
	req, err := http.NewRequest(http.MethodPost, c.getUnsubscribeURL(), nil)
	if err != nil {
		c.Logger.Printf("unable to create request object to unsubscribe %s", err.Error())
		return
	}

	_, err = c.HTTPClient.Do(req)
	if err != nil {
		c.Logger.Printf("unable to unsubscribe %s", err.Error())
	}
}

type subscribeRes struct {
	URL string `json:"url"`
}
