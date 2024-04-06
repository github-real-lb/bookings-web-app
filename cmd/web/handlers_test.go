package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/github-real-lb/bookings-web-app/db"
	"github.com/github-real-lb/bookings-web-app/util"
	"github.com/github-real-lb/bookings-web-app/util/config"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestStaticPageHandlers(t *testing.T) {
	// create slice of test cases
	tests := []struct {
		name      string // name of test
		method    string // http.Method for the http.Request
		url       string // url for the http.Request
		excpected int    // expected status code
	}{
		{"home page", http.MethodGet, "/", http.StatusOK},
		{"/about page", http.MethodGet, "/about", http.StatusOK},
		{"/contact page", http.MethodGet, "/contact", http.StatusOK},
		{"/available-rooms-search page", http.MethodGet, "/available-rooms-search", http.StatusOK},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// create a new test server and request
			ts, _ := NewTestServer(t)
			req := ts.NewRequest(test.method, test.url, nil)

			//  server the request
			rr := ts.ServeRequest(req)

			// assert
			assert.Equal(t, test.excpected, rr.Code)
		})
	}
}

func TestServer_RoomsHandler(t *testing.T) {
	// test displaying the GET /rooms/list
	t.Run("OK List Rooms", func(t *testing.T) {
		// create a new test server, a mock database store and a request
		ts, mockStore := NewTestServer(t)
		req := ts.NewRequestWithSession(t, http.MethodGet, "/rooms/list", nil)

		// create mehod arguments
		arg := db.ListRoomsParams{
			Limit:  LimitRoomsPerPage,
			Offset: 0,
		}

		//create method return arguments
		dbRooms := make([]db.Room, LimitRoomsPerPage)
		for i := 0; i < LimitRoomsPerPage; i++ {
			dbRooms[i] = db.Room{
				ID:            util.RandomID(),
				Name:          util.RandomName(),
				Description:   util.RandomNote(),
				ImageFilename: fmt.Sprint(util.RandomName(), ".png"),
				CreatedAt: pgtype.Timestamptz{
					Time:  util.RandomDatetime(),
					Valid: true,
				},
				UpdatedAt: pgtype.Timestamptz{
					Time:  util.RandomDatetime(),
					Valid: true,
				},
			}
		}

		// build stub
		mockStore.On("ListRooms", mock.Anything, arg).
			Return(dbRooms, nil).
			Once()

		//  server the request
		rr := ts.ServeRequest(req)

		// checks rooms is in session and remove it
		rooms := app.Session.Pop(req.Context(), "rooms").(Rooms)
		require.Len(t, rooms, LimitRoomsPerPage)

		// testify
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	// test handling the GET /rooms/{room index}
	t.Run("OK Room Chosen", func(t *testing.T) {
		// create a new test server, and a new request
		ts, _ := NewTestServer(t)
		req := ts.NewRequestWithSession(t, http.MethodGet, "/rooms/1", nil)

		// create slice of Rooms to put in session
		rooms := randomRooms(LimitRoomsPerPage)

		// put rooms in session
		app.Session.Put(req.Context(), "rooms", rooms)

		// server the request
		rr := ts.ServeRequest(req)

		// remove rooms from session
		app.Session.Remove(req.Context(), "rooms")

		// check room is in session and removes it
		room := app.Session.Pop(req.Context(), "room").(Room)
		require.Equal(t, rooms[1], room)

		// testify
		assert.Equal(t, http.StatusSeeOther, rr.Code)
		ok := strings.HasPrefix(rr.Header().Get("Location"), "/rooms/room/")
		assert.True(t, ok)

	})

	// test handling the GET /rooms/{room index} with rooms missing from session
	t.Run("Error Missing Rooms", func(t *testing.T) {
		// create a new test server, and a new request
		ts, _ := NewTestServer(t)
		req := ts.NewRequestWithSession(t, http.MethodGet, "/rooms/1", nil)

		// server the request
		rr := ts.ServeRequest(req)

		// testify
		assert.Equal(t, http.StatusTemporaryRedirect, rr.Code)
		assert.Equal(t, "/rooms/list", rr.Header().Get("Location"))
	})

	// test handling the GET /rooms/{room index} with index not a number
	t.Run("Error Index Not a Number", func(t *testing.T) {
		// create a new test server, and a new request
		ts, _ := NewTestServer(t)
		req := ts.NewRequestWithSession(t, http.MethodGet, "/rooms/abc", nil)

		// create slice of Rooms to put in session
		rooms := randomRooms(LimitRoomsPerPage)

		// put rooms in session
		app.Session.Put(req.Context(), "rooms", rooms)

		// server the request
		rr := ts.ServeRequest(req)

		// checks rooms is in session and remove it
		sessionRooms := app.Session.Pop(req.Context(), "rooms").(Rooms)
		require.Len(t, sessionRooms, LimitRoomsPerPage)

		// testify
		assert.Equal(t, http.StatusTemporaryRedirect, rr.Code)
		assert.Equal(t, "/rooms/list", rr.Header().Get("Location"))
	})

	// test handling the GET /rooms/{room index} with index bigger than rooms lenght
	t.Run("Error Index Out of Scope", func(t *testing.T) {
		// create a new test server, and a new request
		ts, _ := NewTestServer(t)
		url := fmt.Sprint("/rooms/", LimitRoomsPerPage+10)
		req := ts.NewRequestWithSession(t, http.MethodGet, url, nil)

		// create slice of Rooms to put in session
		rooms := randomRooms(LimitRoomsPerPage)

		// put rooms in session
		app.Session.Put(req.Context(), "rooms", rooms)

		// server the request
		rr := ts.ServeRequest(req)

		// checks rooms is in session and remove it
		sessionRooms := app.Session.Pop(req.Context(), "rooms").(Rooms)
		require.Len(t, sessionRooms, LimitRoomsPerPage)

		// testify
		assert.Equal(t, http.StatusTemporaryRedirect, rr.Code)
		assert.Equal(t, "/rooms/list", rr.Header().Get("Location"))
	})
}

func TestServer_RoomHandler(t *testing.T) {
	// test displaying the GET /rooms/room/{name}
	t.Run("OK", func(t *testing.T) {
		// create a new test server and a request
		ts, _ := NewTestServer(t)
		req := ts.NewRequestWithSession(t, http.MethodGet, "/rooms/room/test", nil)

		// create room with random data to put in the session
		room := randomRoom()

		// put reservation in session
		app.Session.Put(req.Context(), "room", room)

		//  server the request
		rr := ts.ServeRequest(req)

		// remove reservation from session
		app.Session.Remove(req.Context(), "room")

		// testify
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	// test missing room from session
	t.Run("Missing Room from Session", func(t *testing.T) {
		// create a new test server, and a new request
		ts, _ := NewTestServer(t)
		req := ts.NewRequestWithSession(t, http.MethodGet, "/rooms/room/test", nil)

		//  server the request
		rr := ts.ServeRequest(req)

		// testify
		assert.Equal(t, http.StatusTemporaryRedirect, rr.Code)
		assert.Equal(t, "/rooms/list", rr.Header().Get("Location"))
	})
}

func TestServer_PostSearchRoomAvailabilityHandler(t *testing.T) {
	// Test OK: room is available
	t.Run("Room Available", func(t *testing.T) {
		// create room with random data to put in the session
		room := randomRoom()

		// creating dates for the request
		startDate := util.RandomDate()
		endDate := startDate.Add(time.Hour * 24 * 7)

		// create the body of the request
		values := url.Values{}
		values.Set("start_date", startDate.Format(config.DateLayout))
		values.Set("end_date", endDate.Format(config.DateLayout))
		body := strings.NewReader(values.Encode())

		// create a new test server, a mock database store and a request
		ts, mockStore := NewTestServer(t)
		req := ts.NewRequestWithSession(t, http.MethodPost, "/search-room-availability", body)

		// create mehod arguments
		reservation := Reservation{
			StartDate: startDate,
			EndDate:   endDate,
			Room:      room,
		}

		arg := db.CheckRoomAvailabilityParams{}
		err := arg.Unmarshal(reservation.Marshal())
		require.NoError(t, err)

		//build stub
		mockStore.On("CheckRoomAvailability", mock.Anything, arg).
			Return(true, nil).
			Once()

		// put room in session
		app.Session.Put(req.Context(), "room", room)

		//  server the request
		rr := ts.ServeRequest(req)

		// remove room from session
		app.Session.Remove(req.Context(), "room")

		// get the json response
		resp := SearchRoomAvailabilityResponse{}
		jsonResponseUnmarshal(t, rr, &resp)

		// check reservation is in session and removes it
		sessionReservation := app.Session.Pop(req.Context(), "reservation").(Reservation)
		require.WithinDuration(t, reservation.StartDate, sessionReservation.StartDate, time.Second)
		require.WithinDuration(t, reservation.EndDate, sessionReservation.EndDate, time.Second)
		require.Equal(t, reservation.Room, sessionReservation.Room)

		// testify
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.True(t, resp.OK)
		assert.Empty(t, resp.Message)
		assert.Empty(t, resp.Error)
	})

	// Test OK: room is available
	t.Run("Room Unavailable", func(t *testing.T) {
		// create room with random data to put in the session
		room := randomRoom()

		// creating dates for the request
		startDate := util.RandomDate()
		endDate := startDate.Add(time.Hour * 24 * 7)

		// create the body of the request
		values := url.Values{}
		values.Set("start_date", startDate.Format(config.DateLayout))
		values.Set("end_date", endDate.Format(config.DateLayout))
		body := strings.NewReader(values.Encode())

		// create a new test server, a mock database store and a request
		ts, mockStore := NewTestServer(t)
		req := ts.NewRequestWithSession(t, http.MethodPost, "/search-room-availability", body)

		// create mehod arguments
		reservation := Reservation{
			StartDate: startDate,
			EndDate:   endDate,
			Room:      room,
		}

		arg := db.CheckRoomAvailabilityParams{}
		err := arg.Unmarshal(reservation.Marshal())
		require.NoError(t, err)

		//build stub
		mockStore.On("CheckRoomAvailability", mock.Anything, arg).
			Return(false, nil).
			Once()

		// put room in session
		app.Session.Put(req.Context(), "room", room)

		//  server the request
		rr := ts.ServeRequest(req)

		// remove room from session
		app.Session.Remove(req.Context(), "room")

		// get the json response
		resp := SearchRoomAvailabilityResponse{}
		jsonResponseUnmarshal(t, rr, &resp)

		// testify
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.False(t, resp.OK)
		assert.Equal(t, "Room is unavailable. PLease try different dates.", resp.Message)
		assert.Empty(t, resp.Error)
	})

	// Test Error: room missing from session
	t.Run("Missing Room from Session", func(t *testing.T) {
		// create a new test server and a request
		ts, mockStore := NewTestServer(t)
		req := ts.NewRequestWithSession(t, http.MethodPost, "/search-room-availability", nil)

		// build stub
		mockStore.On("CheckRoomAvailability", mock.Anything, mock.Anything).
			Return(mock.Anything, mock.Anything).
			Times(0)

		//  server the request
		rr := ts.ServeRequest(req)

		// remove uncalled stub
		mockStore.On("CheckRoomAvailability", mock.Anything, mock.Anything).Unset()

		// get the json response
		resp := SearchRoomAvailabilityResponse{}
		jsonResponseUnmarshal(t, rr, &resp)

		// testify
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.False(t, resp.OK)
		assert.Empty(t, resp.Message)
		assert.Equal(t, "Internal Error. Please reload and try again.", resp.Error)
	})

	// Test Error: room exists in session but form is invalid
	t.Run("Invalid Form", func(t *testing.T) {
		// create room with random data to put in the session
		room := randomRoom()

		// create the dates to use in the tests
		date := util.RandomDate()
		startDate := date.Format(config.DateLayout)
		endDate := date.Add(-time.Hour * 24 * 7).Format(config.DateLayout)

		// create test cases for the form validation
		tests := []struct {
			Name string
			Data map[string]string
		}{
			{
				Name: "Missing Start Date",
				Data: map[string]string{
					"end_date": endDate,
				},
			},
			{
				Name: "Missing End Date",
				Data: map[string]string{
					"start_date": startDate,
				},
			},
			{
				Name: "End Date Prior to Start Date",
				Data: map[string]string{
					"start_date": startDate,
					"end_date":   endDate,
				},
			},
			{
				Name: "Invalid Start Date",
				Data: map[string]string{
					"start_date": util.RandomName(),
					"end_date":   endDate,
				},
			},
			{
				Name: "Invalid End Date",
				Data: map[string]string{
					"start_date": startDate,
					"end_date":   util.RandomName(),
				},
			},
		}

		// create a new test server and a mock database store
		ts, mockStore := NewTestServer(t)

		for _, v := range tests {
			t.Run(v.Name, func(t *testing.T) {
				// create the body of the request
				values := url.Values{}
				for key, value := range v.Data {
					values.Set(key, value)
				}
				body := strings.NewReader(values.Encode())

				// create a new request
				req := ts.NewRequestWithSession(t, http.MethodPost, "/search-room-availability", body)

				// build stub
				mockStore.On("CheckRoomAvailability", mock.Anything, mock.Anything).
					Return(mock.Anything, mock.Anything).
					Times(0)

				// put room in session
				app.Session.Put(req.Context(), "room", room)

				//  server the request
				rr := ts.ServeRequest(req)

				// remove room from session
				app.Session.Remove(req.Context(), "room")

				// remove uncalled stub
				mockStore.On("CheckRoomAvailability", mock.Anything, mock.Anything).Unset()

				// get the json response
				resp := SearchRoomAvailabilityResponse{}
				jsonResponseUnmarshal(t, rr, &resp)

				// testify
				assert.Equal(t, http.StatusOK, rr.Code)
				assert.False(t, resp.OK)
				assert.NotEmpty(t, resp.Message)
				assert.Empty(t, resp.Error)
			})
		}
	})

	// Test Error: room exists and form is invalid, but internal server error on CheckRoomAvailability
	t.Run("Internal Server Error", func(t *testing.T) {
		// create room with random data to put in the session
		room := randomRoom()

		// creating dates for the request
		startDate := util.RandomDate()
		endDate := startDate.Add(time.Hour * 24 * 7)

		// create the body of the request
		values := url.Values{}
		values.Set("start_date", startDate.Format(config.DateLayout))
		values.Set("end_date", endDate.Format(config.DateLayout))
		body := strings.NewReader(values.Encode())

		// create a new test server, a mock database store and a request
		ts, mockStore := NewTestServer(t)
		req := ts.NewRequestWithSession(t, http.MethodPost, "/search-room-availability", body)

		// create mehod arguments
		reservation := Reservation{
			StartDate: startDate,
			EndDate:   endDate,
			Room:      room,
		}

		arg := db.CheckRoomAvailabilityParams{}
		err := arg.Unmarshal(reservation.Marshal())
		require.NoError(t, err)

		//build stub
		mockStore.On("CheckRoomAvailability", mock.Anything, arg).
			Return(false, errors.New("internal server error")).
			Once()

		// put room in session
		app.Session.Put(req.Context(), "room", room)

		//  server the request
		rr := ts.ServeRequest(req)

		// remove room from session
		app.Session.Remove(req.Context(), "room")

		// get the json response
		resp := SearchRoomAvailabilityResponse{}
		jsonResponseUnmarshal(t, rr, &resp)

		// testify
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.False(t, resp.OK)
		assert.Empty(t, resp.Message)
		assert.Equal(t, "Internal Error. Please reload and try again.", resp.Error)
	})
}

// jsonResponseUnmarshal parses rr body and stores the result in the value pointed to by v.
// Any error is testified.
func jsonResponseUnmarshal(t *testing.T, rr *httptest.ResponseRecorder, v any) {
	// get the json response
	respBody, err := io.ReadAll(rr.Body)
	require.NoErrorf(t, err, "unable to read response body")

	err = json.Unmarshal(respBody, v)
	require.NoErrorf(t, err, "unable to unmarshal json response")
}

func TestServer_MakeReservationHandler(t *testing.T) {
	// Test OK: reservation exists in session
	t.Run("OK", func(t *testing.T) {
		// create a new test server, and a new request
		ts, _ := NewTestServer(t)
		req := ts.NewRequestWithSession(t, http.MethodGet, "/make-reservation", nil)

		// create reservation with random data to put in the session
		date := util.RandomDate()
		reservation := Reservation{
			StartDate: date,
			EndDate:   date.Add(time.Hour * 24 * 7),
			Room:      randomRoom(),
		}

		// put reservation in session
		app.Session.Put(req.Context(), "reservation", reservation)

		//  server the request
		rr := ts.ServeRequest(req)

		// remove reservation from session
		app.Session.Remove(req.Context(), "reservation")

		// testify
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	// Test Error: reservation missing from session
	t.Run("Error", func(t *testing.T) {
		// create a new test server, and a new request
		ts, _ := NewTestServer(t)
		req := ts.NewRequestWithSession(t, http.MethodGet, "/make-reservation", nil)

		//  server the request
		rr := ts.ServeRequest(req)

		// get error message from session and remove it
		errMsg := app.Session.PopString(req.Context(), "error")
		assert.Equal(t, "No reservation exists. Please make a reservation.", errMsg)

		// testify
		assert.Equal(t, http.StatusTemporaryRedirect, rr.Code)
		assert.Equal(t, "/", rr.Header().Get("Location"))
	})
}

func TestServer_PostMakeReservationHandler(t *testing.T) {
	// Test OK: reservation exists in session and form is valid
	t.Run("OK", func(t *testing.T) {
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

		// create a new test server, a mock database store and a request
		ts, mockStore := NewTestServer(t)
		req := ts.NewRequestWithSession(t, http.MethodPost, "/make-reservation", body)

		// create the final reservation that we are expected to get from the session
		finalReservation := initialReservation
		finalReservation.Unmarshal(data)
		err := finalReservation.GenerateReservationCode()
		require.NoError(t, err)

		// create mehod arguments
		arg := db.CreateReservationParams{}
		err = arg.Unmarshal(finalReservation.Marshal())
		require.NoError(t, err)

		//create method return arguments
		dbReservation := db.Reservation{}
		err = dbReservation.Unmarshal(finalReservation.Marshal())
		require.NoError(t, err)

		// build stub
		mockStore.On("CreateReservationTx", mock.Anything, arg).
			Return(dbReservation, nil).
			Once()

		// put reservation in session
		app.Session.Put(req.Context(), "reservation", initialReservation)

		//  server the request
		rr := ts.ServeRequest(req)

		// check reservation is in session and removes it
		sessionReservation := app.Session.Pop(req.Context(), "reservation").(Reservation)
		require.NotEmpty(t, sessionReservation.Code)
		require.Equal(t, finalReservation.FirstName, sessionReservation.FirstName)
		require.Equal(t, finalReservation.LastName, sessionReservation.LastName)
		require.Equal(t, finalReservation.Email, sessionReservation.Email)
		require.Equal(t, finalReservation.Phone, sessionReservation.Phone)
		require.WithinDuration(t, finalReservation.StartDate, sessionReservation.StartDate, time.Second)
		require.WithinDuration(t, finalReservation.EndDate, sessionReservation.EndDate, time.Second)
		require.Equal(t, finalReservation.Room, sessionReservation.Room)
		require.Equal(t, finalReservation.Notes, sessionReservation.Notes)

		// testify
		assert.Equal(t, http.StatusSeeOther, rr.Code)
		assert.Equal(t, "/reservation-summary", rr.Header().Get("Location"))
	})

	// Test Error: reservation missing from session
	t.Run("Missing Reservation from Session", func(t *testing.T) {
		// create a new test server, a mock database store and a request
		ts, mockStore := NewTestServer(t)
		req := ts.NewRequestWithSession(t, http.MethodPost, "/make-reservation", nil)

		// build stub
		mockStore.On("CreateReservationTx", mock.Anything, mock.Anything).
			Return(mock.Anything, mock.Anything).
			Times(0)

		//  server the request
		rr := ts.ServeRequest(req)

		// remove uncalled stub
		mockStore.On("CreateReservationTx", mock.Anything, mock.Anything).Unset()

		// get error message from session and removes it
		errMsg := app.Session.PopString(req.Context(), "error")
		assert.Equal(t, "No reservation exists. Please make a reservation.", errMsg)

		// testify
		assert.Equal(t, http.StatusTemporaryRedirect, rr.Code)
		assert.Equal(t, "/", rr.Header().Get("Location"))

	})

	// Test Error: reservation exists in session but form is invalid
	t.Run("Invalid Form", func(t *testing.T) {
		// create initial reservation with random data to put in the session
		date := util.RandomDate()
		initialReservation := Reservation{
			StartDate: date,
			EndDate:   date.Add(time.Hour * 24 * 7),
			Room:      randomRoom(),
		}

		firstName := util.RandomName()
		lastName := util.RandomName()
		email := util.RandomEmail()

		// create test cases for the form validation
		tests := []struct {
			Name string
			Data map[string]string
		}{
			{
				Name: "Missing First Name",
				Data: map[string]string{
					"last_name": lastName,
					"email":     email,
				},
			},
			{
				Name: "Missing Last Name",
				Data: map[string]string{
					"first_name": firstName,
					"email":      email,
				},
			},
			{
				Name: "Missing Email",
				Data: map[string]string{
					"first_name": firstName,
					"last_name":  lastName,
				},
			},
			{
				Name: "Invalid First Name",
				Data: map[string]string{
					"first_name": "x",
					"last_name":  lastName,
					"email":      email,
				},
			},
			{
				Name: "Invalid Last Name",
				Data: map[string]string{
					"first_name": firstName,
					"last_name":  "x",
					"email":      email,
				},
			},
			{
				Name: "Invalid Email",
				Data: map[string]string{
					"first_name": firstName,
					"last_name":  lastName,
					"email":      "x",
				},
			},
		}

		// create a new test server and a mock database store
		ts, mockStore := NewTestServer(t)

		for _, v := range tests {
			t.Run(v.Name, func(t *testing.T) {
				// create the body of the request
				values := url.Values{}
				for key, value := range v.Data {
					values.Set(key, value)
				}
				body := strings.NewReader(values.Encode())

				// create a new request
				req := ts.NewRequestWithSession(t, http.MethodPost, "/make-reservation", body)

				// build stub
				mockStore.On("CreateReservationTx", mock.Anything, mock.Anything).
					Return(mock.Anything, mock.Anything).
					Times(0)

				// put reservation in session
				app.Session.Put(req.Context(), "reservation", initialReservation)

				//  server the request
				rr := ts.ServeRequest(req)

				// remove reservation from session
				app.Session.Remove(req.Context(), "reservation")

				// remove uncalled stub
				mockStore.On("CreateReservationTx", mock.Anything, mock.Anything).Unset()

				// testify
				assert.Equal(t, http.StatusOK, rr.Code)
			})
		}
	})

	// Test Error: reservation exists and form it valid, but internal server error on CreateReservation
	//TODO:
	t.Run("Internal Server Error", func(t *testing.T) {
		log.Println("*********   TODO   *********")
	})
}
