package main

import (
	"context"
	"time"

	"github.com/github-real-lb/bookings-web-app/db"
)

const ContextTimeout = 3 * time.Second

// AuthenticateUser authenticate the user email and password.
// If successful, it returns the id of the user and nil, otherwise a 0 and error
func (s *Server) AuthenticateUser(email, password string) (int64, error) {
	// create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), ContextTimeout)
	defer cancel()

	// authenticate user
	dbUser, err := s.DatabaseStore.AuthenticateUser(ctx, db.AuthenticateUserParams{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return 0, err
	}

	return dbUser.ID, err
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
	arg := db.CreateReservationParams{
		Code:      r.Code,
		FirstName: r.FirstName,
		LastName:  r.LastName,
		Email:     r.Email,
		RoomID:    r.RoomID,
	}
	arg.Phone.Scan(r.Phone)
	arg.StartDate.Scan(r.StartDate)
	arg.EndDate.Scan(r.EndDate)
	arg.Notes.Scan(r.Notes)

	// create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), ContextTimeout)
	defer cancel()

	// execute database transaction
	_, err := s.DatabaseStore.CreateReservationTx(ctx, arg)

	return err
}

// ListAvailableRooms returns limit amount of avaiable rooms in a date range, with the offset specified
func (s *Server) ListAvailableRooms(limit, offset int, startDate, endData time.Time) ([]Room, error) {
	// parse form's data to query arguments
	arg := db.ListAvailableRoomsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}
	err := arg.StartDate.Scan(startDate)
	if err != nil {
		return []Room{}, err
	}

	err = arg.EndDate.Scan(endData)
	if err != nil {
		return []Room{}, err
	}

	// create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), ContextTimeout)
	defer cancel()

	// get list of availabe rooms
	dbRooms, err := s.DatabaseStore.ListAvailableRooms(ctx, arg)
	if err != nil {
		return []Room{}, err
	}

	l := len(dbRooms)
	if l == 0 {
		return []Room{}, nil
	}

	rooms := make([]Room, l)
	for i := 0; i < l; i++ {
		rooms[i].Import(dbRooms[i])
	}

	return rooms, nil
}

// ListReservations returns limit amount of reservations, with the offset specified.
func (s *Server) ListReservations(limit, offset int) ([]Reservation, error) {
	arg := db.ListReservationsAndRoomsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	// create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), ContextTimeout)
	defer cancel()

	// get list of availabe rooms
	dbRsvs, err := s.DatabaseStore.ListReservationsAndRooms(ctx, arg)
	if err != nil {
		return nil, err
	}

	l := len(dbRsvs)
	if l == 0 {
		return []Reservation{}, nil
	}

	rsvs := make([]Reservation, l)
	for i := 0; i < l; i++ {
		rsvs[i].ImportWithRoom(dbRsvs[i])
	}

	return rsvs, nil
}

// ListRooms returns limit amount of rooms, with the offset specified
func (s *Server) ListRooms(limit, offset int) ([]Room, error) {
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
		return []Room{}, nil
	}

	rooms := make([]Room, l)
	for i := 0; i < l; i++ {
		rooms[i].Import(dbRooms[i])
	}

	return rooms, nil
}

// Import update r with the data from dbr
func (r *Reservation) Import(dbr db.Reservation) {
	r.ID = dbr.ID
	r.Code = dbr.Code
	r.FirstName = dbr.FirstName
	r.LastName = dbr.LastName
	r.Email = dbr.Email
	r.Phone = dbr.Phone.String
	r.StartDate = dbr.StartDate.Time
	r.EndDate = dbr.EndDate.Time
	r.RoomID = dbr.RoomID
	r.Notes = dbr.Notes.String
	r.CreatedAt = dbr.CreatedAt.Time
	r.UpdatedAt = dbr.UpdatedAt.Time
}

// Export update dbr with the data from r
func (r *Reservation) Export(dbr *db.Reservation) {
	dbr.ID = r.ID
	dbr.Code = r.Code
	dbr.FirstName = r.FirstName
	dbr.LastName = r.LastName
	dbr.Email = r.Email
	dbr.Phone.Scan(r.Phone)
	dbr.StartDate.Scan(r.StartDate)
	dbr.EndDate.Scan(r.EndDate)
	dbr.RoomID = r.RoomID
	dbr.Notes.Scan(r.Notes)
	dbr.CreatedAt.Scan(r.CreatedAt)
	dbr.UpdatedAt.Scan(r.UpdatedAt)
}

// ImportWithRoom update r with the data from dbr, imcluding the room data
func (r *Reservation) ImportWithRoom(dbr db.ListReservationsAndRoomsRow) {
	r.Import(dbr.Reservation)
	r.Room.Import(dbr.Room)
}

// ExportWithRoom update dbr with the data from r, imcluding the room data
func (r *Reservation) ExportWithRoom(dbr *db.ListReservationsAndRoomsRow) {
	r.Export(&dbr.Reservation)
	r.Room.Export(&dbr.Room)
}

// Import update r with the data from dbr
func (r *Room) Import(dbr db.Room) {
	r.ID = dbr.ID
	r.Name = dbr.Name
	r.Description = dbr.Description
	r.ImageFilename = dbr.ImageFilename
	r.CreatedAt = dbr.CreatedAt.Time
	r.UpdatedAt = dbr.UpdatedAt.Time
}

// Export update dbr with the data from r
func (r *Room) Export(dbr *db.Room) {
	dbr.ID = r.ID
	dbr.Name = r.Name
	dbr.Description = r.Description
	dbr.ImageFilename = r.ImageFilename
	dbr.CreatedAt.Scan(r.CreatedAt)
	dbr.UpdatedAt.Scan(r.UpdatedAt)
}
