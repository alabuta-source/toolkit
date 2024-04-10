package paseto

import (
	"errors"
	"fmt"
	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

var (
	expiredTokenErr = "token has expired: [%s]"
	invalidTokenErr = "token is invalid: [%s]"
)

type TokenBuilder interface {
	// CreateToken creates a new token for a specific username and duration,
	// and returns the signed token string or an error.
	// The tokenID is used to identify the token you can send a UUID or the userID.
	CreateToken(option ...Option) (string, error)
	// VerifyToken verifies the token string and returns the payload or an error.
	VerifyToken(token string) (*TokenPayload, error)
}

type pasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

// NewTokenMaker creates a new TokenBuilder.
// The symmetricKey is used to sign the token and needs to be a 32 len string
func NewTokenMaker(symmetricKey string) (TokenBuilder, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be exactly %d characters", chacha20poly1305.KeySize)
	}
	maker := &pasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}
	return maker, nil
}

func (maker *pasetoMaker) CreateToken(option ...Option) (string, error) {
	payload := NewPayload(option...)
	return maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
}

func (maker *pasetoMaker) VerifyToken(token string) (*TokenPayload, error) {
	var payload TokenPayload
	err := maker.paseto.Decrypt(token, maker.symmetricKey, &payload, nil)
	if err != nil {
		return nil, errors.New(formatErr(invalidTokenErr, err.Error()))
	}

	err = payload.valid()
	if err != nil {
		return nil, err
	}
	return &payload, nil
}
