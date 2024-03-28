package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStore_CreateReservationTx(t *testing.T) {
	reservationData := randomReservationData(t)

	arg := CreateReservationParams{}
	err := arg.Unmarshal(reservationData)
	require.NoError(t, err)

	restriction := createRandomRestriction(t)

	reservation, err := testStore.CreateReservationTx(context.Background(), arg, restriction.ID)
	require.NoError(t, err)
	assert.NotEmpty(t, reservation.ID)
	assert.Equal(t, arg.FirstName, reservation.FirstName)
	assert.Equal(t, arg.LastName, reservation.LastName)
	assert.Equal(t, arg.Email, reservation.Email)
	assert.Equal(t, arg.Phone, reservation.Phone)
	assert.Equal(t, arg.StartDate, reservation.StartDate)
	assert.Equal(t, arg.EndDate, reservation.EndDate)
	assert.Equal(t, arg.RoomID, reservation.RoomID)
	assert.Equal(t, arg.Notes, reservation.Notes)
	assert.WithinDuration(t, time.Now(), reservation.CreatedAt.Time, time.Second)
	assert.WithinDuration(t, time.Now(), reservation.UpdatedAt.Time, time.Second)
}
