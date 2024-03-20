package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	password := RandomString(6)

	hased_password1, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hased_password1)

	err = CheckPassword(password, hased_password1)
	require.NoError(t, err)

	wrongPassword := RandomString(6)
	err = CheckPassword(wrongPassword, hased_password1)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	hased_password2, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hased_password2)

	require.NotEqual(t, hased_password1, hased_password2)
}
