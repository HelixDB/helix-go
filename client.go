package helix

import (
	"net/http"
	"time"
)

type Client struct {
	BaseURL    string
	Port       int
	HTTPClient *http.Client
}

type ClientOption func(*Client)

func WithPort(port int) ClientOption {
	return func(c *Client) {
		c.Port = port
	}
}

func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.BaseURL = baseURL
	}
}

func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.HTTPClient = httpClient
	}
}

func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.HTTPClient.Timeout = timeout
	}
}

func NewClient(options ...ClientOption) *Client {

	client := &Client{
		BaseURL: "http://localhost",
		Port: "6969"
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		}
	}

	for _, opt := range options {
		opt(client)
	}

	return client
}
