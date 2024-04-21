package main

import (
	"context"
	"time"

	"github.com/github-real-lb/bookings-web-app/db"
)

const ContextTimeout = 3 * time.Second

// AuthenticateUser
func (s *Server) AuthenticateUser(email, password string) (User, error) {
	var user = User{}

	// create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), ContextTimeout)
	defer cancel()

	// authenticate user
	dbUser, err := s.DatabaseStore.AuthenticateUser(ctx, db.AuthenticateUserParams{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return user, err
	}

	err = CopyStructData(dbUser, &user)

	return user, err
}

// CheckRoomAvailability checks if room is available
func (s *Server) CheckRoomAvailability(roomID int64, startDate, endData time.Time) (bool, error) {
	// parse form's data to query arguments
	arg := db.CheckRoomAvailabilityParams{}
	arg.RoomID = roomID
	arg.StartDate.Scan(startDate)
	arg.EndDate.Scan(endData)

	// create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), ContextTimeout)
	defer cancel()

	// get list of availabe rooms
	return s.DatabaseStore.CheckRoomAvailability(ctx, arg)
}

// CreateReservation insert reservation data into database.
// it updates r with new data from database.
func (s *Server) CreateReservation(r *Reservation) error {
	// parse form's data to query arguments
	arg := db.CreateReservationParams{}
	err := arg.Unmarshal(r.Marshal())
	if err != nil {
		return err
	}

	// create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), ContextTimeout)
	defer cancel()

	// create new reservation
	dbReservation, err := s.DatabaseStore.CreateReservationTx(ctx, arg)
	if err != nil {
		return err
	}

	// update reservation with new data from database
	err = r.Unmarshal(dbReservation.Marshal())

	return err
}

// ListAvailableRooms returns limit amount of avaiable rooms for r, with the offset specified
func (s *Server) ListAvailableRooms(r Reservation, limit int, offset int) (Rooms, error) {
	// parse form's data to query arguments
	arg := db.ListAvailableRoomsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}
	err := arg.Unmarshal(r.Marshal())
	if err != nil {
		return nil, err
	}

	// create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), ContextTimeout)
	defer cancel()

	// get list of availabe rooms
	dbRooms, err := s.DatabaseStore.ListAvailableRooms(ctx, arg)
	if err != nil {
		return nil, err
	}

	l := len(dbRooms)
	if l == 0 {
		return Rooms{}, nil
	}

	rooms := make(Rooms, l)
	for i := 0; i < l; i++ {
		err = CopyStructData(dbRooms[i], &rooms[i])
		if err != nil {
			return nil, err
		}
	}

	return rooms, nil
}

func (s *Server) ListRooms(limit, offset int) (Rooms, error) {
	arg := db.ListRoomsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	// create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), ContextTimeout)
	defer cancel()

	// get list of availabe rooms
	dbRooms, err := s.DatabaseStore.ListRooms(ctx, arg)
	if err != nil {
		return nil, err
	}

	l := len(dbRooms)
	if l == 0 {
		return Rooms{}, nil
	}

	rooms := make(Rooms, l)
	for i := 0; i < l; i++ {
		err = CopyStructData(dbRooms[i], &rooms[i])
		if err != nil {
			return nil, err
		}
	}

	return rooms, nil
}
