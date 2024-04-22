package main

import (
	"context"
	"time"

	"github.com/github-real-lb/bookings-web-app/db"
	"github.com/github-real-lb/bookings-web-app/util"
	"github.com/jackc/pgx/v5/pgtype"
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

	err = util.CopyDataUsingJSON(dbUser, &user)

	return user, err
}

// CheckRoomAvailability checks if room is available
func (s *Server) CheckRoomAvailability(roomID int64, startDate, endData time.Time) (bool, error) {
	// parse form's data to query arguments
	arg := db.CheckRoomAvailabilityParams{}
	arg.RoomID = roomID
	err := arg.StartDate.Scan(startDate)
	if err != nil {
		return false, err
	}

	err = arg.EndDate.Scan(endData)
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
func (s *Server) CreateReservation(r Reservation) error {
	// create database transaction arguments
	arg := db.CreateReservationParams{}
	err := CopyStructDataToDBStruct(r, &arg)
	if err != nil {
		return err
	}

	// create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), ContextTimeout)
	defer cancel()

	// execute database transaction
	_, err = s.DatabaseStore.CreateReservationTx(ctx, arg)

	return err
}

// ListAvailableRooms returns limit amount of avaiable rooms for r, with the offset specified
func (s *Server) ListAvailableRooms(limit, offset int, startDate, endData time.Time) (Rooms, error) {
	// parse form's data to query arguments
	arg := db.ListAvailableRoomsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}
	err := arg.StartDate.Scan(startDate)
	if err != nil {
		return Rooms{}, err
	}

	err = arg.EndDate.Scan(endData)
	if err != nil {
		return Rooms{}, err
	}

	// create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), ContextTimeout)
	defer cancel()

	// get list of availabe rooms
	dbRooms, err := s.DatabaseStore.ListAvailableRooms(ctx, arg)
	if err != nil {
		return Rooms{}, err
	}

	l := len(dbRooms)
	if l == 0 {
		return Rooms{}, nil
	}

	rooms := make(Rooms, l)
	for i := 0; i < l; i++ {
		err = util.CopyDataUsingJSON(dbRooms[i], &rooms[i])
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
		err = util.CopyDataUsingJSON(dbRooms[i], &rooms[i])
		if err != nil {
			return nil, err
		}
	}

	return rooms, nil
}

// CopyStructDataToDBStruct uses json package to copy data from main package structs to db package structs.
// It is used to bridge differences of date and time implementations.
// target needs to be a pointer to a struct.
func CopyStructDataToDBStruct(src any, target any) error {
	// create intermediate map to get json data and manipulate differences between src and target
	intermediate, err := util.StructToMapUsingJSON(src)
	if err != nil {
		return err
	}

	// convert time.Time to pgtype.Date
	if v, ok := intermediate["start_date"].(time.Time); ok {
		var dbDate pgtype.Date
		err = dbDate.Scan(v)
		if err != nil {
			return err
		}
		intermediate["start_date"] = dbDate
	}

	// convert time.Time to pgtype.Date
	if v, ok := intermediate["end_date"].(time.Time); ok {
		var dbDate pgtype.Date
		err = dbDate.Scan(v)
		if err != nil {
			return err
		}
		intermediate["end_date"] = dbDate
	}

	return util.CopyDataUsingJSON(intermediate, target)
}
