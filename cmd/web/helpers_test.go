package main

import (
	"net/http"
	"testing"
	"time"

	"github.com/github-real-lb/bookings-web-app/util"
	"github.com/stretchr/testify/assert"
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

func TestIsAuthenticated(t *testing.T) {
	ts := NewTestServer(t)
	req := ts.NewRequestWithSession(t, http.MethodGet, "/", nil)

	t.Run("Not Authenticated", func(t *testing.T) {
		result := IsAuthenticated(req)
		assert.False(t, result)
	})

	t.Run("Authenticated", func(t *testing.T) {
		app.Session.Put(req.Context(), "user_id", 1)
		result := IsAuthenticated(req)
		assert.True(t, result)
	})
}
