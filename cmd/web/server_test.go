package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	dbmocks "github.com/github-real-lb/bookings-web-app/db/mocks"
	"github.com/github-real-lb/bookings-web-app/util"
	loggermocks "github.com/github-real-lb/bookings-web-app/util/loggers/mocks"
	mailermocks "github.com/github-real-lb/bookings-web-app/util/mailers/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewServer(t *testing.T) {
	mockDbStore := dbmocks.NewMockDBStore(t)
	mockLogger := loggermocks.NewMockLogger(t)
	mockMailer := mailermocks.NewMockMailer(t)

	server := NewServer(mockDbStore, mockLogger, mockMailer)
	require.IsType(t, (*Server)(nil), server)
}

func TestServer_StartAndStop(t *testing.T) {
	//TODO
}

func TestServer_LogError(t *testing.T) {
	t.Run("LogChannel nil", func(t *testing.T) {
		// create new test server
		ts := NewTestServer(t)

		// create test error
		err := errors.New("this is a test error")

		// build stubs
		ts.MockLogger.On("MyLogChannel").Return(nil).Times(1)
		ts.MockLogger.On("Log", err).Times(1)

		// log error
		ts.LogError(err)
	})

	t.Run("LogChannel active", func(t *testing.T) {
		// create new test server
		ts := NewTestServer(t)

		// create log channel
		logChan := make(chan any)
		defer close(logChan)

		// build stubs
		ts.MockLogger.On("MyLogChannel").Return(logChan).Times(2)
		ts.MockLogger.On("Log", mock.Anything).Times(0)

		// create test error
		err := errors.New("this is a test error")

		// log error
		go ts.LogError(err)

		// get err back from log channel
		result := <-logChan

		assert.Equal(t, err.Error(), result.(error).Error())

		// remove unused stubs
		ts.MockLogger.On("Log", mock.Anything).Unset()
	})
}

func TestServer_LogErrorAndRedirect(t *testing.T) {
	// create new server, request and response recorder
	ts := NewTestServer(t)
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
	ts := NewTestServer(t)
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
	ts := NewTestServer(t)
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
		ts := NewTestServer(t)
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
