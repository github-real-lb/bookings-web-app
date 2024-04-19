package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

type AuthenticateUserParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Unmarshal parse data into p
func (p *AuthenticateUserParams) Unmarshal(data map[string]string) {
	p.Email = data["email"]
	p.Password = data["password"]
}

// AuthenticateUser validates the email and password of a user.
// Returns nil on success, or an error on failure.
func (store *PostgresDBStore) AuthenticateUser(ctx context.Context, arg AuthenticateUserParams) (User, error) {

	user, err := store.GetUserByEmail(ctx, arg.Email)
	if err != nil {
		return user, err
	}

	if user.ID == 0 {
		return User{}, errors.New("could not authenticate user")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(arg.Password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return User{}, errors.New("could not authenticate user")
	} else if err != nil {
		return User{}, err
	}

	return user, nil
}

func (store *PostgresDBStore) CreateReservationTx(ctx context.Context, arg CreateReservationParams) (Reservation, error) {
	var reservation Reservation

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		// insert new reservation into database
		reservation, err = store.CreateReservation(ctx, arg)
		if err != nil {
			return err
		}

		rrArg := CreateRoomRestrictionParams{
			StartDate: reservation.StartDate,
			EndDate:   reservation.EndDate,
			RoomID:    reservation.RoomID,
			ReservationID: pgtype.Int8{
				Int64: reservation.ID,
				Valid: true,
			},
			Restriction: RestrictionReservation,
		}

		_, err = store.CreateRoomRestriction(ctx, rrArg)
		if err != nil {
			return err
		}

		return nil
	})

	return reservation, err
}
