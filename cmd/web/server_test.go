package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/github-real-lb/bookings-web-app/db/mocks"
	"github.com/github-real-lb/bookings-web-app/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewServer(t *testing.T) {
	store := mocks.NewMockStore(t)

	server := NewServer(store)
	require.IsType(t, (*Server)(nil), server)
}

func TestServer_LogErrorAndRedirect(t *testing.T) {
	ts, _ := NewTestServer(t)
	req := ts.NewRequestWithSession(t, http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	s := NewServer(nil)
	s.LogErrorAndRedirect(rr, req, "test message", errors.New("test error"), "/test_url")

	// check Status Code and redirect url
	assert.Equal(t, http.StatusTemporaryRedirect, rr.Code)
	assert.Equal(t, "/test_url", rr.Header().Get("Location"))

	// check session error message
	errMsg := app.Session.Pop(req.Context(), "error")
	assert.Equal(t, "test message", errMsg)
}

func TestServer_ResponseJSON(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		s := NewServer(nil)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()

		obj := struct {
			FieldA string `json:"field_a"`
			FieldB string `json:"field_b"`
		}{
			FieldA: util.RandomString(10),
			FieldB: util.RandomString(10),
		}

		err := s.ResponseJSON(rr, req, obj)
		require.NoError(t, err)

		jr := fmt.Sprintf(`{"field_a":"%s","field_b":"%s"}`, obj.FieldA, obj.FieldB)

		require.Equal(t, jr, rr.Body.String())
	})
}
