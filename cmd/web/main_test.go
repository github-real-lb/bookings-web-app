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

	// start listenning for errors
	app.Logger.ListenAndLogErrors()
	defer func() {
		app.Logger.Shutdown()
	}()

	// run tests
	code := m.Run()

	os.Exit(code)
}

type TestServer struct {
	*Server
}

// NewTestServer creates and returns a test server connected to a mock database store
func NewTestServer(t *testing.T) (*TestServer, *mocks.MockStore) {
	mockStore := mocks.NewMockStore(t)
	server := NewServer(mockStore)
	return &TestServer{Server: server}, mockStore
}

// NewTestRequest creates a new get request for use in testing
func (ts *TestServer) NewRequest(method string, url string, body io.Reader) *http.Request {
	return httptest.NewRequest(method, url, body)
}

// NewTestRequestWithSession creates a new get request with new session data for use in testing
func (ts *TestServer) NewRequestWithSession(t *testing.T, method string, url string, body io.Reader) *http.Request {
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

// ServeRequest execute a ServerHTTP method and return the response recorder
func (ts *TestServer) ServeRequest(r *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	ts.Router.Handler.ServeHTTP(rr, r)

	return rr
}
