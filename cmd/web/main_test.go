package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/github-real-lb/bookings-web-app/db/mocks"
	"github.com/github-real-lb/bookings-web-app/util/config"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	InitializeApp(config.TestingMode)

	os.Exit(m.Run())
}

// NewTestServer creates and returns a server connected to a mock database store
func NewTestServer(t *testing.T) (*Server, *mocks.MockStore) {
	mockStore := mocks.NewMockStore(t)
	server := NewServer(mockStore)

	return server, mockStore
}

// NewTestRequest creates a new get request for use in testing
func NewTestRequest() *http.Request {
	return httptest.NewRequest(http.MethodGet, "/", nil)
}

// NewTestRequestWithSession creates a new get request with new session data for use in testing
func NewTestRequestWithSession(t *testing.T, method string, url string, body io.Reader) *http.Request {
	// checks that the session manager is loaded
	require.NotNil(t, app.Session)

	// creating new request
	r := httptest.NewRequest(method, url, body)

	if method == http.MethodPost {
		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}

	// adding new session data to context
	ctx, err := app.Session.Load(r.Context(), "X-Session")
	require.NoError(t, err)
	require.NotNil(t, ctx)

	return r.WithContext(ctx)
}
