package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

func (store *PostgresDBStore) CreateReservationTx(ctx context.Context, arg CreateReservationParams, restrictionsID int64) (Reservation, error) {
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
			RestrictionsID: restrictionsID,
		}

		_, err = store.CreateRoomRestriction(ctx, rrArg)
		if err != nil {
			return err
		}

		return nil
	})

	return reservation, err
}
