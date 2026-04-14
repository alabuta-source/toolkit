package inter

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

type requester struct {
	auth           *auth
	baseURL        string
	timeoutSeconds int
	httpClient     *http.Client
	contaCorrente  string

	token    string
	tokenDue time.Time
}

func newRequester(
	clientID, clientSecret, certPath, keyPath string,
	sandbox bool,
	timeoutSeconds int,
	scope, contaCorrente string,
	httpClientOverride *http.Client,
) (*requester, error) {
	base := URLProduction
	if sandbox {
		base = URLSandbox
	}

	if scope == "" {
		scope = DefaultOAuthScope
	}

	var hc *http.Client
	if httpClientOverride != nil {
		hc = httpClientOverride
	} else {
		cert, err := tls.LoadX509KeyPair(certPath, keyPath)
		if err != nil {
			return nil, fmt.Errorf("inter: load client certificate: %w", err)
		}

		tr := &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates: []tls.Certificate{cert},
				MinVersion:   tls.VersionTLS12,
			},
		}
		hc = &http.Client{
			Timeout:   time.Second * time.Duration(timeoutSeconds),
			Transport: tr,
		}
	}

	a := &auth{
		clientID:     clientID,
		clientSecret: clientSecret,
		scope:        scope,
		baseURL:      base,
		httpClient:   hc,
	}

	return &requester{
		auth:           a,
		baseURL:        base,
		timeoutSeconds: timeoutSeconds,
		httpClient:     hc,
		contaCorrente:  digitsOnly(contaCorrente),
	}, nil
}

func digitsOnly(s string) string {
	if s == "" {
		return ""
	}

	var b strings.Builder
	for _, r := range s {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}

	return b.String()
}

func (r *requester) authenticate() error {
	if r.token != "" && time.Now().Before(r.tokenDue) {
		return nil
	}

	tokenData, err := r.auth.getAccessToken()
	if err != nil {
		return err
	}

	r.token = tokenData.AccessToken
	r.tokenDue = time.Now().Add(time.Duration(tokenData.ExpiresIn) * time.Second)

	return nil
}

func (r *requester) request(endpoint, httpVerb string, requestParams map[string]string, body map[string]interface{}) ([]byte, error) {
	params := cloneStringMap(requestParams)
	route := getRoute(endpoint, params)
	route = appendQuery(route, params)

	if err := r.authenticate(); err != nil {
		return nil, err
	}

	var bodyReader io.Reader
	if body != nil && httpVerb != http.MethodGet {
		buf := new(bytes.Buffer)
		if err := json.NewEncoder(buf).Encode(body); err != nil {
			return nil, fmt.Errorf("inter: encode body: %w", err)
		}

		bodyReader = buf
	}

	req, err := http.NewRequest(httpVerb, strings.TrimRight(r.baseURL, "/")+route, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("inter: build request: %w", err)
	}

	if bodyReader != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+r.token)

	if r.contaCorrente != "" {
		req.Header.Set("x-conta-corrente", r.contaCorrente)
	}

	res, err := r.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("inter: do request: %w", err)
	}

	defer res.Body.Close()

	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("inter: read response: %w", err)
	}

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("inter: %s: %s", res.Status, string(respBody))
	}

	return respBody, nil
}

func cloneStringMap(in map[string]string) map[string]string {
	if in == nil {
		return map[string]string{}
	}

	out := make(map[string]string, len(in))
	for k, v := range in {
		out[k] = v
	}

	return out
}

func getRoute(endpoint string, params map[string]string) string {
	pattern := regexp.MustCompile(`:(\w+)`)
	variables := pattern.FindAllStringSubmatch(endpoint, -1)

	for i := 0; i < len(variables); i++ {
		if value, exists := params[variables[i][1]]; exists {
			endpoint = strings.ReplaceAll(endpoint, variables[i][0], value)
			delete(params, variables[i][1])
		}
	}

	return endpoint
}

func appendQuery(path string, params map[string]string) string {
	if len(params) == 0 {
		return path
	}

	sep := "?"
	if strings.Contains(path, "?") {
		sep = "&"
	}

	var b strings.Builder
	b.WriteString(path)
	first := true

	for key, value := range params {
		if first {
			b.WriteString(sep)
			first = false
			sep = "&"
		} else {
			b.WriteString("&")
		}

		b.WriteString(key)
		b.WriteString("=")
		b.WriteString(url.QueryEscape(value))
	}

	return b.String()
}
