package main

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/github-real-lb/bookings-web-app/db"
	"github.com/github-real-lb/bookings-web-app/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// randomReservation returns a Reservation struct with random data
func randomReservation() Reservation {
	rDate := util.RandomDate()
	rRoom := randomRoom()

	return Reservation{
		ID:        util.RandomID(),
		Code:      util.RandomString(ReservationCodeLenght),
		FirstName: util.RandomName(),
		LastName:  util.RandomName(),
		Email:     util.RandomEmail(),
		Phone:     util.RandomPhone(),
		StartDate: rDate.Add(time.Hour * 24 * 30),
		EndDate:   rDate.Add(time.Hour * 24 * 37),
		RoomID:    rRoom.ID,
		Notes:     util.RandomNote(),
		CreatedAt: rDate.Add(time.Hour * 3),
		UpdatedAt: rDate.Add(time.Hour * 3),
		Room:      rRoom,
	}
}

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

// randomUser returns a User struct with random data
func randomUser() User {
	randomTime := util.RandomDatetime()

	return User{
		ID:          util.RandomID(),
		FirstName:   util.RandomName(),
		LastName:    util.RandomName(),
		Email:       util.RandomEmail(),
		Password:    util.RandomPassword(),
		AccessLevel: util.RandomInt64(1, 10),
		CreatedAt:   randomTime,
		UpdatedAt:   randomTime,
	}
}

func TestServer_AuthenticateUser(t *testing.T) {
	// create new test server
	ts := NewTestServer(t)

	// create random user with room data
	user := randomUser()

	// create stub call arguments
	arg := db.AuthenticateUserParams{
		Email:    user.Email,
		Password: user.Password,
	}

	t.Run("OK", func(t *testing.T) {
		// create stub return arguments
		dbUser := db.User{}
		err := util.CopyDataUsingJSON(user, &dbUser)
		require.NoError(t, err)

		// build stub
		ts.MockDBStore.On("AuthenticateUser", mock.Anything, arg).
			Return(dbUser, nil).
			Once()

		result, err := ts.AuthenticateUser(user.Email, user.Password)
		require.NoError(t, err)
		assert.Equal(t, user.ID, result)
	})

	t.Run("Error", func(t *testing.T) {
		// build stub
		ts.MockDBStore.On("AuthenticateUser", mock.Anything, arg).
			Return(db.User{}, errors.New("any error")).
			Once()

		result, err := ts.AuthenticateUser(user.Email, user.Password)
		require.Error(t, err)
		require.Equal(t, int64(0), result)

	})

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
	rsv := randomReservation()

	// create ts.MockDBStore mehod arguments
	arg := db.CheckRoomAvailabilityParams{}
	err := CopyStructDataToDBStruct(rsv, &arg)
	require.NoError(t, err)

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			// create a new server with mock database store
			ts := NewTestServer(t)

			// build stub
			ts.MockDBStore.On("CheckRoomAvailability", mock.Anything, arg).
				Return(test.Available, test.Error).
				Once()

			// execute method
			ok, err := ts.CheckRoomAvailability(rsv.RoomID, rsv.StartDate, rsv.EndDate)

			// testify
			assert.Equal(t, test.Available, ok)
			assert.Equal(t, test.Error, err)
		})
	}
}

func TestServer_CreateReservation(t *testing.T) {
	t.Run("Test OK", func(t *testing.T) {
		// create random reservation with room data
		rsv := randomReservation()

		// create stub call arguments
		arg := db.CreateReservationParams{}
		err := CopyStructDataToDBStruct(rsv, &arg)
		require.NoError(t, err)

		// create stub return arguments
		dbRsv := db.Reservation{}
		err = CopyStructDataToDBStruct(rsv, &dbRsv)
		require.NoError(t, err)

		// create a new server with mock database store
		ts := NewTestServer(t)

		// build stub
		ts.MockDBStore.On("CreateReservationTx", mock.Anything, arg).
			Return(dbRsv, nil).
			Once()

		// execute method
		err = ts.CreateReservation(rsv)

		// tesify
		assert.NoError(t, err)
	})

	t.Run("Test Error", func(t *testing.T) {
		// create random reservation with room data
		rsv := randomReservation()

		// create stub call arguments
		arg := db.CreateReservationParams{}
		err := CopyStructDataToDBStruct(rsv, &arg)
		require.NoError(t, err)

		// create a new server with mock database store
		ts := NewTestServer(t)

		// build stub
		ts.MockDBStore.On("CreateReservationTx", mock.Anything, arg).
			Return(db.Reservation{}, errors.New("any error")).
			Once()

		// execute method
		err = ts.CreateReservation(rsv)

		// tesify
		assert.Error(t, err)
	})
}

