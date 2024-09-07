package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

const defaultTimeout = 5 * time.Second

type RestClient struct {
	httpClient     *http.Client
	requestOptions *clientReQuestOptions
	url            string
	method         string
}

func NewHttpClient() *RestClient {
	return &RestClient{}
}

func (client *RestClient) BuildRequest(url, method string, options ...Option) *RestClient {
	var requestOptions clientReQuestOptions
	for _, o := range options {
		o.Apply(&requestOptions)
	}

	var timeout = defaultTimeout
	if requestOptions.timeout.String() != "0s" {
		timeout = requestOptions.timeout
	}

	checkRedirectFunc := requestOptions.checkRedirectFunc
	if checkRedirectFunc == nil {
		checkRedirectFunc = defaultCheckRedirect
	}

	client.httpClient = &http.Client{Timeout: timeout, CheckRedirect: checkRedirectFunc}
	client.requestOptions = &requestOptions
	client.method = method
	client.url = url

	return client
}

func (client *RestClient) Execute() error {
	var buf bytes.Buffer
	if client.requestOptions.body != nil {
		if err := json.NewEncoder(&buf).Encode(client.requestOptions.body); err != nil {
			return fmt.Errorf("[restClient] error encode request body: %v", err)
		}
	}

	request, err := http.NewRequest(client.method, client.url, &buf)
	if err != nil {
		return fmt.Errorf(
			"error trying to build %s request, message: %v",
			client.method,
			err,
		)
	}

	request.Header.Set("content-type", "application/json")
	for k, v := range client.requestOptions.headers {
		request.Header.Set(k, v)
	}

	data, er := client.doRequest(request)
	if er != nil {
		return er
	}

	if client.requestOptions.decode != nil {
		if err = json.Unmarshal(data, client.requestOptions.decode); err != nil {
			return fmt.Errorf("error trying to Unmarshal response: %v", err)
		}
	}
	return nil
}

func (client *RestClient) doRequest(req *http.Request) ([]byte, error) {
	resp, er := client.httpClient.Do(req)
	if er != nil {
		return nil, fmt.Errorf(
			"error doing the request, message: %s, "+
				"url: %s",
			er.Error(), req.URL.Path,
		)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		bytes, err := client.closeBodyAndSendResponse(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading the body: %v", err)
		}
		return nil, errors.New(string(bytes))
	}
	return client.closeBodyAndSendResponse(resp.Body)
}

func (client *RestClient) closeBodyAndSendResponse(body io.ReadCloser) ([]byte, error) {
	bts, ioErr := io.ReadAll(body)
	if ioErr != nil {
		return nil, ioErr
	}
	return bts, nil
}
