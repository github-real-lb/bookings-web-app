package main

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/github-real-lb/bookings-web-app/db"
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
	tests := []struct {
		Name      string
		Available bool
		Error     error
	}{
		{
			Name:      "Room Available",
			Available: true,
			Error:     nil,
		},
		{
			Name:      "Room Unavailable",
			Available: false,
			Error:     nil,
		},
		{
			Name:      "Error",
			Available: false,
			Error:     errors.New("any error"),
		},
	}

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

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			// create a new server with mock database store
			server, mockStore := NewTestServer(t)

			// build stub
			mockStore.On("CheckRoomAvailability", mock.Anything, arg).
				Return(test.Available, test.Error).
				Once()

			// execute method
			ok, err := server.CheckRoomAvailability(reservation)

			// testify
			assert.Equal(t, test.Available, ok)
			assert.Equal(t, test.Error, err)
		})
	}
}

func TestServer_CreateReservation(t *testing.T) {
	t.Run("Test OK", func(t *testing.T) {
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
		server, mockStore := NewTestServer(t)

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
	})

	t.Run("Test Error", func(t *testing.T) {
		// create random reservation with room data
		reservation := randomReservation()
		reservationData := reservation.Marshal()

		arg := db.CreateReservationParams{}
		err := arg.Unmarshal(reservationData)
		require.NoError(t, err)

		// create a new server with mock database store
		server, mockStore := NewTestServer(t)

		// build stub
		mockStore.On("CreateReservationTx", mock.Anything, arg).
			Return(db.Reservation{}, errors.New("any error")).
			Once()

		// create a copy of reservation to pass to method
		copyReservation := reservation

		// execute method
		err = server.CreateReservation(&copyReservation)

		// tesify
		assert.Error(t, err)
		assert.Equal(t, reservation, copyReservation)
	})
}

func TestServer_ListAvailableRooms(t *testing.T) {
	t.Run("Test Available Rooms", func(t *testing.T) {
		// create random reservation with room data
		reservation := randomReservation()

		//create rooms slice with random data of n rooms
		const N = 5
		rooms := randomRooms(N)

		//create mockStore method return arguments
		arg := db.ListAvailableRoomsParams{
			Limit:  N,
			Offset: 0,
			StartDate: pgtype.Date{
				Time:  reservation.StartDate,
				Valid: true,
			},
			EndDate: pgtype.Date{
				Time:  reservation.EndDate,
				Valid: true,
			},
		}

		//create mockStore method return arguments
		dbRooms := make([]db.Room, N)
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
		server, mockStore := NewTestServer(t)

		// build stub
		mockStore.On("ListAvailableRooms", mock.Anything, arg).
			Return(dbRooms, nil).
			Once()

		// execute method
		resultRooms, err := server.ListAvailableRooms(reservation, N, 0)

		// tesify
		assert.NoError(t, err)
		require.Len(t, resultRooms, N)

		for i := 0; i < N; i++ {
			room := rooms[i]
			resultRoom := resultRooms[i]
			require.Equal(t, room.ID, resultRoom.ID)
			require.Equal(t, room.Name, resultRoom.Name)
			require.Equal(t, room.Description, resultRoom.Description)
			require.Equal(t, room.ImageFilename, resultRoom.ImageFilename)
			require.WithinDuration(t, room.CreatedAt, resultRoom.CreatedAt, time.Second)
			require.WithinDuration(t, room.UpdatedAt, resultRoom.UpdatedAt, time.Second)
		}
	})

	t.Run("Test No Available Rooms", func(t *testing.T) {
		// create random reservation with room data
		reservation := randomReservation()

		//create mockStore method return arguments
		arg := db.ListAvailableRoomsParams{
			Limit:  5,
			Offset: 0,
			StartDate: pgtype.Date{
				Time:  reservation.StartDate,
				Valid: true,
			},
			EndDate: pgtype.Date{
				Time:  reservation.EndDate,
				Valid: true,
			},
		}

		// create a new server with mock database store
		server, mockStore := NewTestServer(t)

		// build stub
		mockStore.On("ListAvailableRooms", mock.Anything, arg).
			Return([]db.Room{}, nil).
			Once()

		// execute method
		resultRooms, err := server.ListAvailableRooms(reservation, 5, 0)

		// tesify
		assert.NoError(t, err)
		require.Len(t, resultRooms, 0)
	})

	t.Run("Test Error", func(t *testing.T) {
		// create random reservation with room data
		reservation := randomReservation()

		//create mockStore method return arguments
		arg := db.ListAvailableRoomsParams{
			Limit:  5,
			Offset: 0,
			StartDate: pgtype.Date{
				Time:  reservation.StartDate,
				Valid: true,
			},
			EndDate: pgtype.Date{
				Time:  reservation.EndDate,
				Valid: true,
			},
		}

		// create a new server with mock database store
		server, mockStore := NewTestServer(t)

		// build stub
		mockStore.On("ListAvailableRooms", mock.Anything, arg).
			Return(nil, errors.New("any error")).
			Once()

		// execute method
		rooms, err := server.ListAvailableRooms(reservation, 5, 0)

		// tesify
		assert.Error(t, err)
		require.Nil(t, rooms)
	})
}

