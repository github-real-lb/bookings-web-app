package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/github-real-lb/bookings-web-app/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPageHandlers(t *testing.T) {
	type params map[string]string

	tests := []struct {
		name   string // name of test
		method string // http.Method for the http.Request
		url    string // url for the http.Request
		params
		putSessionValues    func(t *testing.T, s *Server, r *http.Request)
		excpectedStatusCode int
	}{
		{"/ OK", http.MethodGet, "/", nil, nil, http.StatusOK},
		{"/about", http.MethodGet, "/about", nil, nil, http.StatusOK},

		{"rooms/list", http.MethodGet, "/rooms/list", nil, nil, http.StatusOK},
		{"rooms/none", http.MethodGet, "/rooms/none", nil, nil, http.StatusTemporaryRedirect},
		{"rooms/1", http.MethodGet, "/rooms/1", nil,
			func(t *testing.T, s *Server, r *http.Request) {
				rooms, err := s.ListRooms(5, 0)
				require.NoError(t, err)
				require.NotNil(t, rooms)

				app.Session.Put(r.Context(), "rooms", rooms)
			}, http.StatusSeeOther},
		{"rooms/room/exist", http.MethodGet, "/rooms/room/Generals-Quarters", nil,
			func(t *testing.T, s *Server, r *http.Request) {
				rooms, err := s.ListRooms(5, 0)
				require.NoError(t, err)
				require.NotNil(t, rooms)
				require.NotEmpty(t, rooms)

				room := rooms[0]
				app.Session.Put(r.Context(), "room", room)
			}, http.StatusOK},
		{"rooms/room/invalid", http.MethodGet, "/rooms/room/Generals-Quarters", nil, nil, http.StatusTemporaryRedirect},

		// {"Room Availability", http.MethodGet, "/search-room-availability", nil, http.StatusOK},

		{"Contact", http.MethodGet, "/contact", nil, nil, http.StatusOK},

		// {"Reservation", http.MethodGet, "/make-reservation", nil, http.StatusOK},
		// {"PostAvailability", http.MethodPost, "/search-availability",
		// 	params{
		// 		"start_date": "2024-05-01",
		// 		"end_date":   "2024-05-08",
		// 	}, http.StatusOK},
		// {"PostAvailabilityJSON", http.MethodPost, "/search-availability-json",
		// 	params{
		// 		"start_date": "2024-05-01",
		// 		"end_date":   "2024-05-08",
		// 	}, http.StatusOK},
		// {"PostReservation_OK", http.MethodPost, "/make-reservation",
		// 	params{
		// 		"first_name": "John",
		// 		"last_name":  "Dow",
		// 		"email":      "john.dow@gmail.com",
		// 		"phone":      "5555-5555",
		// 	}, http.StatusOK},
	}

	// load database/mock store
	store, err := db.NewPostgresDBStore(app.ConnectionString)
	require.NoError(t, err)

	// create new server
	server := NewServer(store)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var body io.Reader = nil

			// encode parameters into body
			if test.params != nil {
				data := url.Values{}

				for key, value := range test.params {
					data.Set(key, value)
				}

				body = strings.NewReader(data.Encode())
			}

			// create a new response and request
			recorder := httptest.NewRecorder()
			request := NewTestRequestWithSession(t, test.method, test.url, body)

			// put values into session
			if test.putSessionValues != nil {
				test.putSessionValues(t, server, request)
			}

			// server HTTP
			server.Router.Handler.ServeHTTP(recorder, request)

			// assert
			assert.Equal(t, test.excpectedStatusCode, recorder.Code)
		})
	}
}
