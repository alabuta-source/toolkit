package nosql_document

import "time"

type documentOptions struct {
	timeout time.Duration
	appName string
}

type Option func(*documentOptions)

func WithAppName(appName string) Option {
	return func(o *documentOptions) {
		o.appName = appName
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(o *documentOptions) {
		o.timeout = timeout
	}
}
