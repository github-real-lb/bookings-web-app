package web

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAppLogger(t *testing.T) {
	al := NewAppLogger()
	require.NotNil(t, al)
	assert.NotNil(t, al.ErrorLog)
	assert.NotNil(t, al.InfoLog)
}

func TestAppLogger_LogClientError(t *testing.T) {
	// bypass Stdout for test
	originalStdout, r, w := bypassStdout()

	// run the test
	appLogger := NewAppLogger()
	recorder := httptest.NewRecorder()
	appLogger.LogClientError(recorder, http.StatusBadRequest)

	// restore Stdout after test
	restoreStdout(originalStdout, w)

	// read the output of our prompt() function from our read pipe
	out, _ := io.ReadAll(r)

	// check results
	stdoutResult := fmt.Sprint("Client error with status ", http.StatusBadRequest)
	assert.Contains(t, string(out), stdoutResult)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)

	bodyResult := fmt.Sprint(http.StatusText(http.StatusBadRequest), "\n")
	body, _ := io.ReadAll(recorder.Body)
	assert.Equal(t, bodyResult, string(body))
}

func TestAppLogger_LogServerError(t *testing.T) {
	// bypass Stdout for test
	originalStdout, r, w := bypassStdout()

	// run the test
	appLogger := NewAppLogger()
	recorder := httptest.NewRecorder()

	sErr := "this is a test error message"
	appLogger.LogServerError(recorder, errors.New(sErr))

	// restore Stdout after test
	restoreStdout(originalStdout, w)

	// read the output of our prompt() function from our read pipe
	out, _ := io.ReadAll(r)

	// check results
	assert.Contains(t, string(out), sErr)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)

	bodyResult := fmt.Sprint(http.StatusText(http.StatusInternalServerError), "\n")
	body, _ := io.ReadAll(recorder.Body)
	assert.Equal(t, bodyResult, string(body))
}

func TestAppLogger_LogError(t *testing.T) {
	// bypass Stdout for test
	originalStdout, r, w := bypassStdout()

	// run the test
	appLogger := NewAppLogger()
	sErr := "this is a test error message"
	appLogger.LogError(errors.New(sErr))

	// restore Stdout after test
	restoreStdout(originalStdout, w)

	// read the output of our prompt() function from our read pipe
	out, _ := io.ReadAll(r)

	// check results
	assert.Contains(t, string(out), sErr)
}

// bypassStdout replace original os.Stdout with a connected pair of Files (r, w)
func bypassStdout() (originalStdout, r, w *os.File) {
	// create a read and write pipes
	r, w, _ = os.Pipe()

	// save os.Stdout and replace with writer pipe (w)
	originalStdout = os.Stdout
	os.Stdout = w
	return
}

// restoreStdout closes w and restore original os.Stdout
func restoreStdout(originalStdout, w *os.File) {
	// close writer pipe and reset os.Stdout to original state
	w.Close()
	os.Stdout = originalStdout
}
