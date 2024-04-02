package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/github-real-lb/bookings-web-app/db"
	"github.com/github-real-lb/bookings-web-app/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestServer_MakeReservationHandler(t *testing.T) {
	// create reservation with random data to put in the session
	date := util.RandomDate()
	reservation := Reservation{
		StartDate: date,
		EndDate:   date.Add(time.Hour * 24 * 7),
		Room:      randomRoom(),
	}

	// create a test server
	testServer, _ := NewTestServer(t)

	// Test OK: reservation exists in session
	t.Run("OK", func(t *testing.T) {
		// create a new response and request
		recorder := httptest.NewRecorder()
		request := NewTestRequestWithSession(t, http.MethodGet, "/make-reservation", nil)

		// put reservation in session
		app.Session.Put(request.Context(), "reservation", reservation)

		// server the request
		testServer.Router.Handler.ServeHTTP(recorder, request)

		// remove reservation from session
		app.Session.Remove(request.Context(), "reservation")

		// testify
		assert.Equal(t, http.StatusOK, recorder.Code)
	})

	// Test Error: reservation missing from session
	t.Run("Error", func(t *testing.T) {
		// create a new response and request
		recorder := httptest.NewRecorder()
		request := NewTestRequestWithSession(t, http.MethodGet, "/make-reservation", nil)

		// server the request
		testServer.Router.Handler.ServeHTTP(recorder, request)

		// get error message from session
		errMsg := app.Session.PopString(request.Context(), "error")

		// testify
		assert.Equal(t, http.StatusTemporaryRedirect, recorder.Code)
		assert.Equal(t, "No reservation exists. Please make a reservation.", errMsg)
	})

}

func TestServer_PostMakeReservationHandler(t *testing.T) {
	// create initial reservation with random data to put in the session
	date := util.RandomDate()
	initialReservation := Reservation{
		StartDate: date,
		EndDate:   date.Add(time.Hour * 24 * 7),
		Room:      randomRoom(),
	}

	// create data for the body of the request
	data := make(map[string]string)
	data["first_name"] = util.RandomName()
	data["last_name"] = util.RandomName()
	data["email"] = util.RandomEmail()
	data["phone"] = util.RandomPhone()
	data["notes"] = util.RandomNote()

	// create the body of the request
	values := url.Values{}
	for key, value := range data {
		values.Set(key, value)
	}
	body := strings.NewReader(values.Encode())

	// create the final reservation that we are expected to get from the session
	finalReservation := initialReservation
	finalReservation.Unmarshal(data)
	err := finalReservation.GenerateReservationCode()
	require.NoError(t, err)

	// create a test server and mock database store
	testServer, mockStore := NewTestServer(t)

	// create mehod arguments
	var restrictionID int64 = 1

	arg := db.CreateReservationParams{}
	err = arg.Unmarshal(finalReservation.Marshal())
	require.NoError(t, err)

	//create method return arguments
	dbReservation := db.Reservation{}
	err = dbReservation.Unmarshal(finalReservation.Marshal())
	require.NoError(t, err)

	// build stub
	mockStore.On("CreateReservationTx", mock.Anything, arg, restrictionID).
		Return(dbReservation, nil).
		Once()

	// Test OK: reservation exists in session
	t.Run("OK", func(t *testing.T) {
		// create a new response and request
		recorder := httptest.NewRecorder()
		request := NewTestRequestWithSession(t, http.MethodPost, "/make-reservation", body)

		// put reservation in session
		app.Session.Put(request.Context(), "reservation", initialReservation)

		// server the request
		testServer.Router.Handler.ServeHTTP(recorder, request)

		// remove reservation from session
		app.Session.Remove(request.Context(), "reservation")

		// testify
		assert.Equal(t, http.StatusSeeOther, recorder.Code)
	})

}

// func TestPageHandlers(t *testing.T) {
// 	type params map[string]string

// 	tests := []struct {
// 		name   string // name of test
// 		method string // http.Method for the http.Request
// 		url    string // url for the http.Request
// 		params
// 		putSessionValues    func(t *testing.T, s *Server, r *http.Request)
// 		excpectedStatusCode int
// 	}{
// 		{"/ OK", http.MethodGet, "/", nil, nil, http.StatusOK},
// 		{"/about", http.MethodGet, "/about", nil, nil, http.StatusOK},

// 		{"rooms/list", http.MethodGet, "/rooms/list", nil, nil, http.StatusOK},
// 		{"rooms/none", http.MethodGet, "/rooms/none", nil, nil, http.StatusTemporaryRedirect},
// 		{"rooms/1", http.MethodGet, "/rooms/1", nil,
// 			func(t *testing.T, s *Server, r *http.Request) {
// 				rooms, err := s.ListRooms(5, 0)
// 				require.NoError(t, err)
// 				require.NotNil(t, rooms)

// 				app.Session.Put(r.Context(), "rooms", rooms)
// 			}, http.StatusSeeOther},
// 		{"rooms/room/exist", http.MethodGet, "/rooms/room/Generals-Quarters", nil,
// 			func(t *testing.T, s *Server, r *http.Request) {
// 				rooms, err := s.ListRooms(5, 0)
// 				require.NoError(t, err)
// 				require.NotNil(t, rooms)
// 				require.NotEmpty(t, rooms)

// 				room := rooms[0]
// 				app.Session.Put(r.Context(), "room", room)
// 			}, http.StatusOK},
// 		{"rooms/room/invalid", http.MethodGet, "/rooms/room/Generals-Quarters", nil, nil, http.StatusTemporaryRedirect},

// 		// {"Room Availability", http.MethodGet, "/search-room-availability", nil, http.StatusOK},

// 		{"Contact", http.MethodGet, "/contact", nil, nil, http.StatusOK},

// 		// {"Reservation", http.MethodGet, "/make-reservation", nil, http.StatusOK},
// 		// {"PostAvailability", http.MethodPost, "/search-availability",
// 		// 	params{
// 		// 		"start_date": "2024-05-01",
// 		// 		"end_date":   "2024-05-08",
// 		// 	}, http.StatusOK},
// 		// {"PostAvailabilityJSON", http.MethodPost, "/search-availability-json",
// 		// 	params{
// 		// 		"start_date": "2024-05-01",
// 		// 		"end_date":   "2024-05-08",
// 		// 	}, http.StatusOK},
// 		// {"PostReservation_OK", http.MethodPost, "/make-reservation",
// 		// 	params{
// 		// 		"first_name": "John",
// 		// 		"last_name":  "Dow",
// 		// 		"email":      "john.dow@gmail.com",
// 		// 		"phone":      "5555-5555",
// 		// 	}, http.StatusOK},
// 	}

// 	// load database/mock store
// 	store, err := db.NewPostgresDBStore(app.ConnectionString)
// 	require.NoError(t, err)

// 	// create new server
// 	server := NewServer(store)

// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			var body io.Reader = nil

// 			// encode parameters into body
// 			if test.params != nil {
// 				data := url.Values{}

// 				for key, value := range test.params {
// 					data.Set(key, value)
// 				}

// 				body = strings.NewReader(data.Encode())
// 			}

// 			// create a new response and request
// 			recorder := httptest.NewRecorder()
// 			request := NewTestRequestWithSession(t, test.method, test.url, body)

// 			// put values into session
// 			if test.putSessionValues != nil {
// 				test.putSessionValues(t, server, request)
// 			}

// 			// server HTTP
// 			server.Router.Handler.ServeHTTP(recorder, request)

// 			// assert
// 			assert.Equal(t, test.excpectedStatusCode, recorder.Code)
// 		})
// 	}
// }
