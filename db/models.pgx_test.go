package db

import (
	"fmt"
	"testing"

	"github.com/github-real-lb/bookings-web-app/util"
	"github.com/github-real-lb/bookings-web-app/util/config"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func TestReservation_MarshalAndUnmarhsal(t *testing.T) {
	data := make(map[string]string)
	data["id"] = fmt.Sprint(util.RandomID())
	data["code"] = util.RandomString(ReservationCodeLenght)
	data["first_name"] = util.RandomName()
	data["last_name"] = util.RandomName()
	data["email"] = util.RandomEmail()
	data["phone"] = util.RandomPhone()
	data["start_date"] = util.RandomDate().Format(config.DateLayout)
	data["end_date"] = util.RandomDate().Format(config.DateLayout)
	data["room_id"] = fmt.Sprint(util.RandomID())
	data["notes"] = util.RandomNote()
	data["created_at"] = util.RandomDatetime().Format(config.DateTimeLayout)
	data["updated_at"] = util.RandomDatetime().Format(config.DateTimeLayout)

	r := Reservation{}
	err := r.Unmarshal(data)
	require.NoError(t, err)
	require.Equal(t, data["id"], fmt.Sprint(r.ID))

	require.Equal(t, data["code"], r.Code)
	require.Equal(t, data["first_name"], r.FirstName)
	require.Equal(t, data["last_name"], r.LastName)
	require.Equal(t, data["email"], r.Email)

	require.Equal(t, data["phone"], r.Phone.String)
	require.True(t, r.Phone.Valid)

	require.Equal(t, data["start_date"], r.StartDate.Time.Format(config.DateLayout))
	require.True(t, r.StartDate.Valid)
	require.Equal(t, data["end_date"], r.EndDate.Time.Format(config.DateLayout))
	require.True(t, r.StartDate.Valid)

	require.Equal(t, data["room_id"], fmt.Sprint(r.RoomID))

	require.Equal(t, data["notes"], r.Notes.String)
	require.True(t, r.Notes.Valid)

	require.Equal(t, data["created_at"], r.CreatedAt.Time.Format(config.DateTimeLayout))
	require.True(t, r.CreatedAt.Valid)

	require.Equal(t, data["updated_at"], r.UpdatedAt.Time.Format(config.DateTimeLayout))
	require.True(t, r.UpdatedAt.Valid)

	marshaledData := r.Marshal()
	require.Equal(t, data, marshaledData)
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
