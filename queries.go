package helix

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
)

type Response struct {
	body []byte
}

type QueryOption struct {
	data any
}

type QueryOptionFunc func(*QueryOption)

func WithData(data any) QueryOptionFunc {
	return func(o *QueryOption) {
		o.data = data
	}
}

func (c *Client) Query(endpoint string, opts ...QueryOptionFunc) (*Response, error) {

	option := QueryOption{}
	for _, opt := range opts {
		opt(&option)
	}

	jsonData, err := marshalInput(option.data)
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

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("%d: %s", res.StatusCode, string(body))
	}

	return &Response{
		body: body,
	}, nil
}

func (r *Response) Raw() []byte {
	return r.body
}

func (r *Response) AsMap() (map[string]any, error) {

	var mapResponse map[string]any
	err := json.Unmarshal(r.body, &mapResponse)
	if err != nil {
		return nil, err
	}

	return mapResponse, nil
}

type ScanOption struct {
	dest any
	name string
}

type ScanOptionFunc func(*ScanOption)

func WithDest(name string, dest any) ScanOptionFunc {
	return func(o *ScanOption) {
		o.name = name
		o.dest = dest
	}
}

// Double json.Unmarshal when `len(args) == 1` and `WithDest(...)` is being used and when `len(args) > 1`
func (r *Response) Scan(args ...any) error {

	if len(args) == 0 {
		return fmt.Errorf("scan destination is expected")
	}

	if len(args) == 1 {
		err := validateDestPointer(args[0])
		if err != nil {
			optFunc, err := validateDestOption(args[0])
			if err != nil {
				return err
			}

			var jsonData map[string]json.RawMessage

			err = json.Unmarshal(r.body, &jsonData)
			if err != nil {
				return fmt.Errorf("invalid json response: %w", err)
			}

			err = scanOption(optFunc, jsonData)
			if err != nil {
				return err
			}

			return nil
		}

		return json.Unmarshal(r.body, args[0])

	}

	var jsonData map[string]json.RawMessage

	err := json.Unmarshal(r.body, &jsonData)
	if err != nil {
		return fmt.Errorf("invalid json response: %w", err)
	}

	for i, arg := range args {
		optFunc, ok := arg.(ScanOptionFunc)
		if !ok {
			return fmt.Errorf("invalid scan argument at position %d: got %T (when passing multiple arguments, only WithDest(fieldName, &dest) is allowed)", i, arg)
		}

		err := scanOption(optFunc, jsonData)
		if err != nil {
			return err
		}
	}

	return nil
}

func scanOption(optFunc ScanOptionFunc, jsonData map[string]json.RawMessage) error {
	var opt ScanOption
	optFunc(&opt)

	err := validateDestPointer(opt.dest)
	if err != nil {
		return err
	}

	rawData, ok := jsonData[opt.name]
	if !ok {
		return fmt.Errorf("field \"%s\" not found", opt.name)
	}

	err = json.Unmarshal(rawData, opt.dest)
	if err != nil {
		return fmt.Errorf("failed to scan field \"%s\": %w", opt.name, err)
	}

	return nil
}

func validateDestPointer(dest any) error {
	rv := reflect.ValueOf(dest)
	if rv.Kind() != reflect.Ptr {
		return fmt.Errorf("scan destination must be a pointer")
	}

	if rv.IsNil() {
		return fmt.Errorf("scan destination cannot be nil")
	}

	return nil
}

func validateDestOption(dest any) (ScanOptionFunc, error) {
	optFunc, ok := dest.(ScanOptionFunc)
	if !ok {
		return nil, fmt.Errorf("invalid scan argument type %T (expected struct pointer, map pointer, WithDest(...))", dest)
	}

	return optFunc, nil
}

func marshalInput(input any) ([]byte, error) {
	if input == nil {
		return []byte("{}"), nil
	}

	switch v := input.(type) {
	case string:
		if !json.Valid([]byte(v)) {
			return nil, fmt.Errorf("provided string is not valid JSON")
		}
		return []byte(v), nil
	case []byte:
		if !json.Valid(v) {
			return nil, fmt.Errorf("provided byte slice is not valid JSON")
		}
		return v, nil
	}

	val := reflect.ValueOf(input)

	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil, fmt.Errorf("input pointer cannot be nil")
		}
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Struct, reflect.Map:
		return json.Marshal(input)

	case reflect.Slice, reflect.Array:
		return nil, fmt.Errorf(
			"input data cannot be a slice or array; it must be a struct or map to produce a key-value object",
		)

	default:
		return nil, fmt.Errorf(
			"unsupported input data type: %s. Input must be a struct or a map",
			val.Kind(),
		)
	}
}
