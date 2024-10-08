package paseto

import (
	"crypto/ed25519"
	"errors"
	"fmt"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

var (
	expiredTokenErr = "token has expired"
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
	publicMode   bool
	publicKey    ed25519.PublicKey
	privatekey   ed25519.PrivateKey
	symmetricKey []byte
	footer       map[string]string
}

// NewTokenMaker creates a new TokenBuilder.
// The symmetricKey is used to sign the token and needs to be a 32 len string
func NewTokenMaker(
	symmetricKey string,
	publicMode bool,
	publicKey ed25519.PublicKey,
	privatekey ed25519.PrivateKey,
) (TokenBuilder, error) {
	if publicMode {
		if privatekey == nil || publicKey == nil {
			return nil, errors.New("private/public key should not be null")
		}
	} else {
		if len(symmetricKey) != chacha20poly1305.KeySize {
			return nil, fmt.Errorf("invalid key size: must be exactly %d characters", chacha20poly1305.KeySize)
		}
	}

	maker := &pasetoMaker{
		paseto:       paseto.NewV2(),
		publicMode:   publicMode,
		publicKey:    publicKey,
		privatekey:   privatekey,
		symmetricKey: []byte(symmetricKey),
		footer:       map[string]string{"version": "2", "app": "alabuta-toolkit"},
	}
	return maker, nil
}

func (maker *pasetoMaker) CreateToken(option ...Option) (string, error) {
	payload := NewPayload(option...)

	if maker.publicMode {
		return maker.paseto.Sign(maker.privatekey, payload, &maker.footer)
	}
	return maker.paseto.Encrypt(maker.symmetricKey, payload, &maker.footer)
}

func (maker *pasetoMaker) VerifyToken(token string) (*TokenPayload, error) {
	var payload TokenPayload
	var err error

	if maker.publicMode {
		err = maker.paseto.Verify(token, maker.publicKey, &payload, &maker.footer)
	} else {
		err = maker.paseto.Decrypt(token, maker.symmetricKey, &payload, &maker.footer)
	}

	if err != nil {
		return nil, errors.New(formatErr(invalidTokenErr, err.Error()))
	}

	err = payload.valid()
	if err != nil {
		return nil, err
	}
	return &payload, nil
}
