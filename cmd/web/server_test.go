package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/github-real-lb/bookings-web-app/db"
	"github.com/github-real-lb/bookings-web-app/db/mocks"
	"github.com/github-real-lb/bookings-web-app/util/webapp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewServer(t *testing.T) {
	store := mocks.NewMockStore(t)

	server := NewServer(store)
	require.IsType(t, (*Server)(nil), server)
}

func TestPageHandlers(t *testing.T) {
	type params map[string]string

	tests := []struct {
		name   string // name of test
		method string // http.Method for the http.Request
		url    string // url for the http.Request
		params
		excpectedStatusCode int
	}{
		{"Home", http.MethodGet, "/", nil, http.StatusOK},
		{"About", http.MethodGet, "/about", nil, http.StatusOK},
		{"Generals", http.MethodGet, "/generals-quarters", nil, http.StatusOK},
		{"Majors", http.MethodGet, "/majors-suite", nil, http.StatusOK},
		{"Contact", http.MethodGet, "/contact", nil, http.StatusOK},
		{"Availability", http.MethodGet, "/search-availability", nil, http.StatusOK},
		{"Reservation", http.MethodGet, "/make-reservation", nil, http.StatusOK},
		{"PostAvailability", http.MethodPost, "/search-availability",
			params{
				"start_date": "2024-05-01",
				"end_date":   "2024-05-08",
			}, http.StatusOK},
		{"PostAvailabilityJSON", http.MethodPost, "/search-availability-json",
			params{
				"start_date": "2024-05-01",
				"end_date":   "2024-05-08",
			}, http.StatusOK},
		{"PostReservation_OK", http.MethodPost, "/make-reservation",
			params{
				"first_name": "John",
				"last_name":  "Dow",
				"email":      "john.dow@gmail.com",
				"phone":      "5555-5555",
			}, http.StatusOK},
	}

	// initialize the Application Config and Templates
	InitializeApp(webapp.TestingMode)

	// start test server and send request
	store, err := db.NewPostgresDBStore(app.ConnectionString)
	require.NoError(t, err)

	server := NewServer(store)
	testServer := httptest.NewTLSServer(server.Router.Handler)
	defer testServer.Close()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.method == http.MethodGet {
				response, err := testServer.Client().Get(testServer.URL + test.url)
				require.NoError(t, err)

				assert.Equal(t, test.excpectedStatusCode, response.StatusCode)
			} else {
				data := url.Values{}

				for key, value := range test.params {
					data.Set(key, value)
				}

				response, err := testServer.Client().PostForm(testServer.URL+test.url, data)
				require.NoError(t, err)
				defer response.Body.Close()

				assert.Equal(t, test.excpectedStatusCode, response.StatusCode)
			}
		})
	}
}
