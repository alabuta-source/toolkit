package toolkit

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestPasetoMaker(t *testing.T) {
	maker := NewPasetoMaker(randomString(32))

	username := randomOwner()
	duration := time.Minute
	tokenID := generateUUID()

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
	maker := NewPasetoMaker(randomString(32))
	tokenID := generateUUID()

	token, err := maker.CreateToken(tokenID, randomOwner(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, er := maker.VerifyToken(token)
	require.Error(t, er)
	require.EqualError(t, er, expiredTokenErr.Error())
	require.Nil(t, payload)
}
