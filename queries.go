package helix

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
)

type HelixInput map[string]any

type HelixResponse struct {
	bytes []byte
	err   error
}

type QueryOption struct {
	data     HelixInput
	datatype any
}

type QueryOptionFunc func(*QueryOption)

func WithData(data HelixInput) QueryOptionFunc {
	return func(o *QueryOption) {
		o.data = data
	}
}

func WithTarget(datatype any) QueryOptionFunc {
	return func(o *QueryOption) {
		o.datatype = datatype
	}
}

func (c *Client) Query(endpoint string, opts ...QueryOptionFunc) *HelixResponse {

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
		return &HelixResponse{
			bytes: nil,
			err:   fmt.Errorf("failed to marshal input data: %w", err),
		}
	}

	url := c.host + endpoint
	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonData))
	if err != nil {
		return &HelixResponse{
			bytes: nil,
			err:   fmt.Errorf("failed to create request: %w", err),
		}
	}

	req.Header.Set("Content-Type", "application/json")
	// Authorization token in the future maybe?

	res, err := c.httpClient.Do(req)
	if err != nil {
		return &HelixResponse{
			bytes: nil,
			err:   fmt.Errorf("failed to send request: %w", err),
		}
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return &HelixResponse{
			bytes: nil,
			err:   fmt.Errorf("HTTP error %d: %s", res.StatusCode, string(body)),
		}
	}

	return &HelixResponse{
		bytes: body,
		err:   nil,
	}
}

func (r *HelixResponse) AsMap() (map[string]any, error) {

	if r.err != nil {
		return nil, r.err
	}

	var mapResponse map[string]any
	err := json.Unmarshal(r.bytes, &mapResponse)
	if err != nil {
		return nil, err
	}

	return mapResponse, nil
}

func (r *HelixResponse) Scan(dest any) error {

	if r.err != nil {
		return r.err
	}

	rv := reflect.ValueOf(dest)
	if rv.Kind() != reflect.Ptr {
		return fmt.Errorf("scan destination must be a pointer")
	}

	if rv.IsNil() {
		return fmt.Errorf("scan destination cannot be nil")
	}

	return json.Unmarshal(r.bytes, dest)
}
