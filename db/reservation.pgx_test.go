package db

import (
	"fmt"
	"testing"
	"time"

	"github.com/github-real-lb/bookings-web-app/util"
	"github.com/github-real-lb/bookings-web-app/util/config"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateReservationParams_Unmarshal(t *testing.T) {
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
			Time:  util.RandomDate(),
			Valid: true,
		},
		EndDate: pgtype.Date{
			Time:  util.RandomDate(),
			Valid: true,
		},
		RoomID: util.RandomID(),
		Notes: pgtype.Text{
			String: util.RandomNote(),
			Valid:  true,
		},
	}

	data := map[string]string{
		"code":       arg.Code,
		"first_name": arg.FirstName,
		"last_name":  arg.LastName,
		"email":      arg.Email,
		"phone":      arg.Phone.String,
		"start_date": arg.StartDate.Time.Format(config.DateLayout),
		"end_date":   arg.EndDate.Time.Format(config.DateLayout),
		"room_id":    fmt.Sprint(arg.RoomID),
		"notes":      arg.Notes.String,
	}

	arg2 := CreateReservationParams{}
	err := arg2.Unmarshal(data)
	require.NoError(t, err)
	assert.Equal(t, arg, arg2)

}

func TestUpdateReservationParams_Unmarshal(t *testing.T) {
	arg := UpdateReservationParams{
		ID:        util.RandomID(),
		Code:      util.RandomString(ReservationCodeLenght),
		FirstName: util.RandomName(),
		LastName:  util.RandomName(),
		Email:     util.RandomEmail(),
		Phone: pgtype.Text{
			String: util.RandomPhone(),
			Valid:  true,
		},
		StartDate: pgtype.Date{
			Time:  util.RandomDate(),
			Valid: true,
		},
		EndDate: pgtype.Date{
			Time:  util.RandomDate(),
			Valid: true,
		},
		RoomID: util.RandomID(),
		Notes: pgtype.Text{
			String: util.RandomNote(),
			Valid:  true,
		},
		UpdatedAt: pgtype.Timestamptz{
			Time:  util.RandomDatetime(),
			Valid: true,
		},
	}

	data := map[string]string{
		"id":         fmt.Sprint(arg.ID),
		"code":       arg.Code,
		"first_name": arg.FirstName,
		"last_name":  arg.LastName,
		"email":      arg.Email,
		"phone":      arg.Phone.String,
		"start_date": arg.StartDate.Time.Format(config.DateLayout),
		"end_date":   arg.EndDate.Time.Format(config.DateLayout),
		"room_id":    fmt.Sprint(arg.RoomID),
		"notes":      arg.Notes.String,
		"updated_at": arg.UpdatedAt.Time.Format(config.DateTimeLayout),
	}

	arg2 := UpdateReservationParams{}
	err := arg2.Unmarshal(data)

	require.NoError(t, err)
	assert.Equal(t, arg.ID, arg2.ID)
	assert.Equal(t, arg.Code, arg2.Code)
	assert.Equal(t, arg.FirstName, arg2.FirstName)
	assert.Equal(t, arg.LastName, arg2.LastName)
	assert.Equal(t, arg.Email, arg2.Email)
	assert.Equal(t, arg.Phone, arg2.Phone)
	assert.Equal(t, arg.StartDate, arg2.StartDate)
	assert.Equal(t, arg.EndDate, arg2.EndDate)
	assert.Equal(t, arg.RoomID, arg2.RoomID)
	assert.Equal(t, arg.Notes, arg2.Notes)
	assert.WithinDuration(t, arg.UpdatedAt.Time, arg2.UpdatedAt.Time, time.Second)
}
