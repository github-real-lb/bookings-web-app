package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	dbmocks "github.com/github-real-lb/bookings-web-app/db/mocks"
	"github.com/github-real-lb/bookings-web-app/util"
	loggermocks "github.com/github-real-lb/bookings-web-app/util/loggers/mocks"
	"github.com/github-real-lb/bookings-web-app/util/mailers"
	mailermocks "github.com/github-real-lb/bookings-web-app/util/mailers/mocks"
	"github.com/stretchr/testify/assert"
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
	ts := NewTestServer(t)

	// build stubs for Start()
	logChan := make(chan any, LoggerBufferSize)
	defer close(logChan)

	ts.MockLogger.On("ListenAndLog", LoggerBufferSize).Times(1)
	ts.MockLogger.On("MyLogChannel").Return(logChan).Times(1)
	ts.MockMailer.On("ListenAndMail", logChan, MailerBufferSize).Times(1)

	// Use a goroutine to handle Start because it is blocking
	go ts.Start()

	// wait a bit for the server to start
	time.Sleep(500 * time.Millisecond)

	// make an HTTP request to ensure the server is responding
	fmt.Print("Sending an HTTP request... ")
	resp, err := http.Get("http://" + app.ServerAddress)
	require.NoError(t, err)
	defer resp.Body.Close()

	fmt.Println("Sucess")

	// Check if the response is what you expect
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// build stubs for Stop()
	ts.MockLogger.On("Shutdown").Times(1)
	ts.MockMailer.On("Shutdown").Times(1)

	// stop the server
	ts.Stop()

	// Try to make a request after stopping
	_, err = http.Get("http://" + app.ServerAddress)
	require.Error(t, err) // Expect an error because the server should be closed
}

func TestServer_LogError(t *testing.T) {
	t.Run("LogChannel nil", func(t *testing.T) {
		// create new test server
		ts := NewTestServer(t)

		// create test error
		err := errors.New("this is a test error")

		// build stub
		ts.BuildLogErrorStub(err)

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
		ts.MockLogger.On("MyLogChannel").Return(logChan).Times(1)

		// create test error
		err := errors.New("this is a test error")

		// use a goroutine to handle LogError because it is blocking
		go ts.LogError(err)

		// wait a bit for the LogError to executet
		time.Sleep(500 * time.Millisecond)

		// get err back from log channel
		result := <-logChan

		assert.Equal(t, err.Error(), result.(error).Error())
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

	// build stub
	ts.BuildLogErrorStub(err)

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

	// build stub
	ts.BuildLogErrorStub(err)

	// call method
	ts.LogInternalServerError(rr, err)

	expected := fmt.Sprint(http.StatusText(http.StatusInternalServerError), "\n")

	// check Status Code and redirect url
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, expected, rr.Body.String())
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

func TestServer_SendMail(t *testing.T) {
	t.Run("MailChannel nil, Mail OK", func(t *testing.T) {
		// create new test server
		ts := NewTestServer(t)

		// create test mail data
		data := mailers.MailData{}

		// build stub
		ts.BuildSendMailStub(data)

		// log error
		ts.SendMail(data)
	})

	t.Run("MailChannel nil, Mail Error", func(t *testing.T) {
		// create new test server
		ts := NewTestServer(t)

		// create test mail data
		data := mailers.MailData{}

		err := errors.New("Test error")

		// build stubs
		ts.MockMailer.On("MyMailChannel").Return(nil).Times(1)
		ts.MockMailer.On("SendMail", data).Return(err).Times(1)
		ts.BuildLogErrorStub(err)

		// log error
		ts.SendMail(data)
	})

	t.Run("MailChannel active", func(t *testing.T) {
		// create new test server
		ts := NewTestServer(t)

		// create mail channel
		mailChan := make(chan mailers.MailData)
		defer close(mailChan)

		// build stubs
		ts.MockMailer.On("MyMailChannel").Return(mailChan).Times(1)

		// create test mail data
		data := mailers.MailData{
			To:      "test to",
			From:    "test from",
			Subject: "test subject",
			Content: "test content",
		}

		// use a goroutine to handle SendMail because it is blocking
		go ts.SendMail(data)

		// wait a bit for the SendMail to execute
		time.Sleep(500 * time.Millisecond)

		// get data back from mail channel
		result := <-mailChan

		assert.Equal(t, data, result)
	})
}
