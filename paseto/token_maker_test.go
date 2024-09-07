package paseto

import (
	"testing"
	"time"

	"github.com/alabuta-source/toolkit"
	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	maker, pErr := NewTokenMaker(randomString(32))
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

func TestExpiredPasetoToken(t *testing.T) {
	maker, pErr := NewTokenMaker(randomString(32))
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
	_, err := NewTokenMaker(randomString(1))
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid key size")
}
