package api

import (
	"github.com/artsafin/coda-schema-generator/dto"
	"log"
	"net/http"
)

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (fn roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func NewHTTPClient(options dto.APIOptions) *client {
	c := &client{
		opts:   options,
		http:   nil,
		logger: log.Default(),
	}
	c.http = &http.Client{
		Transport: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			c.logvf("request: %s\n", req.URL.String())
			req.Header.Set("Authorization", "Bearer "+options.Token)
			return http.DefaultTransport.RoundTrip(req)
		}),
		Timeout: options.RequestTimeout,
	}

	return c
}
