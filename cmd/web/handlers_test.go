package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testCase is used as a single Test Case for specific HandlerFunc
type testCase struct {
	name          string // name of test
	method        string // http.Method for the http.Request
	url           string // url for the http.Request
	body          any    // the json body for the http.Request
	checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
}

// sendRequestToTestServer start test server and send the test request
func (tc *testCase) sendRequestToServer(t *testing.T) *httptest.ResponseRecorder {
	// start test server and send request
	server := NewServer(ADDRESS)
	recorder := httptest.NewRecorder()

	var reader io.Reader = nil

	// creating new reader with arguments passed
	if tc.body != nil {
		jsonData, err := json.Marshal(tc.body)
		require.NoError(t, err)

		reader = strings.NewReader(string(jsonData))
	}

	request, err := http.NewRequest(tc.method, tc.url, reader)
	require.NoError(t, err)

	server.Handler.ServeHTTP(recorder, request)

	return recorder
}

func TestHome(t *testing.T) {
	tc := testCase{
		name:   "OK",
		method: http.MethodGet,
		url:    "/",
		body:   nil,
		checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, recorder.Code)
		},
	}

	templatePath = "./../../templates"
	InitApp()

	recorder := tc.sendRequestToServer(t)
	tc.checkResponse(t, recorder)
}
