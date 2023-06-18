package toolkit

import (
	"errors"
	"github.com/o1egl/paseto"
	"time"
)

var (
	expiredTokenErr = errors.New("token has expired")
	invalidTokenErr = errors.New("token is invalid")
)

type PasetoTokenBuilder interface {
	// CreateToken creates a new token for a specific username and duration,
	// and returns the signed token string or an error.
	// The tokenID is used to identify the token you can send a UUID or the userID.
	CreateToken(tokenID, username string, duration time.Duration) (string, error)
	// VerifyToken verifies the token string and returns the payload or an error.
	VerifyToken(token string) (*TokenPayload, error)
}

type pasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

// NewPasetoMaker creates a new PasetoTokenBuilder.
// The symmetricKey is used to sign the token and needs to be a 32 len string
func NewPasetoMaker(symmetricKey string) PasetoTokenBuilder {
	maker := &pasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}
	return maker
}

func (maker *pasetoMaker) CreateToken(tokenID, username string, duration time.Duration) (string, error) {
	payload := NewPayload(tokenID, username, duration)
	return maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
}

func (maker *pasetoMaker) VerifyToken(token string) (*TokenPayload, error) {
	payload := new(TokenPayload)

	err := maker.paseto.Decrypt(token, maker.symmetricKey, payload, nil)
	if err != nil {
		return nil, invalidTokenErr
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}
	return payload, nil
}
