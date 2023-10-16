package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
)

func (client *restfullClient) POST(url string, body interface{}, headers map[string]string, decode interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(body); err != nil {
		return fmt.Errorf("error enconde body, message: %v", err)
	}

	resp, er := client.buildRequestWithBody(http.MethodPost, url, buf, headers)
	if er != nil {
		return er
	}

	if reflect.ValueOf(decode).IsNil() {
		return nil
	}
	return json.Unmarshal(resp, decode)
}

func (client *restfullClient) GET(url string, headers map[string]string, decode interface{}) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("error build POST rest request, message: %v", err)
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	bts, er := client.doRequest(req)
	if er != nil {
		return fmt.Errorf("error doing %s request: %v", req.Method, er)
	}

	if reflect.ValueOf(decode).IsNil() {
		return nil
	}
	return json.Unmarshal(bts, decode)
}

func (client *restfullClient) PUT(url string, body interface{}, headers map[string]string, decode interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(body); err != nil {
		return fmt.Errorf("error enconde body, message: %v", err)
	}

	resp, er := client.buildRequestWithBody(http.MethodPut, url, buf, headers)
	if er != nil {
		return er
	}

	if reflect.ValueOf(decode).IsNil() {
		return nil
	}
	return json.Unmarshal(resp, decode)
}
