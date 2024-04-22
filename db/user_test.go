package db

import (
	"context"
	"errors"
	"testing"

	"github.com/github-real-lb/bookings-web-app/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func createRandomUser(t *testing.T, password string) User {
	arg := CreateUserParams{
		FirstName:   util.RandomName(),
		LastName:    util.RandomName(),
		Email:       util.RandomEmail(),
		Password:    password,
		AccessLevel: util.RandomInt64(1, 10),
	}

	user, err := testStore.CreateNewUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user.ID)
	require.Equal(t, arg.FirstName, user.FirstName)
	require.Equal(t, arg.LastName, user.LastName)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.AccessLevel, user.AccessLevel)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(arg.Password))
	require.NoError(t, err)

	return user
}

func TestPostgresDBStore_CreateNewUser(t *testing.T) {
	createRandomUser(t, util.RandomPassword())
}

func TestPostgresDBStore_AuthenticateUser(t *testing.T) {
	password := util.RandomPassword()
	user := createRandomUser(t, password)

	t.Run("OK", func(t *testing.T) {
		arg := AuthenticateUserParams{
			Email:    user.Email,
			Password: password,
		}

		result, err := testStore.AuthenticateUser(context.Background(), arg)
		require.NoError(t, err)
		assert.Equal(t, user, result)
	})

	t.Run("Wrong Email", func(t *testing.T) {
		arg := AuthenticateUserParams{
			Email:    util.RandomEmail(),
			Password: password,
		}

		result, err := testStore.AuthenticateUser(context.Background(), arg)
		require.Equal(t, err, errors.New("could not authenticate user"))
		assert.Empty(t, result)
	})

	t.Run("Wrong Password", func(t *testing.T) {
		arg := AuthenticateUserParams{
			Email:    user.Email,
			Password: util.RandomPassword(),
		}

		result, err := testStore.AuthenticateUser(context.Background(), arg)
		require.Equal(t, err, errors.New("could not authenticate user"))
		assert.Empty(t, result)
	})
}
