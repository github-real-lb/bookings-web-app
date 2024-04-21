package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/github-real-lb/bookings-web-app/util"
	"github.com/github-real-lb/bookings-web-app/util/config"
	"github.com/stretchr/testify/require"
)

func TestReservation_GenerateReservationCode(t *testing.T) {
	// N states the number of times to test randomness
	const N int = 10

	codes := make(util.KeysMap)
	for i := 0; i < N; i++ {
		date := util.RandomDate()

		r := Reservation{
			LastName:  util.RandomName(),
			StartDate: date,
			EndDate:   date.Add(time.Hour * 24 * 7),
		}

		r.GenerateReservationCode()
		util.RequireUnique(t, r.Code, codes)
	}
}

func TestReservation_MarshalAndUnmarshal(t *testing.T) {
	r := Reservation{
		ID:        util.RandomID(),
		Code:      util.RandomString(ReservationCodeLenght),
		FirstName: util.RandomName(),
		LastName:  util.RandomName(),
		Email:     util.RandomEmail(),
		Phone:     util.RandomPhone(),
		StartDate: util.RandomDate(),
		EndDate:   util.RandomDate(),
		Room:      randomRoom(),
		Notes:     util.RandomNote(),
		CreatedAt: util.RandomDatetime(),
		UpdatedAt: util.RandomDatetime(),
	}

	data := r.Marshal()
	require.Equal(t, fmt.Sprint(r.ID), data["id"])

	require.Equal(t, r.Code, data["code"])
	require.Equal(t, r.FirstName, data["first_name"])
	require.Equal(t, r.LastName, data["last_name"])
	require.Equal(t, r.Email, data["email"])
	require.Equal(t, r.Phone, data["phone"])

	require.Equal(t, r.StartDate.Format(config.DateLayout), data["start_date"])
	require.Equal(t, r.EndDate.Format(config.DateLayout), data["end_date"])

	require.Equal(t, fmt.Sprint(r.Room.ID), data["room_id"])
	require.Equal(t, r.Room.Name, data["room_name"])
	require.Equal(t, r.Room.Description, data["room_description"])
	require.Equal(t, r.Room.ImageFilename, data["room_image_filename"])

	require.Equal(t, r.Notes, data["notes"])

	require.Equal(t, r.CreatedAt.Format(config.DateTimeLayout), data["created_at"])
	require.Equal(t, r.UpdatedAt.Format(config.DateTimeLayout), data["updated_at"])

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

	require.Equal(t, r.Room.ID, r2.Room.ID)
	require.Equal(t, r.Room.Name, r2.Room.Name)
	require.Equal(t, r.Room.Description, r2.Room.Description)
	require.Equal(t, r.Room.ImageFilename, r2.Room.ImageFilename)

	require.Equal(t, r.Notes, r2.Notes)

	require.WithinDuration(t, r.CreatedAt, r2.CreatedAt, time.Second)
	require.WithinDuration(t, r.UpdatedAt, r2.UpdatedAt, time.Second)
}

// func TestRoom_MarshalAndUnmarshal(t *testing.T) {
// 	r := Room{
// 		ID:            util.RandomID(),
// 		Name:          util.RandomName(),
// 		Description:   util.RandomNote(),
// 		ImageFilename: fmt.Sprint(util.RandomName(), ".png"),
// 		CreatedAt:     util.RandomDatetime(),
// 		UpdatedAt:     util.RandomDatetime(),
// 	}

// 	data := r.Marshal()
// 	require.Len(t, data, 6)
// 	require.Equal(t, fmt.Sprint(r.ID), data["id"])
// 	require.Equal(t, r.Name, data["name"])
// 	require.Equal(t, r.Description, data["description"])
// 	require.Equal(t, r.ImageFilename, data["image_filename"])

// 	require.Equal(t, r.CreatedAt.Format(config.DateTimeLayout), data["created_at"])
// 	require.Equal(t, r.UpdatedAt.Format(config.DateTimeLayout), data["updated_at"])

// 	r2 := Room{}
// 	err := r2.Unmarshal(data)
// 	require.NoError(t, err)
// 	require.Equal(t, r.ID, r2.ID)
// 	require.Equal(t, r.Name, r2.Name)
// 	require.Equal(t, r.Description, r2.Description)
// 	require.Equal(t, r.ImageFilename, r2.ImageFilename)
// 	require.WithinDuration(t, r.CreatedAt, r2.CreatedAt, time.Second)
// 	require.WithinDuration(t, r.UpdatedAt, r2.UpdatedAt, time.Second)
// }

func TestRestriction_Scan(t *testing.T) {
	var r Restriction = RestrictionReservation

	err := r.Scan("owner_block")
	require.NoError(t, err)
	require.Equal(t, RestrictionOwnerBlock, r)

	err = r.Scan([]byte("reservation"))
	require.NoError(t, err)
	require.Equal(t, RestrictionReservation, r)

	err = r.Scan(Reservation{})
	require.ErrorContains(t, err, "unsupported scan type for Restriction:")
}
