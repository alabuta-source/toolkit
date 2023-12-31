package paseto

import (
	"errors"
	"time"
)

type TokenPayload struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func NewPayload(tokenID, username string, duration time.Duration) *TokenPayload {
	payload := &TokenPayload{
		ID:        tokenID,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	return payload
}

func (payload *TokenPayload) valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return errors.New(formatErr(expiredTokenErr))
	}
	return nil
}
