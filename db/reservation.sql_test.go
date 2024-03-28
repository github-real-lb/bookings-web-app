package db

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/github-real-lb/bookings-web-app/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const ReservationCodeLenght = 6

func randomReservationData(t *testing.T) StringMap {
	name := util.RandomName()

	code, err := util.GenerateReservationCode(name, ReservationCodeLenght)
	require.NoError(t, err)

	startDate := util.RandomDatetime()

	room := createRandomRoom(t)

	data := make(StringMap)
	data["code"] = code
	data["first_name"] = util.RandomName()
	data["last_name"] = name
	data["email"] = util.RandomEmail()
	data["phone"] = util.RandomPhoneNumber()
	data["start_date"] = startDate.Format("2006-01-02")
	data["end_date"] = startDate.Add(time.Hour * 24 * 7).Format("2006-01-02")
	data["room_id"] = fmt.Sprint(room.ID)
	data["notes"] = util.RandomNote()
	return data
}

func createRandomReservation(t *testing.T) Reservation {
	data := randomReservationData(t)

	arg := CreateReservationParams{}
	err := arg.Unmarshal(data)
	require.NoError(t, err)

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
	createRandomReservation(t)
}
