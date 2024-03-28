package main

import (
	"context"
	"time"

	"github.com/github-real-lb/bookings-web-app/db"
	"github.com/github-real-lb/bookings-web-app/util"
)

// CreateReservation insert reservation data into database.
// it updates r with new data from database.
func (s *Server) CreateReservation(r *Reservation, restrictionID int64) error {
	// define create reservation parameters
	arg := db.CreateReservationParams{}
	err := arg.Unmarshal(r.Marshal())
	if err != nil {
		return err
	}

	// generate reservation code to add to parameters
	arg.Code, err = util.GenerateReservationCode(r.LastName, ReservationCodeLenght)
	if err != nil {
		return err
	}

	// create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// create new reservation
	reservation, err := s.DatabaseStore.CreateReservationTx(ctx, arg, restrictionID)
	if err != nil {
		return err
	}

	// update reservation with new data from database
	err = r.Unmarshal(reservation.Marshal())

	return err
}
