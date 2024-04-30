package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/github-real-lb/bookings-web-app/db"
	"github.com/github-real-lb/bookings-web-app/util"
	"github.com/github-real-lb/bookings-web-app/util/config"
	"github.com/github-real-lb/bookings-web-app/util/forms"
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
			ts := NewTestServer(t)
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
	t.Run("OK List []Room", func(t *testing.T) {
		// create a new test server, a mock database store and a request
		ts := NewTestServer(t)
		req := ts.NewRequestWithSession(t, http.MethodGet, "/rooms/list", nil)

		// screate stub call arguments
		arg := db.ListRoomsParams{
			Limit:  LimitRoomsPerPage,
			Offset: 0,
		}

		//create stub return arguments
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

		// build stubs
		ts.MockDBStore.On("ListRooms", mock.Anything, arg).
			Return(dbRooms, nil).
			Once()

		//  server the request
		rr := ts.ServeRequest(req)

		// checks rooms is in session and remove it
		rooms := app.Session.Pop(req.Context(), "rooms").([]Room)
		require.Len(t, rooms, LimitRoomsPerPage)

		// testify
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	// test handling the GET /rooms/{room index}
	t.Run("OK Room Chosen", func(t *testing.T) {
		// create a new test server, and a new request
		ts := NewTestServer(t)
		req := ts.NewRequestWithSession(t, http.MethodGet, "/rooms/1", nil)

		// create slice of []Room to put in session
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

	// test database error while displaying the GET /rooms/list
	t.Run("Error In DB", func(t *testing.T) {
		// create a new test server, a mock database store and a request
		ts := NewTestServer(t)
		req := ts.NewRequestWithSession(t, http.MethodGet, "/rooms/list", nil)

		// create stubs arguments
		arg := db.ListRoomsParams{
			Limit:  LimitRoomsPerPage,
			Offset: 0,
		}

		err := errors.New("any error")

		sErr := ServerError{
			Prompt: "Unable to load rooms from database.",
			URL:    req.URL.Path,
			Err:    err,
		}

		// build stubs
		ts.MockDBStore.On("ListRooms", mock.Anything, arg).
			Return(nil, err).
			Once()
		ts.BuildLogErrorStub(sErr)

		//  server the request
		rr := ts.ServeRequest(req)

		// get error message from session and remove it
		errMsg := app.Session.PopString(req.Context(), "error")
		assert.Equal(t, sErr.Prompt, errMsg)

		// testify
		assert.Equal(t, http.StatusTemporaryRedirect, rr.Code)
		assert.Equal(t, "/", rr.Header().Get("Location"))
	})

	// test handling the GET /rooms/{room index} with rooms missing from session
	t.Run("Error Missing []Room", func(t *testing.T) {
		// create a new test server, and a new request
		ts := NewTestServer(t)
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
		ts := NewTestServer(t)
		req := ts.NewRequestWithSession(t, http.MethodGet, "/rooms/abc", nil)

		// create slice of []Room to put in session
		rooms := randomRooms(LimitRoomsPerPage)

		// build stub
		ts.BuildLogAnyErrorStub()

		// put rooms in session
		app.Session.Put(req.Context(), "rooms", rooms)

		// server the request
		rr := ts.ServeRequest(req)

		// checks rooms is in session and remove it
		sessionRooms := app.Session.Pop(req.Context(), "rooms").([]Room)
		require.Len(t, sessionRooms, LimitRoomsPerPage)

		// testify
		assert.Equal(t, http.StatusTemporaryRedirect, rr.Code)
		assert.Equal(t, "/rooms/list", rr.Header().Get("Location"))
	})

	// test handling the GET /rooms/{room index} with index bigger than rooms lenght
	t.Run("Error Index Out of Scope", func(t *testing.T) {
		// create a new test server, and a new request
		ts := NewTestServer(t)
		url := fmt.Sprint("/rooms/", LimitRoomsPerPage+10)
		req := ts.NewRequestWithSession(t, http.MethodGet, url, nil)

		// create slice of []Room to put in session
		rooms := randomRooms(LimitRoomsPerPage)

		// put rooms in session
		app.Session.Put(req.Context(), "rooms", rooms)

		// server the request
		rr := ts.ServeRequest(req)

		// checks rooms is in session and remove it
		sessionRooms := app.Session.Pop(req.Context(), "rooms").([]Room)
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
		ts := NewTestServer(t)
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
		ts := NewTestServer(t)
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
		values := url.Values{
			"start_date": {startDate.Format(config.DateLayout)},
			"end_date":   {endDate.Format(config.DateLayout)},
		}
		body := strings.NewReader(values.Encode())

		// create a new test server, a mock database store and a request
		ts := NewTestServer(t)
		req := ts.NewRequestWithSession(t, http.MethodPost, "/search-room-availability", body)

		// screate stub call arguments
		rsv := Reservation{
			StartDate: startDate,
			EndDate:   endDate,
			RoomID:    room.ID,
			Room:      room,
		}

		arg := db.CheckRoomAvailabilityParams{
			RoomID: rsv.RoomID,
		}
		arg.StartDate.Scan(rsv.StartDate)
		arg.EndDate.Scan(rsv.EndDate)

		//build stub
		ts.MockDBStore.On("CheckRoomAvailability", mock.Anything, arg).
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
		scsRsv := app.Session.Pop(req.Context(), "reservation").(Reservation)
		require.WithinDuration(t, rsv.StartDate, scsRsv.StartDate, time.Second)
		require.WithinDuration(t, rsv.EndDate, scsRsv.EndDate, time.Second)
		require.Equal(t, rsv.RoomID, scsRsv.RoomID)
		require.Equal(t, rsv.Room, scsRsv.Room)

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
		values := url.Values{
			"start_date": {startDate.Format(config.DateLayout)},
			"end_date":   {endDate.Format(config.DateLayout)},
		}
		body := strings.NewReader(values.Encode())

		// create a new test server, a mock database store and a request
		ts := NewTestServer(t)
		req := ts.NewRequestWithSession(t, http.MethodPost, "/search-room-availability", body)

		// screate stub call arguments
		rsv := Reservation{
			StartDate: startDate,
			EndDate:   endDate,
			RoomID:    room.ID,
			Room:      room,
		}

		arg := db.CheckRoomAvailabilityParams{
			RoomID: rsv.RoomID,
		}
		arg.StartDate.Scan(rsv.StartDate)
		arg.EndDate.Scan(rsv.EndDate)

		//build stub
		ts.MockDBStore.On("CheckRoomAvailability", mock.Anything, arg).
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
		ts := NewTestServer(t)
		req := ts.NewRequestWithSession(t, http.MethodPost, "/search-room-availability", nil)

		// build stub
		sErr := ServerError{
			Prompt: "Unable to get room from session.",
			URL:    req.URL.Path,
			Err:    errors.New("wrong routing"),
		}

		ts.BuildLogErrorStub(sErr)

		//  server the request
		rr := ts.ServeRequest(req)

		// get the json response
		resp := SearchRoomAvailabilityResponse{}
		jsonResponseUnmarshal(t, rr, &resp)

		// testify
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.False(t, resp.OK)
		assert.Empty(t, resp.Message)
		assert.Equal(t, "Internal Error. Please reload and try again.", resp.Error)
	})

	// Test Error: invalid body data cause error in ParseForm()
	t.Run("Ivalid Body Data", func(t *testing.T) {
		// create room with random data to put in the session
		room := randomRoom()

		// creating invalid body
		body := strings.NewReader("%^")

		// create a new test server and a request
		ts := NewTestServer(t)
		req := ts.NewRequestWithSession(t, http.MethodPost, "/search-room-availability", body)

		// build stub
		ts.BuildLogAnyErrorStub()

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
			Name   string
			Values url.Values
		}{
			{
				Name: "Missing Start Date",
				Values: url.Values{
					"end_date": {endDate},
				},
			},
			{
				Name: "Missing End Date",
				Values: url.Values{
					"start_date": {startDate},
				},
			},
			{
				Name: "End Date Prior to Start Date",
				Values: url.Values{
					"start_date": {startDate},
					"end_date":   {endDate},
				},
			},
			{
				Name: "Invalid Start Date",
				Values: url.Values{
					"start_date": {util.RandomName()},
					"end_date":   {endDate},
				},
			},
			{
				Name: "Invalid End Date",
				Values: url.Values{
					"start_date": {startDate},
					"end_date":   {util.RandomName()},
				},
			},
		}

		// create a new test server and a mock database store
		ts := NewTestServer(t)

		for _, test := range tests {
			t.Run(test.Name, func(t *testing.T) {
				// create the body of the request
				body := strings.NewReader(test.Values.Encode())

				// create a new request
				req := ts.NewRequestWithSession(t, http.MethodPost, "/search-room-availability", body)

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
		values := url.Values{
			"start_date": {startDate.Format(config.DateLayout)},
			"end_date":   {endDate.Format(config.DateLayout)},
		}
		body := strings.NewReader(values.Encode())

		// create a new test server, a mock database store and a request
		ts := NewTestServer(t)
		req := ts.NewRequestWithSession(t, http.MethodPost, "/search-room-availability", body)

		// screate stub call arguments
		rsv := Reservation{
			StartDate: startDate,
			EndDate:   endDate,
			RoomID:    room.ID,
			Room:      room,
		}

		arg := db.CheckRoomAvailabilityParams{
			RoomID: rsv.RoomID,
		}
		arg.StartDate.Scan(rsv.StartDate)
		arg.EndDate.Scan(rsv.EndDate)

		err := errors.New("any error")

		sErr := ServerError{
			Prompt: "Unable to check room availability.",
			URL:    req.URL.Path,
			Err:    err,
		}

		//build stubs
		ts.MockDBStore.On("CheckRoomAvailability", mock.Anything, arg).
			Return(false, err).
			Once()
		ts.BuildLogErrorStub(sErr)

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
	err := json.Unmarshal(rr.Body.Bytes(), v)
	require.NoErrorf(t, err, "unable to unmarshal json response")
}

func TestServer_PostAvailableRoomsSearchHandler(t *testing.T) {
	// Test OK: no rooms available for form request dates
	t.Run("No Available []Room", func(t *testing.T) {
		// creating dates for the request
		startDate := util.RandomDate()
		endDate := startDate.Add(time.Hour * 24 * 7)

		// create the body of the request
		values := url.Values{
			"start_date": {startDate.Format(config.DateLayout)},
			"end_date":   {endDate.Format(config.DateLayout)},
		}
		body := strings.NewReader(values.Encode())

		// create a new test server, a mock database store and a request
		ts := NewTestServer(t)
		req := ts.NewRequestWithSession(t, http.MethodPost, "/available-rooms-search", body)

		// screate stub call arguments
		rsv := Reservation{
			StartDate: startDate,
			EndDate:   endDate,
		}

		arg := db.ListAvailableRoomsParams{
			Limit:  LimitRoomsPerPage,
			Offset: 0,
		}
		err := arg.StartDate.Scan(rsv.StartDate)
		require.NoError(t, err)
		err = arg.EndDate.Scan(rsv.EndDate)
		require.NoError(t, err)

		//build stub
		ts.MockDBStore.On("ListAvailableRooms", mock.Anything, arg).
			Return([]db.Room{}, nil).
			Once()

		//  server the request
		rr := ts.ServeRequest(req)

		// testify
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), `message: "No rooms are available. Please try different dates."`)
	})

	// Test OK: rooms are available
	t.Run("Available []Room", func(t *testing.T) {
		// creating dates for the request
		startDate := util.RandomDate()
		endDate := startDate.Add(time.Hour * 24 * 7)

		// create the body of the request
		values := url.Values{
			"start_date": {startDate.Format(config.DateLayout)},
			"end_date":   {endDate.Format(config.DateLayout)},
		}
		body := strings.NewReader(values.Encode())

		// create a new test server, a mock database store and a request
		ts := NewTestServer(t)
		req := ts.NewRequestWithSession(t, http.MethodPost, "/available-rooms-search", body)

		// screate stub call arguments
		rsv := Reservation{
			StartDate: startDate,
			EndDate:   endDate,
		}

		arg := db.ListAvailableRoomsParams{
			Limit:  LimitRoomsPerPage,
			Offset: 0,
		}
		err := arg.StartDate.Scan(rsv.StartDate)
		require.NoError(t, err)
		err = arg.EndDate.Scan(rsv.EndDate)
		require.NoError(t, err)

		// create stub return arguments
		rooms := make([]Room, LimitRoomsPerPage)
		dbRooms := make([]db.Room, LimitRoomsPerPage)

		for i := 0; i < LimitRoomsPerPage; i++ {
			rooms[i] = randomRoom()
			rooms[i].Export(&dbRooms[i])
		}

		//build stub
		ts.MockDBStore.On("ListAvailableRooms", mock.Anything, arg).
			Return(dbRooms, nil).
			Once()

		//  server the request
		rr := ts.ServeRequest(req)

		// get reservation from session and remove it
		sessionRsv := app.Session.Pop(req.Context(), "reservation").(Reservation)
		assert.WithinDuration(t, rsv.StartDate, sessionRsv.StartDate, time.Second)
		assert.WithinDuration(t, rsv.EndDate, sessionRsv.EndDate, time.Second)

		// get rooms from session and remove it
		sessionRooms := app.Session.Pop(req.Context(), "rooms").([]Room)
		assert.Len(t, sessionRooms, LimitRoomsPerPage)

		// testify
		assert.Equal(t, http.StatusSeeOther, rr.Code)
		assert.Equal(t, "/available-rooms/available", rr.Header().Get("Location"))
	})

	// Test Error: invalid body data cause error in ParseForm()
	t.Run("Ivalid Body Data", func(t *testing.T) {
		// creating invalid body
		body := strings.NewReader("%^")

		// create a new test server and a request
		ts := NewTestServer(t)
		req := ts.NewRequestWithSession(t, http.MethodPost, "/available-rooms-search", body)

		// build stub
		ts.BuildLogAnyErrorStub()

		//  server the request
		rr := ts.ServeRequest(req)

		// get error message from session and remove it
		errMsg := app.Session.PopString(req.Context(), "error")

		sErr := CreateServerError(ErrorParseForm, req.URL.Path, nil)
		assert.Equal(t, sErr.Prompt, errMsg)

		// testify
		assert.Equal(t, http.StatusTemporaryRedirect, rr.Code)
		assert.Equal(t, "/available-rooms-search", rr.Header().Get("Location"))
	})

	// Test Error: invalid form data in post request
	t.Run("Invalid Form", func(t *testing.T) {
		// create the dates to use in the tests
		date := util.RandomDate()
		startDate := date.Format(config.DateLayout)
		endDate := date.Add(-time.Hour * 24 * 7).Format(config.DateLayout)

		// create test cases for the form validation
		tests := []struct {
			Name   string
			Values url.Values
		}{
			{
				Name: "Missing Start Date",
				Values: url.Values{
					"end_date": {endDate},
				},
			},
			{
				Name: "Missing End Date",
				Values: url.Values{
					"start_date": {startDate},
				},
			},
			{
				Name: "End Date Prior to Start Date",
				Values: url.Values{
					"start_date": {startDate},
					"end_date":   {endDate},
				},
			},
			{
				Name: "Invalid Start Date",
				Values: url.Values{
					"start_date": {util.RandomName()},
					"end_date":   {endDate},
				},
			},
			{
				Name: "Invalid End Date",
				Values: url.Values{
					"start_date": {startDate},
					"end_date":   {util.RandomName()},
				},
			},
		}

		// create a new test server and a mock database store
		ts := NewTestServer(t)

		for _, test := range tests {
			t.Run(test.Name, func(t *testing.T) {
				// create the body of the request
				body := strings.NewReader(test.Values.Encode())

				// create a new request
				req := ts.NewRequestWithSession(t, http.MethodPost, "/available-rooms-search", body)

				//  server the request
				rr := ts.ServeRequest(req)

				// testify
				assert.Equal(t, http.StatusOK, rr.Code)
			})
		}
	})

	// Test Error: form is invalid, but internal server error on ListAvailableRooms
	t.Run("Internal Server Error", func(t *testing.T) {
		// creating dates for the request
		startDate := util.RandomDate()
		endDate := startDate.Add(time.Hour * 24 * 7)

		// create the body of the request
		values := url.Values{
			"start_date": {startDate.Format(config.DateLayout)},
			"end_date":   {endDate.Format(config.DateLayout)},
		}
		body := strings.NewReader(values.Encode())

		// create a new test server, a mock database store and a request
		ts := NewTestServer(t)
		req := ts.NewRequestWithSession(t, http.MethodPost, "/available-rooms-search", body)

		// screate stub call arguments
		rsv := Reservation{
			StartDate: startDate,
			EndDate:   endDate,
		}

		arg := db.ListAvailableRoomsParams{
			Limit:  LimitRoomsPerPage,
			Offset: 0,
		}
		err := arg.StartDate.Scan(rsv.StartDate)
		require.NoError(t, err)
		err = arg.EndDate.Scan(rsv.EndDate)
		require.NoError(t, err)

		err = errors.New("any error")

		sErr := ServerError{
			Prompt: "Unable to load available rooms.",
			URL:    req.URL.Path,
			Err:    err,
		}

		//build stub
		ts.MockDBStore.On("ListAvailableRooms", mock.Anything, arg).
			Return(nil, err).
			Once()
		ts.BuildLogErrorStub(sErr)

		//  server the request
		rr := ts.ServeRequest(req)

		// get error message from session and remove it
		errMsg := app.Session.PopString(req.Context(), "error")
		assert.Equal(t, sErr.Prompt, errMsg)

		// testify
		assert.Equal(t, http.StatusTemporaryRedirect, rr.Code)
		assert.Equal(t, "/", rr.Header().Get("Location"))
	})
}

func TestServer_AvailableRoomsListHandler(t *testing.T) {

	// Test Error: available rooms are missing from session
	t.Run("Error Missing Available []Room", func(t *testing.T) {
		// create a new test server, and a new request
		ts := NewTestServer(t)
		req := ts.NewRequestWithSession(t, http.MethodGet, "/available-rooms/available", nil)

		// build stub
		sErr := CreateServerError(ErrorMissingReservation, req.URL.Path, nil)
		ts.BuildLogErrorStub(sErr)

		//  server the request
		rr := ts.ServeRequest(req)

		// get error message from session and remove it
		errMsg := app.Session.PopString(req.Context(), "error")
		assert.Equal(t, sErr.Prompt, errMsg)

		// testify
		assert.Equal(t, http.StatusTemporaryRedirect, rr.Code)
		assert.Equal(t, "/available-rooms-search", rr.Header().Get("Location"))
	})

	// Test OK: available rooms are missing from session
	t.Run("OK Available []Room", func(t *testing.T) {
		//create rooms slice with random data of n rooms
		const N = 5
		rooms := randomRooms(N)

		// create a new test server, and a new request
		ts := NewTestServer(t)
		req := ts.NewRequestWithSession(t, http.MethodGet, "/available-rooms/available", nil)

		// put rooms in session
		app.Session.Put(req.Context(), "rooms", rooms)

		//  server the request
		rr := ts.ServeRequest(req)

		// remove rooms from session
		app.Session.Remove(req.Context(), "rooms")

		// testify
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	// test handling the GET /available-rooms/{index} with index not a number
	t.Run("Error Index Not a Number", func(t *testing.T) {
		//create rooms slice with random data of n rooms
		const N = 5
		rooms := randomRooms(N)

		// create a new test server, and a new request
		ts := NewTestServer(t)
		req := ts.NewRequestWithSession(t, http.MethodGet, "/available-rooms/abc", nil)

		// build stub
		ts.BuildLogAnyErrorStub()

		// put rooms in session
		app.Session.Put(req.Context(), "rooms", rooms)

		//  server the request
		rr := ts.ServeRequest(req)

		// remove rooms from session
		app.Session.Remove(req.Context(), "rooms")

		// testify
		assert.Equal(t, http.StatusTemporaryRedirect, rr.Code)
		assert.Equal(t, "/available-rooms/available", rr.Header().Get("Location"))
	})

	// Test Error: reservation data is missing from session
	t.Run("Error Missing Reservation Data", func(t *testing.T) {
		//create rooms slice with random data of n rooms
		const N = 5
		rooms := randomRooms(N)

		// create a new test server, and a new request
		ts := NewTestServer(t)
		req := ts.NewRequestWithSession(t, http.MethodGet, "/available-rooms/1", nil)

		// build stub
		sErr := CreateServerError(ErrorMissingReservation, req.URL.Path, nil)
		ts.BuildLogErrorStub(sErr)

		// put rooms in session
		app.Session.Put(req.Context(), "rooms", rooms)

		//  server the request
		rr := ts.ServeRequest(req)

		// remove rooms from session
		app.Session.Remove(req.Context(), "rooms")

		// get error message from session and remove it
		errMsg := app.Session.PopString(req.Context(), "error")
		assert.Equal(t, sErr.Prompt, errMsg)

		// testify
		assert.Equal(t, http.StatusTemporaryRedirect, rr.Code)
		assert.Equal(t, "/", rr.Header().Get("Location"))
	})

	// Test OK
	t.Run("OK", func(t *testing.T) {
		//create rooms slice with random data of n rooms
		const N = 5
		rooms := randomRooms(N)

		// create random reservation
		rsv := randomReservation()

		// create a new test server, and a new request
		ts := NewTestServer(t)
		req := ts.NewRequestWithSession(t, http.MethodGet, "/available-rooms/1", nil)

		// put rooms and reservation in session
		app.Session.Put(req.Context(), "rooms", rooms)
		app.Session.Put(req.Context(), "reservation", rsv)

		//  server the request
		rr := ts.ServeRequest(req)

		// remove rooms and reservation from session
		app.Session.Remove(req.Context(), "rooms")
		app.Session.Remove(req.Context(), "reservation")

		// testify
		assert.Equal(t, http.StatusSeeOther, rr.Code)
		assert.Equal(t, "/make-reservation", rr.Header().Get("Location"))
	})
}

func TestServer_MakeReservationHandler(t *testing.T) {
	// Test OK: reservation exists in session
	t.Run("OK", func(t *testing.T) {
		// create a new test server, and a new request
		ts := NewTestServer(t)
		req := ts.NewRequestWithSession(t, http.MethodGet, "/make-reservation", nil)

		// create reservation with random data to put in the session
		rDate := util.RandomDate()
		rRoom := randomRoom()

		rsv := Reservation{
			StartDate: rDate,
			EndDate:   rDate.Add(time.Hour * 24 * 7),
			RoomID:    rRoom.ID,
			Room:      rRoom,
		}

		// put reservation in session
		app.Session.Put(req.Context(), "reservation", rsv)

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
		ts := NewTestServer(t)
		req := ts.NewRequestWithSession(t, http.MethodGet, "/make-reservation", nil)

		// build stub
		sErr := CreateServerError(ErrorMissingReservation, req.URL.Path, nil)
		ts.BuildLogErrorStub(sErr)

		//  server the request
		rr := ts.ServeRequest(req)

		// get error message from session and remove it
		errMsg := app.Session.PopString(req.Context(), "error")
		assert.Equal(t, sErr.Prompt, errMsg)

		// testify
		assert.Equal(t, http.StatusTemporaryRedirect, rr.Code)
		assert.Equal(t, "/", rr.Header().Get("Location"))
	})
}

func TestServer_PostMakeReservationHandler(t *testing.T) {
	// create initial reservation with random data to put in the session
	rDate := util.RandomDate()
	rRoom := randomRoom()

	initRsv := Reservation{
		StartDate: rDate,
		EndDate:   rDate.Add(time.Hour * 24 * 7),
		RoomID:    rRoom.ID,
		Room:      rRoom,
	}

	// Test OK: reservation exists in session and form is valid
	t.Run("OK", func(t *testing.T) {
		// create the final reservation that we are expected to get from the session
		finalRsv := initRsv
		finalRsv.FirstName = util.RandomName()
		finalRsv.LastName = util.RandomName()
		finalRsv.Email = util.RandomEmail()
		finalRsv.Phone = util.RandomPhone()
		finalRsv.Notes = util.RandomNote()
		finalRsv.GenerateReservationCode()

		// create form data for the body of the request
		f := forms.New(nil)
		f.Add("first_name", finalRsv.FirstName)
		f.Add("last_name", finalRsv.LastName)
		f.Add("email", finalRsv.Email)
		f.Add("phone", finalRsv.Phone)
		f.Add("notes", finalRsv.Notes)

		// create the body of the request
		body := strings.NewReader(f.Encode())

		// create a new test server, a mock database store and a request
		ts := NewTestServer(t)
		req := ts.NewRequestWithSession(t, http.MethodPost, "/make-reservation", body)

		//create stub return arguments for CreateReservationTx
		dbRsv := db.Reservation{}
		finalRsv.Export(&dbRsv)

		// build stub for CreateReservationTx
		ts.MockDBStore.On("CreateReservationTx", mock.Anything, mock.Anything).
			Return(dbRsv, nil).
			Once()

		// build stubs for mailing and logging of mail sent to guest and admin
		ts.BuildSendAnyMailStub()
		ts.BuildLogAnyInfoStub()
		ts.BuildSendAnyMailStub()
		ts.BuildLogAnyInfoStub()

		// put reservation in session
		app.Session.Put(req.Context(), "reservation", initRsv)

		//  server the request
		rr := ts.ServeRequest(req)

		// check reservation is in session and removes it
		scsRsv := app.Session.Pop(req.Context(), "reservation").(Reservation)
		require.NotEmpty(t, scsRsv.Code)
		require.Equal(t, finalRsv.FirstName, scsRsv.FirstName)
		require.Equal(t, finalRsv.LastName, scsRsv.LastName)
		require.Equal(t, finalRsv.Email, scsRsv.Email)
		require.Equal(t, finalRsv.Phone, scsRsv.Phone)
		require.WithinDuration(t, finalRsv.StartDate, scsRsv.StartDate, time.Second)
		require.WithinDuration(t, finalRsv.EndDate, scsRsv.EndDate, time.Second)
		require.Equal(t, finalRsv.RoomID, scsRsv.RoomID)
		require.Equal(t, finalRsv.Notes, scsRsv.Notes)
		require.Equal(t, finalRsv.Room, scsRsv.Room)

		// testify
		assert.Equal(t, http.StatusSeeOther, rr.Code)
		assert.Equal(t, "/reservation-summary", rr.Header().Get("Location"))
	})

	// Test Error: reservation missing from session
	t.Run("Missing Reservation from Session", func(t *testing.T) {
		// create a new test server, a mock database store and a request
		ts := NewTestServer(t)
		req := ts.NewRequestWithSession(t, http.MethodPost, "/make-reservation", nil)

		// build stub
		sErr := CreateServerError(ErrorMissingReservation, req.URL.Path, nil)
		ts.BuildLogErrorStub(sErr)

		//  server the request
		rr := ts.ServeRequest(req)

		// get error message from session and removes it
		errMsg := app.Session.PopString(req.Context(), "error")
		assert.Equal(t, sErr.Prompt, errMsg)

		// testify
		assert.Equal(t, http.StatusTemporaryRedirect, rr.Code)
		assert.Equal(t, "/", rr.Header().Get("Location"))

	})

	// Test Error: invalid body data cause error in ParseForm()
	t.Run("Invalid Body Data", func(t *testing.T) {
		// creating invalid body
		body := strings.NewReader("%^")

		// create a new test server, a mock database store and a request
		ts := NewTestServer(t)
		req := ts.NewRequestWithSession(t, http.MethodPost, "/make-reservation", body)

		// build stub
		ts.BuildLogAnyErrorStub()

		// put reservation in session
		app.Session.Put(req.Context(), "reservation", initRsv)

		//  server the request
		rr := ts.ServeRequest(req)

		// remove reservation from session
		app.Session.Remove(req.Context(), "reservation")

		// get error message from session and removes it
		errMsg := app.Session.PopString(req.Context(), "error")
		sErr := CreateServerError(ErrorParseForm, req.URL.Path, nil)
		assert.Equal(t, sErr.Prompt, errMsg)

		// testify
		assert.Equal(t, http.StatusTemporaryRedirect, rr.Code)
		assert.Equal(t, "/make-reservation", rr.Header().Get("Location"))

	})

	// Test Error: reservation exists in session but form is invalid
	t.Run("Invalid Form", func(t *testing.T) {
		firstName := util.RandomName()
		lastName := util.RandomName()
		email := util.RandomEmail()

		// create test cases for the form validation
		tests := []struct {
			Name   string
			Values url.Values
		}{
			{
				Name: "Missing First Name",
				Values: url.Values{
					"last_name": {lastName},
					"email":     {email},
				},
			},
			{
				Name: "Missing Last Name",
				Values: url.Values{
					"first_name": {firstName},
					"email":      {email},
				},
			},
			{
				Name: "Missing Email",
				Values: url.Values{
					"first_name": {firstName},
					"last_name":  {lastName},
				},
			},
			{
				Name: "Invalid First Name",
				Values: url.Values{
					"first_name": {"x"},
					"last_name":  {lastName},
					"email":      {email},
				},
			},
			{
				Name: "Invalid Last Name",
				Values: url.Values{
					"first_name": {firstName},
					"last_name":  {"x"},
					"email":      {email},
				},
			},
			{
				Name: "Invalid Email",
				Values: url.Values{
					"first_name": {firstName},
					"last_name":  {lastName},
					"email":      {"x"},
				},
			},
		}

		// create a new test server and a mock database store
		ts := NewTestServer(t)

		for _, test := range tests {
			t.Run(test.Name, func(t *testing.T) {
				// create the body of the request
				body := strings.NewReader(test.Values.Encode())

				// create a new request
				req := ts.NewRequestWithSession(t, http.MethodPost, "/make-reservation", body)

				// put reservation in session
				app.Session.Put(req.Context(), "reservation", initRsv)

				//  server the request
				rr := ts.ServeRequest(req)

				// remove reservation from session
				app.Session.Remove(req.Context(), "reservation")

				// testify
				assert.Equal(t, http.StatusOK, rr.Code)
			})
		}
	})

	// Test Error: reservation exists and form it valid, but internal server error on CreateReservation
	t.Run("Internal Server Error", func(t *testing.T) {
		// create the final reservation that we are expected to get from the session
		finalRsv := initRsv
		finalRsv.FirstName = util.RandomName()
		finalRsv.LastName = util.RandomName()
		finalRsv.Email = util.RandomEmail()
		finalRsv.Phone = util.RandomPhone()
		finalRsv.Notes = util.RandomNote()
		finalRsv.GenerateReservationCode()

		// create form data for the body of the request
		f := forms.New(nil)
		f.Add("first_name", finalRsv.FirstName)
		f.Add("last_name", finalRsv.LastName)
		f.Add("email", finalRsv.Email)
		f.Add("phone", finalRsv.Phone)
		f.Add("notes", finalRsv.Notes)

		// create the body of the request
		body := strings.NewReader(f.Encode())

		// create a new test server, a mock database store and a request
		ts := NewTestServer(t)
		req := ts.NewRequestWithSession(t, http.MethodPost, "/make-reservation", body)

		// create stub return arguments
		err := errors.New("this is a test error")

		sErr := ServerError{
			Prompt: "Unable to create reservation.",
			URL:    req.URL.Path,
			Err:    err,
		}

		// build stub
		ts.MockDBStore.On("CreateReservationTx", mock.Anything, mock.Anything).
			Return(db.Reservation{}, err).
			Once()
		ts.BuildLogErrorStub(sErr)

		// put reservation in session
		app.Session.Put(req.Context(), "reservation", initRsv)

		//  server the request
		rr := ts.ServeRequest(req)

		// get error message from session and remove it
		errMsg := app.Session.PopString(req.Context(), "error")
		assert.Equal(t, sErr.Prompt, errMsg)

		// testify
		assert.Equal(t, http.StatusTemporaryRedirect, rr.Code)
		assert.Equal(t, "/", rr.Header().Get("Location"))
	})
}

func TestServer_ReservationSummaryHandler(t *testing.T) {
	// Test Error: reservation data is missing from session
	t.Run("Error Missing Reservation Data", func(t *testing.T) {
		// create a new test server, and a new request
		ts := NewTestServer(t)
		req := ts.NewRequestWithSession(t, http.MethodGet, "/reservation-summary", nil)

		// build stub
		sErr := CreateServerError(ErrorMissingReservation, req.URL.Path, nil)
		ts.BuildLogErrorStub(sErr)

		//  server the request
		rr := ts.ServeRequest(req)

		// get error message from session and remove it
		errMsg := app.Session.PopString(req.Context(), "error")
		assert.Equal(t, sErr.Prompt, errMsg)

		// testify
		assert.Equal(t, http.StatusTemporaryRedirect, rr.Code)
		assert.Equal(t, "/", rr.Header().Get("Location"))
	})

	// Test OK
	t.Run("OK", func(t *testing.T) {
		//create rooms slice with random data of n rooms
		const N = 5
		rooms := randomRooms(N)

		// create random reservation
		rsv := randomReservation()

		// create a new test server, and a new request
		ts := NewTestServer(t)
		req := ts.NewRequestWithSession(t, http.MethodGet, "/reservation-summary", nil)

		// put rooms and reservation in session
		app.Session.Put(req.Context(), "rooms", rooms)
		app.Session.Put(req.Context(), "reservation", rsv)

		//  server the request
		rr := ts.ServeRequest(req)

		// checks that rooms and reservation are not in session
		ok := app.Session.Exists(req.Context(), "rooms")
		require.False(t, ok)

		ok = app.Session.Exists(req.Context(), "reservation")
		require.False(t, ok)

		// testify
		assert.Equal(t, http.StatusOK, rr.Code)
	})
}
