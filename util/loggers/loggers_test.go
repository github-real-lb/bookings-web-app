package loggers

import (
	"errors"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/github-real-lb/bookings-web-app/util"
	"github.com/stretchr/testify/assert"
)

func TestNewAppLogger(t *testing.T) {
	al := NewAppLogger(os.Stdout)
	assert.Nil(t, al.ErrorChannel)
	assert.Nil(t, al.done)
	assert.NotNil(t, al.ErrorLog)
	assert.NotNil(t, al.InfoLog)

	al2 := NewAppLogger(nil)
	assert.Nil(t, al2.ErrorChannel)
	assert.Nil(t, al2.done)
	assert.NotNil(t, al2.ErrorLog)
	assert.NotNil(t, al2.InfoLog)
}

func TestAppLogger_ListenAndLogErrorsAndShutdown(t *testing.T) {
	// create a read and write pipes
	r, w, _ := os.Pipe()

	// create new AppLogger
	appLogger := NewAppLogger(w)

	// start listenning for errors
	appLogger.ListenAndLogErrors()
	defer func() {
		appLogger.Shutdown()
	}()

	// create an error
	text := util.NewText().
		AddLine("this is the first error data").
		AddLine("this is the second error data").
		AddLine("this is the second error data")

	// send error through channel
	appLogger.ErrorChannel <- text

	// wait for appLogger to log error
	time.Sleep(1 * time.Second)

	// close writer pipe
	w.Close()

	// read the output of our prompt() function from our read pipe
	out, _ := io.ReadAll(r)

	// check results
	assert.Contains(t, string(out), text.String())
}

func TestAppLogger_ListenAndLogErrorsAndShutdown_Buffer(t *testing.T) {
	//create a read and write pipes
	r, w, _ := os.Pipe()

	// create new AppLogger
	appLogger := NewAppLogger(w)

	// start listenning for errors
	appLogger.ListenAndLogErrors()

	// send N errors through channel
	for i := 0; i < 100; i++ {
		appLogger.ErrorChannel <- errors.New(fmt.Sprint(i))
	}

	// shutdown error listenning
	appLogger.Shutdown()

	// wait to log all errors in buffer
	time.Sleep(3 * time.Second)

	// close writer pipe
	w.Close()

	// read the output of our prompt() function from our read pipe
	out, _ := io.ReadAll(r)

	// testify
	result := string(out)
	assert.Len(t, result, 2890)

	result = result[2861:]
	assert.Equal(t, "ERROR", result[:5])
	assert.Equal(t, "99\n", result[26:])
}

func TestAppLogger_LogError(t *testing.T) {
	t.Run("LogDebugStack False", func(t *testing.T) {
		// create a read and write pipes
		r, w, _ := os.Pipe()

		// run the test
		appLogger := NewAppLogger(w)

		text := util.NewText().
			AddLine("this is the first error data").
			AddLine("this is the second error data").
			AddLine("this is the second error data")
		appLogger.LogError(text)

		// close writer pipe
		w.Close()

		// read the output of our prompt() function from our read pipe
		out, _ := io.ReadAll(r)

		// check results
		assert.Contains(t, string(out), text.String())
	})

	t.Run("LogDebugStack True", func(t *testing.T) {
		// create a read and write pipes
		r, w, _ := os.Pipe()

		// run the test
		appLogger := NewAppLogger(w)
		appLogger.LogDebugStack = true

		text := util.NewText().
			AddLine("this is the first error data").
			AddLine("this is the second error data").
			AddLine("this is the second error data")
		appLogger.LogError(text)

		// close writer pipe
		w.Close()

		// read the output of our prompt() function from our read pipe
		out, _ := io.ReadAll(r)

		// check results
		assert.Contains(t, string(out), text.String())
	})
}
