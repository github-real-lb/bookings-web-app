package db

import (
	"testing"
)

func createRandomReservation(t *testing.T) Reservation {
	r := Reservation{}

	return r

}

func TestQueries_CreateReservation(t *testing.T) {
	createRandomReservation(t)
}
