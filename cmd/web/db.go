package main

import (
	"context"
	"time"

	"github.com/github-real-lb/bookings-web-app/db"
)

const ContextTimeout = 3 * time.Second

// AuthenticateUser
func (s *Server) AuthenticateUser(u *User) error {
	// parse user email and password to query arguments
	arg := db.AuthenticateUserParams{}
	arg.Unmarshal(u.Marshal())

	// create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), ContextTimeout)
	defer cancel()

	// authenticate user
	dbUser, err := s.DatabaseStore.AuthenticateUser(ctx, arg)
	if err != nil {
		return err
	}

	// unmarhsal user data into User
	return u.Unmarshal(dbUser.Marshal())
}

// CheckRoomAvailability checks if room in reservation is available
func (s *Server) CheckRoomAvailability(r Reservation) (bool, error) {
	// parse form's data to query arguments
	var arg db.CheckRoomAvailabilityParams
	err := arg.Unmarshal(r.Marshal())
	if err != nil {
		return false, err
	}

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

	// unmarhsal list of availabe rooms into Rooms
	var rooms Rooms
	for i, v := range dbRooms {
		rooms = append(rooms, Room{})

		err = rooms[i].Unmarshal(v.Marshal())
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

	// unmarhsal list of availabe rooms into Rooms
	var rooms Rooms
	for i, v := range dbRooms {
		rooms = append(rooms, Room{})

		err = rooms[i].Unmarshal(v.Marshal())
		if err != nil {
			return nil, err
		}
	}

	return rooms, nil
}
