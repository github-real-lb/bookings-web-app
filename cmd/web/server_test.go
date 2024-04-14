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
	// create new server, request and response recorder
	ts, _ := NewTestServer(t)
	req := ts.NewRequestWithSession(t, http.MethodGet, "/test_url", nil)
	rr := httptest.NewRecorder()

	// create an error
	err := ServerError{
		Prompt: "test prompt",
		URL:    "/test_url",
		Err:    errors.New("test error"),
	}

	// call method
	ts.LogErrorAndRedirect(rr, req, err, "/url")

	// check Status Code and redirect url
	assert.Equal(t, http.StatusTemporaryRedirect, rr.Code)
	assert.Equal(t, "/url", rr.Header().Get("Location"))

	// check session error message
	errMsg := app.Session.Pop(req.Context(), "error")
	assert.Equal(t, "test prompt", errMsg)
}

func TestServer_LogInternalServerError(t *testing.T) {
	// create new server and response recorder
	ts, _ := NewTestServer(t)
	rr := httptest.NewRecorder()

	// create an error
	err := ServerError{
		Prompt: "test prompt",
		URL:    "/test_url",
		Err:    errors.New("test error"),
	}

	// call method
	ts.LogInternalServerError(rr, err)

	expected := fmt.Sprint(http.StatusText(http.StatusInternalServerError), "\n")

	// check Status Code and redirect url
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, expected, rr.Body.String())
}

func TestServer_LogRenderErrorAndRedirect(t *testing.T) {
	// create new server, request and response recorder
	ts, _ := NewTestServer(t)
	req := ts.NewRequestWithSession(t, http.MethodGet, "/test_url", nil)
	rr := httptest.NewRecorder()

	// call method
	ts.LogRenderErrorAndRedirect(rr, req, "filename.page.gohtml", errors.New("test error"), "/url")

	// check Status Code and redirect url
	assert.Equal(t, http.StatusTemporaryRedirect, rr.Code)
	assert.Equal(t, "/url", rr.Header().Get("Location"))

	// check session error message
	errMsg := app.Session.Pop(req.Context(), "error")
	assert.Equal(t, `unable to render "filename.page.gohtml" template`, errMsg)
}

func TestServer_ResponseJSON(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		// create new server, request and response recorder
		ts, _ := NewTestServer(t)
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()

		obj := struct {
			FieldA string `json:"field_a"`
			FieldB string `json:"field_b"`
		}{
			FieldA: util.RandomString(10),
			FieldB: util.RandomString(10),
		}

		err := ts.ResponseJSON(rr, req, obj)
		require.Nil(t, err)

		jr := fmt.Sprintf(`{"field_a":"%s","field_b":"%s"}`, obj.FieldA, obj.FieldB)

		require.Equal(t, jr, rr.Body.String())
	})
}

func TestServerError_Error(t *testing.T) {
	// create an error
	err := ServerError{
		Prompt: "test prompt",
		URL:    "/test_url",
		Err:    errors.New("test error"),
	}

	assert.Equal(t, "\ttest error\n\tPROMPT: test prompt\n\tURL: /test_url\n", err.Error())
}
