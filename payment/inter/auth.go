package inter

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type auth struct {
	clientID     string
	clientSecret string
	scope        string
	baseURL      string
	httpClient   *http.Client
}

type authResponseBody struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

func (a *auth) getAccessToken() (authResponseBody, error) {
	var out authResponseBody
	form := url.Values{}
	form.Set("client_id", a.clientID)
	form.Set("client_secret", a.clientSecret)
	form.Set("grant_type", "client_credentials")
	form.Set("scope", a.scope)

	req, err := http.NewRequest(http.MethodPost, strings.TrimRight(a.baseURL, "/")+"/oauth/v2/token", strings.NewReader(form.Encode()))
	if err != nil {
		return out, fmt.Errorf("inter auth: build request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	res, err := a.httpClient.Do(req)
	if err != nil {
		return out, fmt.Errorf("inter auth: do request: %w", err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return out, fmt.Errorf("inter auth: read body: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return out, fmt.Errorf("inter auth: %s: %s", res.Status, string(body))
	}

	if err := json.Unmarshal(body, &out); err != nil {
		return out, fmt.Errorf("inter auth: decode token: %w", err)
	}

	if out.AccessToken == "" {
		return out, fmt.Errorf("inter auth: empty access_token in response")
	}

	if out.ExpiresIn <= 0 {
		out.ExpiresIn = 3600
	}

	return out, nil
}
