package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/github-real-lb/bookings-web-app/db"
	"github.com/github-real-lb/bookings-web-app/db/mocks"
	"github.com/github-real-lb/bookings-web-app/util"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// randomRoom returns a Room struct with random data
func randomRoom() Room {
	randomTime := util.RandomDatetime()
	return Room{
		ID:            util.RandomInt64(1, 100),
		Name:          util.RandomName(),
		Description:   util.RandomNote(),
		ImageFilename: fmt.Sprint(util.RandomName(), ".png"),
		CreatedAt:     randomTime,
		UpdatedAt:     randomTime,
	}
}

// randomRooms returns a Rooms slice with n random rooms data
func randomRooms(n int) Rooms {
	rooms := make([]Room, n)

	for i := 0; i < n; i++ {
		rooms[i] = randomRoom()
	}

	return rooms
}

// randomReservation returns a Reservation struct with random data
func randomReservation() Reservation {
	randomDate := util.RandomDate()

	return Reservation{
		ID:        util.RandomInt64(1, 1000),
		Code:      util.RandomString(ReservationCodeLenght),
		FirstName: util.RandomName(),
		LastName:  util.RandomName(),
		Email:     util.RandomEmail(),
		Phone:     util.RandomPhoneNumber(),
		StartDate: randomDate.Add(time.Hour * 24 * 30),
		EndDate:   randomDate.Add(time.Hour * 24 * 37),
		Room:      randomRoom(),
		Notes:     util.RandomNote(),
		CreatedAt: randomDate.Add(time.Hour * 3),
		UpdatedAt: randomDate.Add(time.Hour * 3),
	}
}

func TestServer_CheckRoomAvailability(t *testing.T) {
	// create random reservation with room data
	reservation := randomReservation()

	// create mockStore mehod arguments
	arg := db.CheckRoomAvailabilityParams{
		RoomID: reservation.Room.ID,
		EndDate: pgtype.Date{
			Time:  reservation.EndDate,
			Valid: true,
		},
		StartDate: pgtype.Date{
			Time:  reservation.StartDate,
			Valid: true,
		},
	}

	// define the func to receive reservation and return true + nil
	mockStore := mocks.NewMockStore(t)
	mockStore.On("CheckRoomAvailability", mock.Anything, arg).
		Return(true, nil).
		Once()

	// test the func
	server := NewServer(mockStore)
	ok, err := server.CheckRoomAvailability(reservation)
	assert.NoError(t, err)
	assert.True(t, ok)
}

func TestServer_CreateReservation(t *testing.T) {
	// create random reservation with room data
	reservation := randomReservation()

	// create mockStore mehod arguments
	restrictionID := util.RandomInt64(1, 100)
	arg := db.CreateReservationParams{
		Code:      reservation.Code,
		FirstName: reservation.FirstName,
		LastName:  reservation.LastName,
		Email:     reservation.Email,
		Phone: pgtype.Text{
			String: reservation.Phone,
			Valid:  true,
		},
		StartDate: pgtype.Date{
			Time:  reservation.StartDate,
			Valid: true,
		},
		EndDate: pgtype.Date{
			Time:  reservation.EndDate,
			Valid: true,
		},
		RoomID: reservation.Room.ID,
		Notes: pgtype.Text{
			String: reservation.Notes,
			Valid:  true,
		},
	}

	//create mockStore method return arguments
	dbReservation := db.Reservation{
		ID:        reservation.ID,
		Code:      reservation.Code,
		FirstName: reservation.FirstName,
		LastName:  reservation.LastName,
		Email:     reservation.Email,
		Phone: pgtype.Text{
			String: reservation.Phone,
			Valid:  true,
		},
		StartDate: pgtype.Date{
			Time:  reservation.StartDate,
			Valid: true,
		},
		EndDate: pgtype.Date{
			Time:  reservation.EndDate,
			Valid: true,
		},
		RoomID: reservation.Room.ID,
		Notes: pgtype.Text{
			String: reservation.Notes,
			Valid:  true,
		},
		CreatedAt: pgtype.Timestamptz{
			Time:  reservation.CreatedAt,
			Valid: true,
		},
		UpdatedAt: pgtype.Timestamptz{
			Time:  reservation.UpdatedAt,
			Valid: true,
		},
	}

	// define the func to receive reservation and return nil
	mockStore := mocks.NewMockStore(t)
	mockStore.On("CreateReservationTx", mock.Anything, arg, restrictionID).
		Return(dbReservation, nil).
		Once()

	// create a copy of reservation to pass to method
	copyReservation := reservation

	// test the func
	server := NewServer(mockStore)
	err := server.CreateReservation(&copyReservation, restrictionID)
	assert.NoError(t, err)
	assert.Equal(t, reservation, copyReservation)
}

func TestServer_ListAvailableRooms(t *testing.T) {
	// create random reservation with room data
	reservation := randomReservation()

	//create rooms slice with random data of n rooms
	n := 5
	rooms := randomRooms(n)

	//create mockStore method return arguments
	arg := db.ListAvailableRoomsParams{
		EndDate: pgtype.Date{
			Time:  reservation.EndDate,
			Valid: true,
		},
		StartDate: pgtype.Date{
			Time:  reservation.StartDate,
			Valid: true,
		},
	}

	//create mockStore method return arguments
	dbRooms := make([]db.Room, n)
	for i, v := range rooms {
		dbRooms[i] = db.Room{
			ID:          v.ID,
			Name:        v.Name,
			Description: v.Description,
			ImageFilename: pgtype.Text{
				String: v.ImageFilename,
				Valid:  true,
			},
			CreatedAt: pgtype.Timestamptz{
				Time:  v.CreatedAt,
				Valid: true,
			},
			UpdatedAt: pgtype.Timestamptz{
				Time:  v.UpdatedAt,
				Valid: true,
			},
		}
	}

	// define the func to receive reservation and return nil
	mockStore := mocks.NewMockStore(t)
	mockStore.On("ListAvailableRooms", mock.Anything, arg).
		Return(dbRooms, nil).
		Once()

	// test the func
	server := NewServer(mockStore)
	returnedRooms, err := server.ListAvailableRooms(reservation)
	assert.NoError(t, err)
	require.Len(t, returnedRooms, n)

	for i := 0; i < n; i++ {
		room := rooms[i]
		returnedRoom := returnedRooms[i]
		require.Equal(t, room.ID, returnedRoom.ID)
		require.Equal(t, room.Name, returnedRoom.Name)
		require.Equal(t, room.Description, returnedRoom.Description)
		require.Equal(t, room.ImageFilename, returnedRoom.ImageFilename)
		require.WithinDuration(t, room.CreatedAt, returnedRoom.CreatedAt, time.Second)
		require.WithinDuration(t, room.UpdatedAt, returnedRoom.UpdatedAt, time.Second)
	}
}
