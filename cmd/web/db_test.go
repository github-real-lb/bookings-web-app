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
		ID:            util.RandomID(),
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
		ID:        util.RandomID(),
		Code:      util.RandomString(ReservationCodeLenght),
		FirstName: util.RandomName(),
		LastName:  util.RandomName(),
		Email:     util.RandomEmail(),
		Phone:     util.RandomPhone(),
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

	// create a new server with mock database store
	mockStore := mocks.NewMockStore(t)
	server := NewServer(mockStore)

	// build stub
	mockStore.On("CheckRoomAvailability", mock.Anything, arg).
		Return(true, nil).
		Once()

	// execute method
	ok, err := server.CheckRoomAvailability(reservation)

	// testify
	assert.NoError(t, err)
	assert.True(t, ok)
}

func TestServer_CreateReservation(t *testing.T) {
	// create random reservation with room data
	reservation := randomReservation()
	reservationData := reservation.Marshal()

	arg := db.CreateReservationParams{}
	err := arg.Unmarshal(reservationData)
	require.NoError(t, err)

	//create method return arguments
	dbReservation := db.Reservation{}
	err = dbReservation.Unmarshal(reservationData)
	require.NoError(t, err)

	// create a new server with mock database store
	mockStore := mocks.NewMockStore(t)
	server := NewServer(mockStore)

	// build stub
	mockStore.On("CreateReservationTx", mock.Anything, arg).
		Return(dbReservation, nil).
		Once()

	// create a copy of reservation to pass to method
	copyReservation := reservation

	// execute method
	err = server.CreateReservation(&copyReservation)

	// tesify
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
			ID:            v.ID,
			Name:          v.Name,
			Description:   v.Description,
			ImageFilename: v.ImageFilename,
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

	// create a new server with mock database store
	mockStore := mocks.NewMockStore(t)
	server := NewServer(mockStore)

	// build stub
	mockStore.On("ListAvailableRooms", mock.Anything, arg).
		Return(dbRooms, nil).
		Once()

	// execute method
	returnedRooms, err := server.ListAvailableRooms(reservation)

	// tesify
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