func TestServer_ListRooms(t *testing.T) {
	t.Run("Test All Rooms", func(t *testing.T) {
		//create rooms slice with random data of n rooms
		const N = 5
		rooms := randomRooms(N)

		//create mockStore method return arguments
		arg := db.ListRoomsParams{
			Limit:  int32(N),
			Offset: 0,
		}

		//create mockStore method return arguments
		dbRooms := make([]db.Room, N)
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
		server, mockStore := NewTestServer(t)

		// build stub
		mockStore.On("ListRooms", mock.Anything, arg).
			Return(dbRooms, nil).
			Once()

		// execute method
		returnedRooms, err := server.ListRooms(N, 0)

		// tesify
		assert.NoError(t, err)
		require.Len(t, returnedRooms, N)

		for i := 0; i < N; i++ {
			room := rooms[i]
			returnedRoom := returnedRooms[i]
			require.Equal(t, room.ID, returnedRoom.ID)
			require.Equal(t, room.Name, returnedRoom.Name)
			require.Equal(t, room.Description, returnedRoom.Description)
			require.Equal(t, room.ImageFilename, returnedRoom.ImageFilename)
			require.WithinDuration(t, room.CreatedAt, returnedRoom.CreatedAt, time.Second)
			require.WithinDuration(t, room.UpdatedAt, returnedRoom.UpdatedAt, time.Second)
		}
	})

	t.Run("Test No Rooms", func(t *testing.T) {
		//create mockStore method return arguments
		arg := db.ListRoomsParams{
			Limit:  5,
			Offset: 0,
		}

		// create a new server with mock database store
		server, mockStore := NewTestServer(t)

		// build stub
		mockStore.On("ListRooms", mock.Anything, arg).
			Return([]db.Room{}, nil).
			Once()

		// execute method
		ResultRooms, err := server.ListRooms(5, 0)

		// tesify
		assert.NoError(t, err)
		require.Len(t, ResultRooms, 0)
	})

	t.Run("Test Error", func(t *testing.T) {
		//create mockStore method return arguments
		arg := db.ListRoomsParams{
			Limit:  5,
			Offset: 0,
		}

		// create a new server with mock database store
		server, mockStore := NewTestServer(t)

		// build stub
		mockStore.On("ListRooms", mock.Anything, arg).
			Return(nil, errors.New("any error")).
			Once()

		// execute method
		rooms, err := server.ListRooms(int(arg.Limit), 0)

		// tesify
		assert.Error(t, err)
		require.Nil(t, rooms)
	})
}
