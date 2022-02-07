package api

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// request defines an API request
type request struct {
	method   string
	endpoint string
	query    url.Values
	header   http.Header
	body     io.Reader
	fullURL  string
}

func (r *request) compile() {
	r.fullURL = r.endpoint

	q := r.query.Encode()
	if q != "" {
		r.fullURL = fmt.Sprintf("%s?%s", r.fullURL, q)
	}
}

// setParam set param with key/value to query string
func (r *request) setParam(key string, value interface{}) *request {
	if r.query == nil {
		r.query = url.Values{}
	}
	r.query.Set(key, fmt.Sprintf("%v", value))
	return r
}
