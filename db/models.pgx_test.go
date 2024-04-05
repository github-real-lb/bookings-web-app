package db

import (
	"fmt"
	"testing"
	"time"

	"github.com/github-real-lb/bookings-web-app/util"
	"github.com/github-real-lb/bookings-web-app/util/config"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func TestReservation_MarshalAndUnmarhsal(t *testing.T) {
	r := Reservation{
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
		CreatedAt: pgtype.Timestamptz{
			Time:  util.RandomDatetime(),
			Valid: true,
		},
		UpdatedAt: pgtype.Timestamptz{
			Time:  util.RandomDatetime(),
			Valid: true,
		},
	}

	data := r.Marshal()
	require.Equal(t, fmt.Sprint(r.ID), data["id"])

	require.Equal(t, r.Code, data["code"])
	require.Equal(t, r.FirstName, data["first_name"])
	require.Equal(t, r.LastName, data["last_name"])
	require.Equal(t, r.Email, data["email"])

	require.Equal(t, r.Phone.String, data["phone"])
	require.True(t, r.Phone.Valid)

	require.Equal(t, r.StartDate.Time.Format(config.DateLayout), data["start_date"])
	require.True(t, r.StartDate.Valid)
	require.Equal(t, r.EndDate.Time.Format(config.DateLayout), data["end_date"])
	require.True(t, r.StartDate.Valid)

	require.Equal(t, fmt.Sprint(r.RoomID), data["room_id"])

	require.Equal(t, r.Notes.String, data["notes"])
	require.True(t, r.Notes.Valid)

	require.Equal(t, r.CreatedAt.Time.Format(config.DateTimeLayout), data["created_at"])
	require.True(t, r.CreatedAt.Valid)

	require.Equal(t, r.UpdatedAt.Time.Format(config.DateTimeLayout), data["updated_at"])
	require.True(t, r.UpdatedAt.Valid)

	r2 := Reservation{}
	err := r2.Unmarshal(data)
	require.NoError(t, err)

	require.Equal(t, r.ID, r2.ID)
	require.Equal(t, r.Code, r2.Code)
	require.Equal(t, r.FirstName, r2.FirstName)
	require.Equal(t, r.LastName, r2.LastName)
	require.Equal(t, r.Email, r2.Email)
	require.Equal(t, r.Phone, r2.Phone)
	require.Equal(t, r.StartDate, r2.StartDate)
	require.Equal(t, r.EndDate, r2.EndDate)
	require.Equal(t, r.RoomID, r2.RoomID)
	require.Equal(t, r.Notes, r2.Notes)

	require.WithinDuration(t, r.CreatedAt.Time, r2.CreatedAt.Time, time.Second)
	require.True(t, r2.CreatedAt.Valid)

	require.WithinDuration(t, r.UpdatedAt.Time, r2.UpdatedAt.Time, time.Second)
	require.True(t, r2.UpdatedAt.Valid)
}

func TestRoom_Marshal(t *testing.T) {
	r := Room{
		ID:            util.RandomID(),
		Name:          util.RandomName(),
		Description:   util.RandomNote(),
		ImageFilename: fmt.Sprint(util.RandomName(), ".png"),
		CreatedAt: pgtype.Timestamptz{
			Time:  util.RandomDatetime(),
			Valid: true,
		},
		UpdatedAt: pgtype.Timestamptz{
			Time:  util.RandomDatetime(),
			Valid: true,
		},
	}

	data := r.Marshal()
	require.Len(t, data, 6)
	require.Equal(t, fmt.Sprint(r.ID), data["id"])
	require.Equal(t, r.Name, data["name"])
	require.Equal(t, r.Description, data["description"])
	require.Equal(t, r.ImageFilename, data["image_filename"])

	require.Equal(t, r.CreatedAt.Time.Format(config.DateTimeLayout), data["created_at"])
	require.True(t, r.CreatedAt.Valid)

	require.Equal(t, r.UpdatedAt.Time.Format(config.DateTimeLayout), data["updated_at"])
	require.True(t, r.UpdatedAt.Valid)
}

func TestRoomRestriction_Marshal(t *testing.T) {
	r := RoomRestriction{
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
		CreatedAt: pgtype.Timestamptz{
			Time:  util.RandomDatetime(),
			Valid: true,
		},
		UpdatedAt: pgtype.Timestamptz{
			Time:  util.RandomDatetime(),
			Valid: true,
		},
	}

	data := r.Marshal()
	require.Len(t, data, 8)
	require.Equal(t, fmt.Sprint(r.ID), data["id"])

	require.Equal(t, r.StartDate.Time.Format(config.DateLayout), data["start_date"])
	require.Equal(t, r.EndDate.Time.Format(config.DateLayout), data["end_date"])

	require.Equal(t, fmt.Sprint(r.RoomID), data["room_id"])

	require.Equal(t, fmt.Sprint(r.ReservationID.Int64), data["reservation_id"])
	require.True(t, r.ReservationID.Valid)

	require.Equal(t, r.Restriction, RestrictionOwnerBlock)

	require.Equal(t, r.CreatedAt.Time.Format(config.DateTimeLayout), data["created_at"])
	require.True(t, r.CreatedAt.Valid)

	require.Equal(t, r.UpdatedAt.Time.Format(config.DateTimeLayout), data["updated_at"])
	require.True(t, r.UpdatedAt.Valid)
}
