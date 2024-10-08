package paseto

import (
	"errors"
	"time"
)

type TokenPayload struct {
	ID        string
	Metadata  map[string]any
	IssuedAt  time.Time
	ExpiredAt time.Time
}

func NewPayload(options ...Option) *TokenPayload {
	var payload TokenPayload
	for _, opt := range options {
		opt(&payload)
	}
	return &payload
}

func (payload *TokenPayload) SetMetadata(key string, value any) {
	payload.Metadata[key] = value
}

func (payload *TokenPayload) GetString(key string) string {
	if value, ok := payload.Metadata[key]; ok {
		return value.(string)
	}
	return ""
}

func (payload *TokenPayload) GetBool(key string) bool {
	if value, ok := payload.Metadata[key]; ok {
		return value.(bool)
	}
	return false
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
