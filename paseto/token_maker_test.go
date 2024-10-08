package paseto

import (
	"crypto/ed25519"
	"encoding/hex"
	"testing"
	"time"

	"github.com/alabuta-source/toolkit"
	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	maker, pErr := NewTokenMaker(randomString(32), false, nil, nil)
	require.NoError(t, pErr)

	username := randomOwner()
	duration := time.Minute
	tokenID := toolkit.GenerateID()

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	metaData := make(map[string]interface{})
	metaData["username"] = username

	token, err := maker.CreateToken(
		WithID(tokenID),
		WithIssueDate(issuedAt),
		WithMetadata(metaData),
		WithDuration(duration),
	)
	require.NotEmpty(t, token)
	require.NoError(t, err)

	payload, er := maker.VerifyToken(token)
	require.NoError(t, er)
	require.NotEmpty(t, token)

	require.NotZero(t, payload)
	require.Equal(t, username, payload.GetString("username"))
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestPublicMode_error_cases(t *testing.T) {
	maker, err := NewTokenMaker("", true, nil, nil)

	require.Nil(t, maker)
	require.Error(t, err)
	require.Equal(t, "private/public key should not be null", err.Error())
}

func TestPublicMode(t *testing.T) {
	privateKeyBytes, er := hex.DecodeString("b4cbfb43df4ce210727d953e4a713307fa19bb7d9f85041438d9e11b942a37741eb9dbbbbc047c03fd70604e0071f0987e16b28b757225c11f00415d0e20b1a2")
	require.NoError(t, er)

	publicKeyBytes, err := hex.DecodeString("1eb9dbbbbc047c03fd70604e0071f0987e16b28b757225c11f00415d0e20b1a2")
	require.NoError(t, err)

	privateKey := ed25519.PrivateKey(privateKeyBytes)
	publicKey := ed25519.PublicKey(publicKeyBytes)

	maker, tErr := NewTokenMaker("", true, publicKey, privateKey)
	require.NoError(t, tErr)

	username := randomOwner()
	duration := time.Minute
	tokenID := toolkit.GenerateID()

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	metaData := make(map[string]interface{})
	metaData["username"] = username

	token, err := maker.CreateToken(
		WithID(tokenID),
		WithIssueDate(issuedAt),
		WithMetadata(metaData),
		WithDuration(duration),
	)
	require.NotEmpty(t, token)
	require.NoError(t, err)

	payload, er := maker.VerifyToken(token)
	require.NoError(t, er)
	require.NotEmpty(t, token)

	require.NotZero(t, payload)
	require.Equal(t, username, payload.GetString("username"))
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	maker, pErr := NewTokenMaker(randomString(32), false, nil, nil)
	require.NoError(t, pErr)
	tokenID := toolkit.GenerateID()

	token, err := maker.CreateToken(
		WithID(tokenID),
		WithDuration(-time.Minute),
	)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, er := maker.VerifyToken(token)
	require.Error(t, er)
	require.Contains(t, er.Error(), "token has expired")
	require.Nil(t, payload)
}

func TestInvalidSymmetricKey(t *testing.T) {
	_, err := NewTokenMaker(randomString(1), false, nil, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid key size")
}
