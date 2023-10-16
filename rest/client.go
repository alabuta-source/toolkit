package rest

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Restfull interface {
	POST(url string, body interface{}, headers map[string]string, decode interface{}) error
	GET(url string, headers map[string]string, decode interface{}) error
	PUT(url string, body interface{}, headers map[string]string, decode interface{}) error
}

type restfullClient struct {
	client *http.Client
}

func NewRestClient(timeout int) Restfull {
	return &restfullClient{
		client: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
	}
}

func (client *restfullClient) buildRequestWithBody(method string, url string, buffer bytes.Buffer, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequest(method, url, &buffer)
	if err != nil {
		return nil, fmt.Errorf("error build %s rest request, message: %v", method, err)
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}
	return client.doRequest(req)
}

func (client *restfullClient) doRequest(req *http.Request) ([]byte, error) {
	resp, er := client.client.Do(req)
	if er != nil {
		return nil, fmt.Errorf(
			"error doing client request, message: %s, "+
				"url: %s",
			er.Error(), req.URL.Path,
		)
	}
	defer resp.Body.Close()

	bts, ioErr := io.ReadAll(resp.Body)
	if ioErr != nil {
		return nil, fmt.Errorf("error reading body response")
	}
	return bts, nil
}
