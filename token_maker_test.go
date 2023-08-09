package toolkit

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestPasetoMaker(t *testing.T) {
	maker, pErr := NewTokenMaker(randomString(32))
	require.NoError(t, pErr)

	username := randomOwner()
	duration := time.Minute
	tokenID := GenerateUUID()

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, err := maker.CreateToken(tokenID, username, duration)
	require.NotEmpty(t, token)
	require.NoError(t, err)

	payload, er := maker.VerifyToken(token)
	require.NoError(t, er)
	require.NotEmpty(t, token)

	require.NotZero(t, payload)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	maker, pErr := NewTokenMaker(randomString(32))
	require.NoError(t, pErr)
	tokenID := GenerateUUID()

	token, err := maker.CreateToken(tokenID, randomOwner(), -time.Minute)
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
