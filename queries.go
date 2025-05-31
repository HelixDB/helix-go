package helix

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type HelixInput map[string]any

type HelixResponse map[string]any

type QueryOption struct {
	data HelixInput
}

type QueryOptionFunc func(*QueryOption)

func WithData(data HelixInput) QueryOptionFunc {
	return func(q *QueryOption) {
		q.data = data
	}
}

func (c *Client) Query(endpoint string, opts ...QueryOptionFunc) (HelixResponse, error) {

	option := QueryOption{}
	for _, opt := range opts {
		opt(&option)
	}

	data := option.data
	if data == nil {
		data = make(HelixInput)
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal input data: %w", err)
	}

	url := c.host + endpoint
	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	// Authorization token in the future maybe?

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("HTTP error %d: %s", res.StatusCode, string(body))
	}

	var response HelixResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return response, nil
}
