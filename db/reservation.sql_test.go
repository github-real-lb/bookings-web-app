package db

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/github-real-lb/bookings-web-app/util"
	"github.com/github-real-lb/bookings-web-app/util/config"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const ReservationCodeLenght = 14

func randomReservationData(t *testing.T) map[string]string {

	startDate := util.RandomDatetime()
	room := createRandomRoom(t)

	data := make(map[string]string)
	data["code"] = util.RandomString(ReservationCodeLenght)
	data["first_name"] = util.RandomName()
	data["last_name"] = util.RandomName()
	data["email"] = util.RandomEmail()
	data["phone"] = util.RandomPhone()
	data["start_date"] = startDate.Format(config.DateLayout)
	data["end_date"] = startDate.Add(time.Hour * 24 * 7).Format(config.DateLayout)
	data["room_id"] = fmt.Sprint(room.ID)
	data["notes"] = util.RandomNote()
	return data
}

func createRandomReservation(t *testing.T, room Room) Reservation {
	startDate := util.RandomDate()

	arg := CreateReservationParams{
		Code:      util.RandomString(ReservationCodeLenght),
		FirstName: util.RandomName(),
		LastName:  util.RandomName(),
		Email:     util.RandomEmail(),
		Phone: pgtype.Text{
			String: util.RandomPhone(),
			Valid:  true,
		},
		StartDate: pgtype.Date{
			Time:  startDate,
			Valid: true,
		},
		EndDate: pgtype.Date{
			Time:  startDate.Add(time.Hour * 24 * 7),
			Valid: true,
		},
		RoomID: room.ID,
		Notes: pgtype.Text{
			String: util.RandomNote(),
			Valid:  true,
		},
	}

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
		Phone: pgtype.Text{
			String: util.RandomPhone(),
			Valid:  true,
		},
		StartDate: pgtype.Date{
			Time:  startDate,
			Valid: true,
		},
		EndDate: pgtype.Date{
			Time:  startDate.Add(time.Hour * 24 * 7),
			Valid: true,
		},
		RoomID: room.ID,
		Notes: pgtype.Text{
			String: util.RandomNote(),
			Valid:  true,
		},
	}

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
