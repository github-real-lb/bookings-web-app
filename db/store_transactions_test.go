package db

import (
	"context"
	"testing"
	"time"

	"github.com/github-real-lb/bookings-web-app/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStore_CreateReservationTx(t *testing.T) {
	t.Run("Test OK", func(t *testing.T) {
		room := createRandomRoom(t)

		rDate := util.RandomDate()
		arg := CreateReservationParams{
			Code:      util.RandomString(ReservationCodeLenght),
			FirstName: util.RandomName(),
			LastName:  util.RandomName(),
			Email:     util.RandomEmail(),
			RoomID:    room.ID,
		}
		arg.Phone.Scan(util.RandomPhone())
		arg.StartDate.Scan(rDate)
		arg.EndDate.Scan(rDate.Add(time.Hour * 24 * 7))
		arg.Notes.Scan(util.RandomNote())

		// execute transaction
		rsv, err := testStore.CreateReservationTx(context.Background(), arg)

		// testify reservation
		require.NoError(t, err)
		assert.NotEmpty(t, rsv.ID)
		assert.Equal(t, arg.FirstName, rsv.FirstName)
		assert.Equal(t, arg.LastName, rsv.LastName)
		assert.Equal(t, arg.Email, rsv.Email)
		assert.Equal(t, arg.Phone, rsv.Phone)
		assert.Equal(t, arg.StartDate, rsv.StartDate)
		assert.Equal(t, arg.EndDate, rsv.EndDate)
		assert.Equal(t, arg.RoomID, rsv.RoomID)
		assert.Equal(t, arg.Notes, rsv.Notes)
		assert.WithinDuration(t, time.Now(), rsv.CreatedAt.Time, time.Second)
		assert.True(t, rsv.CreatedAt.Valid)
		assert.WithinDuration(t, time.Now(), rsv.UpdatedAt.Time, time.Second)
		assert.True(t, rsv.UpdatedAt.Valid)

		// get last room restriciton
		rr, err := testStore.GetLastRoomRestriction(context.Background(), rsv.RoomID)

		// testify room restriction
		require.NoError(t, err)
		assert.WithinDuration(t, rsv.StartDate.Time, rr.StartDate.Time, time.Second)
		assert.True(t, rr.StartDate.Valid)
		assert.WithinDuration(t, rsv.EndDate.Time, rr.EndDate.Time, time.Second)
		assert.True(t, rr.EndDate.Valid)
		assert.Equal(t, rsv.ID, rr.ReservationID.Int64)
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
		rsv, err := testStore.CreateReservationTx(context.Background(), arg)

		//testify
		require.Error(t, err)
		require.Empty(t, rsv)
	})
}
