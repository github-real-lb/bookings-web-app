package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

// NewTestRequest creates a new get request for use in testing
func NewTestRequest() *http.Request {
	return httptest.NewRequest(http.MethodGet, "/", nil)
}

// NewTestRequestWithSession creates a new get request with new session data for use in testing
func NewTestRequestWithSession(t *testing.T) *http.Request {
	// checks that the session manager is loaded
	require.NotNil(t, app.Session)

	// creating new request
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	// adding new session data to context
	ctx, err := app.Session.Load(r.Context(), "X-Session")
	require.NoError(t, err)
	require.NotNil(t, ctx)

	return r.WithContext(ctx)
}
