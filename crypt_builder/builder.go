package cryptbuilder

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

var (
	DecodeStringError           = errors.New("error decoding the cookie")
	DecryptDataError            = errors.New("error trying open aes Galois Counter Mode")
	InvalidNonceSizeError       = errors.New("invalid nonce Galois Counter Mode size")
	WrapCypherError             = errors.New("error Wrap the cipher block")
	CreateCypherFromSecretError = errors.New("error creating cypher from secret")
)

type EncryptService interface {
	Encrypt(value string) (string, error)
	Decrypt(value string) (string, error)
}

type encryptService struct {
	secret []byte
}

func NewEncryptService(secret []byte) EncryptService {
	return &encryptService{
		secret: secret,
	}
}

func (e encryptService) Encrypt(value string) (string, error) {
	block, err := aes.NewCipher(e.secret)
	if err != nil {
		return "", err
	}

	aesGcm, er := cipher.NewGCM(block)
	if er != nil {
		return "", er
	}
	nonce := make([]byte, aesGcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	encryptedValue := aesGcm.Seal(nonce, nonce, []byte(value), nil)
	return e.writeAndEncodeCookie(string(encryptedValue)), nil
}

func (e encryptService) Decrypt(value string) (string, error) {
	encryptedValue, er := e.readAndDecodeCookie(value)
	if er != nil {
		return "", er
	}

	// // Create a new AES cipher block from the secret key.
	block, err := aes.NewCipher(e.secret)
	if err != nil {
		return "", errors.Join(CreateCypherFromSecretError, err)
	}

	// Wrap the cipher block in Galois Counter Mode.
	aesGcm, cErr := cipher.NewGCM(block)
	if cErr != nil {
		return "", errors.Join(WrapCypherError, cErr)
	}

	// To avoid a potential 'index out of range' panic in the next step
	nonceSize := aesGcm.NonceSize()
	if len(encryptedValue) < nonceSize {
		return "", InvalidNonceSizeError
	}

	//decrypt and authenticate the data
	nonce := encryptedValue[:nonceSize]
	ciphertext := encryptedValue[nonceSize:]
	plaintext, opErr := aesGcm.Open(nil, []byte(nonce), []byte(ciphertext), nil)
	if opErr != nil {
		return "", errors.Join(DecryptDataError, opErr)
	}
	return string(plaintext), nil
}

func (encryptService) writeAndEncodeCookie(encryptedValue string) string {
	value := base64.URLEncoding.EncodeToString([]byte(encryptedValue))
	return value
}

func (encryptService) readAndDecodeCookie(encryptedValue string) (string, error) {
	value, err := base64.URLEncoding.DecodeString(encryptedValue)
	if err != nil {
		return "", errors.Join(DecodeStringError, err)
	}
	return string(value), nil
}
