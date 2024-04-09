package db

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// createRandomRoomRestriction creates a random room restriction in the database for a reservation
func createRandomRoomRestriction(t *testing.T, r Reservation) RoomRestriction {
	arg := CreateRoomRestrictionParams{
		StartDate: r.StartDate,
		EndDate:   r.EndDate,
		RoomID:    r.RoomID,
		ReservationID: pgtype.Int8{
			Int64: r.ID,
			Valid: true,
		},
		Restriction: RestrictionReservation,
	}

	rr, err := testStore.CreateRoomRestriction(context.Background(), arg)
	require.NoError(t, err)
	assert.NotEmpty(t, rr.ID)

	assert.Equal(t, arg.StartDate, rr.StartDate)
	assert.Equal(t, arg.EndDate, rr.EndDate)

	assert.Equal(t, arg.RoomID, rr.RoomID)
	assert.Equal(t, arg.ReservationID, rr.ReservationID)
	assert.Equal(t, arg.Restriction, rr.Restriction)

	assert.WithinDuration(t, time.Now(), rr.CreatedAt.Time, time.Second)
	assert.True(t, rr.CreatedAt.Valid)
	assert.WithinDuration(t, time.Now(), rr.UpdatedAt.Time, time.Second)
	assert.True(t, rr.UpdatedAt.Valid)

	return rr
}

func TestQueries_CreateRoomRestriction(t *testing.T) {
	room := createRandomRoom(t)
	reservation := createRandomReservation(t, room)
	createRandomRoomRestriction(t, reservation)
}
