package db

import (
	"fmt"
	"testing"

	"github.com/github-real-lb/bookings-web-app/util"
	"github.com/github-real-lb/bookings-web-app/util/config"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func TestCheckRoomAvailabilityParams_Unmarshal(t *testing.T) {
	arg := CheckRoomAvailabilityParams{
		RoomID: util.RandomID(),
		StartDate: pgtype.Date{
			Time:  util.RandomDate(),
			Valid: true,
		},
		EndDate: pgtype.Date{
			Time:  util.RandomDate(),
			Valid: true,
		},
	}

	data := map[string]string{
		"room_id":    fmt.Sprint(arg.RoomID),
		"start_date": arg.StartDate.Time.Format(config.DateLayout),
		"end_date":   arg.EndDate.Time.Format(config.DateLayout),
	}

	arg2 := CheckRoomAvailabilityParams{}
	err := arg2.Unmarshal(data)
	require.NoError(t, err)
	require.Equal(t, arg.RoomID, arg2.RoomID)
	require.Equal(t, arg.StartDate, arg2.StartDate)
	require.Equal(t, arg.EndDate, arg2.EndDate)
}

func TestCreateRoomParams_Unmarshal(t *testing.T) {
	arg := CreateRoomParams{
		Name:          util.RandomName(),
		Description:   util.RandomNote(),
		ImageFilename: fmt.Sprint(util.RandomName(), ".png"),
	}

	data := map[string]string{
		"name":           arg.Name,
		"description":    arg.Description,
		"image_filename": arg.ImageFilename,
	}

	arg2 := CreateRoomParams{}
	arg2.Unmarshal(data)

	require.Equal(t, arg.Name, arg2.Name)
	require.Equal(t, arg.Description, arg2.Description)
	require.Equal(t, arg.ImageFilename, arg2.ImageFilename)

}

func TestListAvailableRoomsParams_Unmarshal(t *testing.T) {
	arg := ListAvailableRoomsParams{
		Limit:  int32(util.RandomInt64(0, 10000)),
		Offset: int32(util.RandomInt64(0, 10000)),
		StartDate: pgtype.Date{
			Time:  util.RandomDate(),
			Valid: true,
		},
		EndDate: pgtype.Date{
			Time:  util.RandomDate(),
			Valid: true,
		},
	}

	data := map[string]string{
		"limit":      fmt.Sprint(arg.Limit),
		"offset":     fmt.Sprint(arg.Offset),
		"start_date": arg.StartDate.Time.Format(config.DateLayout),
		"end_date":   arg.EndDate.Time.Format(config.DateLayout),
	}

	arg2 := ListAvailableRoomsParams{}
	err := arg2.Unmarshal(data)
	require.NoError(t, err)

	require.Equal(t, arg.Limit, arg2.Limit)
	require.Equal(t, arg.Offset, arg2.Offset)
	require.Equal(t, arg.StartDate, arg2.StartDate)
	require.Equal(t, arg.EndDate, arg2.EndDate)
}
