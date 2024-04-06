package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStore_CreateReservationTx(t *testing.T) {
	t.Run("Test OK", func(t *testing.T) {
		rsrvData := randomReservationData(t)

		arg := CreateReservationParams{}
		err := arg.Unmarshal(rsrvData)
		require.NoError(t, err)

		// execute transaction
		rsrv, err := testStore.CreateReservationTx(context.Background(), arg)

		// testify reservation
		require.NoError(t, err)
		assert.NotEmpty(t, rsrv.ID)
		assert.Equal(t, arg.FirstName, rsrv.FirstName)
		assert.Equal(t, arg.LastName, rsrv.LastName)
		assert.Equal(t, arg.Email, rsrv.Email)
		assert.Equal(t, arg.Phone, rsrv.Phone)
		assert.Equal(t, arg.StartDate, rsrv.StartDate)
		assert.Equal(t, arg.EndDate, rsrv.EndDate)
		assert.Equal(t, arg.RoomID, rsrv.RoomID)
		assert.Equal(t, arg.Notes, rsrv.Notes)
		assert.WithinDuration(t, time.Now(), rsrv.CreatedAt.Time, time.Second)
		assert.True(t, rsrv.CreatedAt.Valid)
		assert.WithinDuration(t, time.Now(), rsrv.UpdatedAt.Time, time.Second)
		assert.True(t, rsrv.UpdatedAt.Valid)

		// get last room restriciton
		rr, err := testStore.GetLastRoomRestriction(context.Background(), rsrv.RoomID)

		// testify room restriction
		require.NoError(t, err)
		assert.WithinDuration(t, rsrv.StartDate.Time, rr.StartDate.Time, time.Second)
		assert.True(t, rr.StartDate.Valid)
		assert.WithinDuration(t, rsrv.EndDate.Time, rr.EndDate.Time, time.Second)
		assert.True(t, rr.EndDate.Valid)
		assert.Equal(t, rsrv.ID, rr.ReservationID.Int64)
		assert.True(t, rr.ReservationID.Valid)
		assert.Equal(t, RestrictionReservation, rr.Restriction)
		assert.WithinDuration(t, time.Now(), rr.CreatedAt.Time, time.Second)
		assert.True(t, rr.CreatedAt.Valid)
		assert.WithinDuration(t, time.Now(), rr.UpdatedAt.Time, time.Second)
		assert.True(t, rr.CreatedAt.Valid)
	})

	t.Run("Test Error", func(t *testing.T) {
		arg := CreateReservationParams{}

		// execute transaction
		rsrv, err := testStore.CreateReservationTx(context.Background(), arg)

		//testify
		require.Error(t, err)
		require.Empty(t, rsrv)
	})
}
