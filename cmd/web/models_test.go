package main

import (
	"testing"
	"time"

	"github.com/github-real-lb/bookings-web-app/util"
	"github.com/stretchr/testify/require"
)

func TestReservation_GenerateReservationCode(t *testing.T) {
	// N states the number of times to test randomness
	const N int = 10

	codes := make(util.KeysMap)
	for i := 0; i < N; i++ {
		date := util.RandomDate()

		r := Reservation{
			LastName:  util.RandomName(),
			StartDate: date,
			EndDate:   date.Add(time.Hour * 24 * 7),
		}

		r.GenerateReservationCode()
		util.RequireUnique(t, r.Code, codes)
	}
}

func TestRestriction_Scan(t *testing.T) {
	var r Restriction = RestrictionReservation

	err := r.Scan("owner_block")
	require.NoError(t, err)
	require.Equal(t, RestrictionOwnerBlock, r)

	err = r.Scan([]byte("reservation"))
	require.NoError(t, err)
	require.Equal(t, RestrictionReservation, r)

	err = r.Scan(Reservation{})
	require.ErrorContains(t, err, "unsupported scan type for Restriction:")
}
