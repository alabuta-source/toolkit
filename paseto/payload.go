package paseto

import (
	"errors"
	"time"
)

type TokenPayload struct {
	ID        string         `json:"id"`
	Metadata  map[string]any `json:"metadata"`
	IssuedAt  time.Time      `json:"issued_at"`
	ExpiredAt time.Time      `json:"expired_at"`
}

func NewPayload(options ...Option) *TokenPayload {
	var payload TokenPayload
	for _, opt := range options {
		opt(&payload)
	}
	return &payload
}

func (payload *TokenPayload) GetString(key string) string {
	if value, ok := payload.Metadata[key]; ok {
		return value.(string)
	}
	return ""
}

func (payload *TokenPayload) GetData(key string) any {
	if value, ok := payload.Metadata[key]; ok {
		return value
	}
	return nil
}

func (payload *TokenPayload) valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return errors.New(expiredTokenErr)
	}
	return nil
}
