package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewServer(t *testing.T) {
	server := NewServer(ADDRESS)
	assert.IsType(t, (*Server)(nil), server)
}

func TestNewHandler(t *testing.T) {
	handler := NewHandler(Store{})
	assert.Implements(t, (*http.Handler)(nil), handler)
}

func TestPageHandlers(t *testing.T) {
	tests := []struct {
		name   string // name of test
		method string // http.Method for the http.Request
		url    string // url for the http.Request
		params []struct {
			key   string
			value string
		}
		excpectedStatusCode int
	}{
		{"Home", http.MethodGet, "/", nil, http.StatusOK},
		{"About", http.MethodGet, "/about", nil, http.StatusOK},
		{"Generals", http.MethodGet, "/generals-quarters", nil, http.StatusOK},
		{"Majors", http.MethodGet, "/majors-suite", nil, http.StatusOK},
		{"Contact", http.MethodGet, "/contact", nil, http.StatusOK},
		{"Availability", http.MethodGet, "/search-availability", nil, http.StatusOK},
		{"Reservation", http.MethodGet, "/make-reservation", nil, http.StatusOK},
	}

	// initialize the Application Config and Templates
	templatePath = "./../../templates"
	InitApp()

	// start test server and send request
	testServer := httptest.NewTLSServer(NewHandler(Store{}))
	defer testServer.Close()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.method == http.MethodGet {
				response, err := testServer.Client().Get(testServer.URL + test.url)
				require.NoError(t, err)

				assert.Equal(t, test.excpectedStatusCode, response.StatusCode)
			} else {
				// TODO
			}

		})
	}
}
