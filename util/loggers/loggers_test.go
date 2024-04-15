package loggers

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/github-real-lb/bookings-web-app/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLogger(t *testing.T) {
	t.Run("Test With os.Pipe", func(t *testing.T) {
		var buf bytes.Buffer
		sl := NewSmartLogger(&buf, "TEST")

		// testify
		assert.Nil(t, sl.LogChannel)
		require.NotNil(t, sl.Logger)

		// log a message
		sl.Logger.Print("this is a test message")

		// testify
		assert.Contains(t, buf.String(), "TEST")
		assert.Contains(t, buf.String(), "this is a test message")
	})

	t.Run("Test With os.Stdout", func(t *testing.T) {
		//create a read and write pipes
		r, w, _ := os.Pipe()

		//replace Stdout with writer pipe
		originalStdout := os.Stdout
		os.Stdout = w

		// create new smart logger with default output (Stdout) and INFO prefix
		sl := NewSmartLogger(nil, "TEST")

		// testify
		assert.Nil(t, sl.LogChannel)
		require.NotNil(t, sl.Logger)

		// log a message
		sl.Logger.Print("this is a test message")

		// close writer pipe and restore Stdout
		w.Close()
		os.Stdout = originalStdout

		// read the output of our prompt() function from our read pipe
		out, _ := io.ReadAll(r)
		result := string(out)

		// testify
		assert.Contains(t, result, "TEST")
		assert.Contains(t, result, "this is a test message")
	})
}

func TestAppLogger_Log(t *testing.T) {
	t.Run("LogDebugStack False", func(t *testing.T) {
		var buf bytes.Buffer
		sl := NewSmartLogger(&buf, "TEST")

		text := util.NewText().
			AddLine("this is the first error data").
			AddLine("this is the second error data").
			AddLine("this is the second error data")
		sl.Log(text)

		// check results
		assert.Contains(t, buf.String(), text.String())
	})

	t.Run("LogDebugStack True", func(t *testing.T) {
		var buf bytes.Buffer
		sl := NewSmartLogger(&buf, "TEST")

		sl.LogDebugStack = true

		text := util.NewText().
			AddLine("this is the first error data").
			AddLine("this is the second error data").
			AddLine("this is the second error data")
		sl.Log(text)

		// check results
		assert.Contains(t, buf.String(), text.String())
	})
}

func TestAppLogger_ListenAndLogAndShutdown(t *testing.T) {
	var buf bytes.Buffer
	sl := NewSmartLogger(&buf, "TEST")

	// Start ListenAndLog in a goroutine with a buffer size of 100
	go sl.ListenAndLog(100)

	// wait to ensure ListenAndLog has started
	time.Sleep(100 * time.Millisecond)

	// Send a message to the logger
	sl.LogChannel <- "Test message"

	// Allow some time for the message to be processed
	time.Sleep(1 * time.Second)

	// Shutdown the logger
	sl.Shutdown()

	// Ensure all messages are processed before channel is closed
	assert.Contains(t, buf.String(), "Test message")

	// Make sure the channel is closed
	_, ok := <-sl.LogChannel
	assert.Falsef(t, ok, "LogChannel should be closed after Shutdown")
}

func TestAppLogger_ListenAndLogAndShutdown_Buffer(t *testing.T) {
	var buf bytes.Buffer
	sl := NewSmartLogger(&buf, "TEST")

	// start listening for errors on a separate go routine
	go sl.ListenAndLog(100)

	// wait to ensure ListenAndLog has started
	time.Sleep(100 * time.Millisecond)

	// send N errors through channel
	for i := 0; i < 100; i++ {
		sl.LogChannel <- errors.New(fmt.Sprint(i))
	}

	// wait for smart logger to log all errors in buffer
	time.Sleep(3 * time.Second)

	// shutdown listening
	sl.Shutdown()

	// testify
	result := buf.String()
	require.Len(t, result, 2690)

	result = result[2663:]
	assert.Equal(t, "TEST", result[:4])
	assert.Equal(t, "99\n", result[24:])

	// Make sure the channel is closed
	_, ok := <-sl.LogChannel
	assert.Falsef(t, ok, "LogChannel should be closed after Shutdown")
}
