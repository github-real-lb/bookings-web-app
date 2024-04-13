package loggers

import (
	"io"
	"os"
	"testing"

	"github.com/github-real-lb/bookings-web-app/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAppLogger(t *testing.T) {
	al := NewAppLogger()
	require.NotNil(t, al)
	assert.NotNil(t, al.ErrorChannel)
	assert.NotNil(t, al.ErrorLog)
	assert.NotNil(t, al.InfoLog)
}

func TestAppLogger_ListenAndLogErrorsAndShutdown(t *testing.T) {
	// bypass Stdout for test
	originalStdout, r, w := bypassStdout()

	// start ListenAndLogErrors
	appLogger := NewAppLogger()
	appLogger.ListenAndLogErrors()

	// send error through channel
	text := util.NewText().
		AddLine("this is the first error data").
		AddLine("this is the second error data").
		AddLine("this is the second error data")
	appLogger.ErrorChannel <- text

	// shutdown ListenAndLogErrors
	appLogger.Shutdown()

	// restore Stdout after test
	restoreStdout(originalStdout, w)

	// read the output of our prompt() function from our read pipe
	out, _ := io.ReadAll(r)

	// check results
	assert.Contains(t, string(out), text.String())
}

func TestAppLogger_LogError(t *testing.T) {
	t.Run("LogDebugStack False", func(t *testing.T) {
		// bypass Stdout for test
		originalStdout, r, w := bypassStdout()

		// run the test
		appLogger := NewAppLogger()

		text := util.NewText().
			AddLine("this is the first error data").
			AddLine("this is the second error data").
			AddLine("this is the second error data")
		appLogger.LogError(text)

		// restore Stdout after test
		restoreStdout(originalStdout, w)

		// read the output of our prompt() function from our read pipe
		out, _ := io.ReadAll(r)

		// check results
		assert.Contains(t, string(out), text.String())
	})

	t.Run("LogDebugStack True", func(t *testing.T) {
		// bypass Stdout for test
		originalStdout, r, w := bypassStdout()

		// run the test
		appLogger := NewAppLogger()
		appLogger.LogDebugStack = true

		text := util.NewText().
			AddLine("this is the first error data").
			AddLine("this is the second error data").
			AddLine("this is the second error data")
		appLogger.LogError(text)

		// restore Stdout after test
		restoreStdout(originalStdout, w)

		// read the output of our prompt() function from our read pipe
		out, _ := io.ReadAll(r)

		// check results
		assert.Contains(t, string(out), text.String())
	})
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
