package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/techschool/simplebank/util"
)

func CreateRandomTestUser(t *testing.T) User {
	return createRandomTestUser(t)
}

func createRandomTestUser(t *testing.T) User {
	hased_password, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	arg := CreateUserParams{
		Username:      util.RandomOwner(),
		HasedPassword: hased_password,
		FullName:      util.RandomOwner(),
		Email:         util.RandomEmail(),
	}

	user, error := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, error)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HasedPassword, user.HasedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomTestUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomTestUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.HasedPassword, user2.HasedPassword)

	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
}

func TestUpdateUser(t *testing.T) {
	oldUser := CreateRandomTestUser(t)

	newFullname := util.RandomOwner()

	newUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		FullName: sql.NullString{
			String: newFullname,
			Valid:  true,
		},
	})
	require.NoError(t, err)
	require.NotEmpty(t, newUser)
	require.Equal(t, newFullname, newUser.FullName)
	require.Equal(t, oldUser.Email, newUser.Email)
	require.Equal(t, oldUser.HasedPassword, newUser.HasedPassword)
}
