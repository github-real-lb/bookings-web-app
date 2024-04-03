package main

import (
	"context"
	"time"

	"github.com/github-real-lb/bookings-web-app/db"
)

const ContextTimeout = 3 * time.Second

// CheckRoomAvailability checks if room in reservation is available
func (s *Server) CheckRoomAvailability(r Reservation) (bool, error) {
	// parse form's data to query arguments
	var arg db.CheckRoomAvailabilityParams
	arg.Unmarshal(r.Marshal())

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
	reservation, err := s.DatabaseStore.CreateReservationTx(ctx, arg)
	if err != nil {
		return err
	}

	// update reservation with new data from database
	err = r.Unmarshal(reservation.Marshal())

	return err
}

func (s *Server) ListAvailableRooms(r Reservation) (Rooms, error) {
	// parse form's data to query arguments
	var arg db.ListAvailableRoomsParams
	arg.Unmarshal(r.Marshal())

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
