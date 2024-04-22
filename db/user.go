package db

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func (store *PostgresDBStore) CreateNewUser(ctx context.Context, arg CreateUserParams) (User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(arg.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}

	arg.Password = string(hash)

	return store.CreateUser(ctx, arg)
}

type AuthenticateUserParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthenticateUser validates the email and password of a user.
// Returns nil on success, or an error on failure.
func (store *PostgresDBStore) AuthenticateUser(ctx context.Context, arg AuthenticateUserParams) (User, error) {
	// get user from database using email
	user, err := store.GetUserByEmail(ctx, arg.Email)
	if err != nil {
		return user, errors.New("could not authenticate user")
	}

	// compate database hash to passed password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(arg.Password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return User{}, errors.New("could not authenticate user")
	} else if err != nil {
		return User{}, err
	}

	return user, nil
}
