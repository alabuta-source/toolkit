package paseto

import "time"

type payloadOption struct {
	iD       string
	metadata map[string]interface{}
	issuedAt time.Time
	duration time.Duration
}

type Option func(*payloadOption)

func WithID(id string) Option {
	return func(p *payloadOption) {
		p.iD = id
	}
}

func WithMetadata(data map[string]interface{}) Option {
	return func(p *payloadOption) {
		p.metadata = data
	}
}

func WithIssueDate(date time.Time) Option {
	return func(p *payloadOption) {
		p.issuedAt = date
	}
}

func WithDuration(date time.Duration) Option {
	return func(p *payloadOption) {
		p.duration = date
	}
}
