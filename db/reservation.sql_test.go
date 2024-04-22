package db

import (
	"context"
	"testing"
	"time"

	"github.com/github-real-lb/bookings-web-app/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const ReservationCodeLenght = 14

func createRandomReservation(t *testing.T, room Room) Reservation {
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

	r, err := testStore.CreateReservation(context.Background(), arg)
	require.NoError(t, err)
	assert.NotEmpty(t, r.ID)
	assert.Equal(t, arg.FirstName, r.FirstName)
	assert.Equal(t, arg.LastName, r.LastName)
	assert.Equal(t, arg.Email, r.Email)
	assert.Equal(t, arg.Phone, r.Phone)
	assert.Equal(t, arg.StartDate, r.StartDate)
	assert.Equal(t, arg.EndDate, r.EndDate)
	assert.Equal(t, arg.RoomID, r.RoomID)
	assert.Equal(t, arg.Notes, r.Notes)
	assert.WithinDuration(t, time.Now(), r.CreatedAt.Time, time.Second)
	assert.WithinDuration(t, time.Now(), r.UpdatedAt.Time, time.Second)

	return r
}

func createRandomWeekReservation(t *testing.T, room Room, startDate time.Time) Reservation {
	arg := CreateReservationParams{
		Code:      util.RandomString(ReservationCodeLenght),
		FirstName: util.RandomName(),
		LastName:  util.RandomName(),
		Email:     util.RandomEmail(),
		RoomID:    room.ID,
	}
	arg.Phone.Scan(util.RandomPhone())
	arg.StartDate.Scan(startDate)
	arg.EndDate.Scan(startDate.Add(time.Hour * 24 * 7))
	arg.Notes.Scan(util.RandomNote())

	r, err := testStore.CreateReservation(context.Background(), arg)
	require.NoError(t, err)
	assert.NotEmpty(t, r.ID)
	assert.Equal(t, arg.FirstName, r.FirstName)
	assert.Equal(t, arg.LastName, r.LastName)
	assert.Equal(t, arg.Email, r.Email)
	assert.Equal(t, arg.Phone, r.Phone)
	assert.Equal(t, arg.StartDate, r.StartDate)
	assert.Equal(t, arg.EndDate, r.EndDate)
	assert.Equal(t, arg.RoomID, r.RoomID)
	assert.Equal(t, arg.Notes, r.Notes)
	assert.WithinDuration(t, time.Now(), r.CreatedAt.Time, time.Second)
	assert.WithinDuration(t, time.Now(), r.UpdatedAt.Time, time.Second)

	return r
}

func TestQueries_CreateReservation(t *testing.T) {
	room := createRandomRoom(t)
	createRandomReservation(t, room)
}
