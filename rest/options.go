package rest

import (
	"errors"
	"net/http"
	"time"
)

type clientReQuestOptions struct {
	body              any
	decode            any
	headers           map[string]string
	timeout           time.Duration
	checkRedirectFunc func(*http.Request, []*http.Request) error
}

type OptionFunc func(*clientReQuestOptions)

type Option interface {
	Apply(*clientReQuestOptions)
}

func defaultCheckRedirect(req *http.Request, via []*http.Request) error {
	if len(via) >= 3 {
		return errors.New("stopped after 3 redirects")
	}
	return nil
}

func (f OptionFunc) Apply(client *clientReQuestOptions) {
	f(client)
}

func RequestWithBody(body any) Option {
	return OptionFunc(func(c *clientReQuestOptions) {
		c.body = body
	})
}

func RequestWithHeaders(headers map[string]string) Option {
	return OptionFunc(func(c *clientReQuestOptions) {
		c.headers = headers
	})
}

func RequestWithDecodeValue(decode any) Option {
	return OptionFunc(func(c *clientReQuestOptions) {
		c.decode = decode
	})
}

func RequestWithTimeout(timeout time.Duration) Option {
	return OptionFunc(func(c *clientReQuestOptions) {
		c.timeout = timeout
	})
}

func RequestWithCheckRedirectFunc(fn func(*http.Request, []*http.Request) error) Option {
	return OptionFunc(func(c *clientReQuestOptions) {
		c.checkRedirectFunc = fn
	})
}
