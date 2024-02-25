package paseto

import (
	"errors"
	"time"
)

type TokenPayload struct {
	ID        string                 `json:"id"`
	Metadata  map[string]interface{} `json:"metadata"`
	IssuedAt  time.Time              `json:"issued_at"`
	ExpiredAt time.Time              `json:"expired_at"`
}

func NewPayload(option ...Option) *TokenPayload {
	options := setupOptions(option...)

	payload := &TokenPayload{
		ID:        options.iD,
		Metadata:  options.metadata,
		IssuedAt:  options.issuedAt,
		ExpiredAt: time.Now().Add(options.duration),
	}
	return payload
}

func (payload *TokenPayload) valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return errors.New(formatErr(expiredTokenErr))
	}
	return nil
}
