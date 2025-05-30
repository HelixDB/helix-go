package helix

import (
	"net/http"
	"time"
)

type Client struct {
	host       string
	httpClient *http.Client
}

type ClientOption struct {
	host    string
	timeout time.Duration
}

type ClientOptionFunc func(*ClientOption)

func WithTimeout(timeout time.Duration) ClientOptionFunc {
	return func(c *ClientOption) {
		c.timeout = timeout
	}
}

func NewClient(host string, options ...ClientOptionFunc) *Client {

	opts := &ClientOption{
		host:    host,
		timeout: 10 * time.Second, // default time
	}

	for _, opt := range options {
		opt(opts)
	}

	return &Client{
		host: opts.host,
		httpClient: &http.Client{
			Timeout: opts.timeout,
		},
	}
}
