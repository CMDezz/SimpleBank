package token

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/require"
	"github.com/techschool/simplebank/util"
)

func TestCreateJWTTokenOk(t *testing.T) {
	secretKey := util.RandomString(32)
	username := util.RandomOwner()
	duration := time.Minute
	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	jwtMarker, err := NewMaker(secretKey)
	require.NoError(t, err)

	token, err := jwtMarker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := jwtMarker.ValidToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredJWTToken(t *testing.T) {
	secretKey := util.RandomString(32)
	username := util.RandomOwner()
	duration := time.Minute

	jwtMarker, err := NewMaker(secretKey)
	require.NoError(t, err)

	token, err := jwtMarker.CreateToken(username, -duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := jwtMarker.ValidToken(token)
	require.Error(t, err)
	require.Nil(t, payload)

	require.EqualError(t, err, ErrExpiredToken.Error())

}

func TestInvalidJWTToken(t *testing.T) {
	username := util.RandomOwner()
	duration := time.Minute

	payload, err := NewPayload(username, duration)
	require.NoError(t, err)

	//create a token without signing method
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	//create maker verify the token
	maker, err := NewMaker(util.RandomString(32))
	require.NoError(t, err)

	payload, err = maker.ValidToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)

}
