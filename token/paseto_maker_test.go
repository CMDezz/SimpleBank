package token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/techschool/simplebank/db/util"
)

func TestCreatePasetoTokenOk(t *testing.T) {
	symmetricKey := util.RandomString(32)
	username := util.RandomOwner()
	duration := time.Minute
	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	pasetoMaker, err := NewPasetoMaker(symmetricKey)
	require.NoError(t, err)

	token, err := pasetoMaker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := pasetoMaker.ValidToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}
func TestExpiredPasetoToken(t *testing.T) {
	symmetricKey := util.RandomString(32)
	username := util.RandomOwner()
	duration := time.Minute

	pasetoMaker, err := NewMaker(symmetricKey)
	require.NoError(t, err)

	token, err := pasetoMaker.CreateToken(username, -duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := pasetoMaker.ValidToken(token)
	require.Error(t, err)
	require.Nil(t, payload)

	require.EqualError(t, err, ErrExpiredToken.Error())

}
