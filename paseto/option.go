package paseto

import "time"

type Option func(payload *TokenPayload)

func WithID(id string) Option {
	return func(p *TokenPayload) {
		p.ID = id
	}
}

func WithMetadata(data map[string]interface{}) Option {
	return func(p *TokenPayload) {
		p.Metadata = data
	}
}

func WithIssueDate(date time.Time) Option {
	return func(p *TokenPayload) {
		p.IssuedAt = date
	}
}

func WithDuration(duration time.Duration) Option {
	return func(p *TokenPayload) {
		p.ExpiredAt = time.Now().Add(duration)
	}
}
