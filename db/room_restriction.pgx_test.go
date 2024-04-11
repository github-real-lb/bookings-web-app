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

func TestCreateRoomRestrictionParams_Unmarshal(t *testing.T) {
	arg := CreateRoomRestrictionParams{
		StartDate: pgtype.Date{
			Time:  util.RandomDate(),
			Valid: true,
		},
		EndDate: pgtype.Date{
			Time:  util.RandomDate(),
			Valid: true,
		},
		RoomID: util.RandomID(),
		ReservationID: pgtype.Int8{
			Int64: util.RandomID(),
			Valid: true,
		},
		Restriction: RestrictionOwnerBlock,
	}

	data := map[string]string{
		"start_date":     arg.StartDate.Time.Format(config.DateLayout),
		"end_date":       arg.EndDate.Time.Format(config.DateLayout),
		"room_id":        fmt.Sprint(arg.RoomID),
		"reservation_id": fmt.Sprint(arg.ReservationID.Int64),
		"restriction":    string(arg.Restriction),
	}

	arg2 := CreateRoomRestrictionParams{}
	err := arg2.Unmarshal(data)
	require.NoError(t, err)
	assert.Equal(t, arg, arg2)
}

func TestUpdateRoomRestrictionParams_Unmarshal(t *testing.T) {
	arg := UpdateRoomRestrictionParams{
		ID: util.RandomID(),
		StartDate: pgtype.Date{
			Time:  util.RandomDate(),
			Valid: true,
		},
		EndDate: pgtype.Date{
			Time:  util.RandomDate(),
			Valid: true,
		},
		RoomID: util.RandomID(),
		ReservationID: pgtype.Int8{
			Int64: util.RandomID(),
			Valid: true,
		},
		Restriction: RestrictionOwnerBlock,
		UpdatedAt: pgtype.Timestamptz{
			Time:  util.RandomDatetime(),
			Valid: true,
		},
	}

	data := map[string]string{
		"id":             fmt.Sprint(arg.ID),
		"start_date":     arg.StartDate.Time.Format(config.DateLayout),
		"end_date":       arg.EndDate.Time.Format(config.DateLayout),
		"room_id":        fmt.Sprint(arg.RoomID),
		"reservation_id": fmt.Sprint(arg.ReservationID.Int64),
		"restriction":    string(arg.Restriction),
		"updated_at":     arg.UpdatedAt.Time.Format(config.DateTimeLayout),
	}

	arg2 := UpdateRoomRestrictionParams{}
	err := arg2.Unmarshal(data)
	require.NoError(t, err)

	assert.Equal(t, arg.ID, arg2.ID)
	assert.Equal(t, arg.StartDate, arg2.StartDate)
	assert.Equal(t, arg.EndDate, arg2.EndDate)
	assert.Equal(t, arg.RoomID, arg2.RoomID)
	assert.Equal(t, arg.ReservationID, arg2.ReservationID)
	assert.Equal(t, arg.Restriction, arg2.Restriction)
	assert.WithinDuration(t, arg.UpdatedAt.Time, arg2.UpdatedAt.Time, time.Second)
}