func TestServer_ListAvailableRooms(t *testing.T) {
	// create random reservation with room data
	rsv := randomReservation()

	//create db stub call arguments
	arg := db.ListAvailableRoomsParams{}
	err := CopyStructDataToDBStruct(rsv, &arg)
	require.NoError(t, err)

	arg.Limit = LimitRoomsPerPage
	arg.Offset = 0

	t.Run("Test Available Rooms", func(t *testing.T) {
		// create stub return arguments
		rooms := make(Rooms, LimitRoomsPerPage)
		dbRooms := make([]db.Room, LimitRoomsPerPage)

		for i := 0; i < LimitRoomsPerPage; i++ {
			rooms[i] = randomRoom()
			err = util.CopyDataUsingJSON(rooms[i], &dbRooms[i])
			require.NoError(t, err)
		}

		// create a new server with mock database store
		ts := NewTestServer(t)

		// build stub
		ts.MockDBStore.On("ListAvailableRooms", mock.Anything, arg).
			Return(dbRooms, nil).
			Once()

		// execute method
		resultRooms, err := ts.ListAvailableRooms(int(arg.Limit), int(arg.Offset), rsv.StartDate, rsv.EndDate)

		// tesify
		assert.NoError(t, err)
		require.Len(t, resultRooms, LimitRoomsPerPage)

		for i := 0; i < LimitRoomsPerPage; i++ {
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
		// create a new server with mock database store
		ts := NewTestServer(t)

		// build stub
		ts.MockDBStore.On("ListAvailableRooms", mock.Anything, arg).
			Return([]db.Room{}, nil).
			Once()

		// execute method
		resultRooms, err := ts.ListAvailableRooms(int(arg.Limit), int(arg.Offset), rsv.StartDate, rsv.EndDate)

		// tesify
		assert.NoError(t, err)
		require.Len(t, resultRooms, 0)
	})

	t.Run("Test Error", func(t *testing.T) {
		// create random reservation with room data
		rsv := randomReservation()

		//create db stub call arguments
		arg := db.ListAvailableRoomsParams{}
		err := CopyStructDataToDBStruct(rsv, &arg)
		require.NoError(t, err)

		arg.Limit = LimitRoomsPerPage
		arg.Offset = 0

		// create a new server with mock database store
		ts := NewTestServer(t)

		// build stub
		ts.MockDBStore.On("ListAvailableRooms", mock.Anything, arg).
			Return([]db.Room{}, errors.New("any error")).
			Once()

		// execute method
		rooms, err := ts.ListAvailableRooms(int(arg.Limit), int(arg.Offset), rsv.StartDate, rsv.EndDate)

		// tesify
		assert.Error(t, err)
		require.Len(t, rooms, 0)
	})
}

func TestServer_ListRooms(t *testing.T) {
	//create stub db call arguments
	arg := db.ListRoomsParams{
		Limit:  LimitRoomsPerPage,
		Offset: 0,
	}

	t.Run("Test All Rooms", func(t *testing.T) {
		// create stub return arguments
		rooms := make(Rooms, LimitRoomsPerPage)
		dbRooms := make([]db.Room, LimitRoomsPerPage)
		var err error

		for i := 0; i < LimitRoomsPerPage; i++ {
			rooms[i] = randomRoom()
			err = util.CopyDataUsingJSON(rooms[i], &dbRooms[i])
			require.NoError(t, err)
		}

		// create a new server with mock database store
		ts := NewTestServer(t)

		// build stub
		ts.MockDBStore.On("ListRooms", mock.Anything, arg).
			Return(dbRooms, nil).
			Once()

		// execute method
		returnedRooms, err := ts.ListRooms(int(arg.Limit), int(arg.Offset))

		// tesify
		assert.NoError(t, err)
		require.Len(t, returnedRooms, LimitRoomsPerPage)

		for i := 0; i < LimitRoomsPerPage; i++ {
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

		// create a new server with mock database store
		ts := NewTestServer(t)

		// build stub
		ts.MockDBStore.On("ListRooms", mock.Anything, arg).
			Return([]db.Room{}, nil).
			Once()

		// execute method
		ResultRooms, err := ts.ListRooms(int(arg.Limit), int(arg.Offset))

		// tesify
		assert.NoError(t, err)
		require.Len(t, ResultRooms, 0)
	})

	t.Run("Test Error", func(t *testing.T) {
		// create a new server with mock database store
		ts := NewTestServer(t)

		// build stub
		ts.MockDBStore.On("ListRooms", mock.Anything, arg).
			Return(nil, errors.New("any error")).
			Once()

		// execute method
		rooms, err := ts.ListRooms(int(arg.Limit), int(arg.Offset))

		// tesify
		assert.Error(t, err)
		require.Nil(t, rooms)
	})
}
