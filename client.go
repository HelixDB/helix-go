package helix

import (
	"net/http"
	"strings"
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
	return func(o *ClientOption) {
		o.timeout = timeout
	}
}

func NewClient(host string, opts ...ClientOptionFunc) *Client {

	if !strings.HasSuffix(host, "/") {
		host = host + "/"
	}

	option := ClientOption{
		host:    host,
		timeout: 10 * time.Second, // default time
	}

	for _, opt := range opts {
		opt(&option)
	}

	return &Client{
		host: option.host,
		httpClient: &http.Client{
			Timeout: option.timeout,
		},
	}
}
