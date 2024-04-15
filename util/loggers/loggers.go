package loggers

import (
	"io"
	"log"
	"os"
	"runtime/debug"
	"sync"

	"github.com/github-real-lb/bookings-web-app/util"
)

type Loggerer interface {
	MyLogChannel() chan any
	IsLogDebugStack() bool
	Log(v any)
	ListenAndLog(buffer int)
	Shutdown()
}

// SmartLogger is a configurable logger
type SmartLogger struct {
	*log.Logger // this is the loggers

	// ErrorChannel is achannel to pass errors to the error logger.
	// Use after calling ListenAndLogErrors()
	LogChannel chan any

	LogDebugStack bool //determines if error logger logs the debug.stack() information

	shutdown sync.Once // ensures Shutdown() is only performed once
}

// NewSmartLogger returns an initialized SmartLogger with configurable output and a prefix string to log.
// Example: NewLogger(os.Stdout,"ERROR")
// Logger Output: ERROR 2009/01/23 01:23:23	message
func NewSmartLogger(output io.Writer, prefix string) *SmartLogger {
	if output == nil {
		output = os.Stdout
	}

	return &SmartLogger{
		Logger:        log.New(output, prefix, log.Ldate|log.LstdFlags),
		LogDebugStack: false,
	}
}

func (sl *SmartLogger) MyLogChannel() chan any {
	return sl.LogChannel
}

func (sl *SmartLogger) IsLogDebugStack() bool {
	return sl.LogDebugStack
}

// LogError logs server side errors.
func (sl *SmartLogger) Log(v any) {
	if sl.LogDebugStack {
		text := util.NewText().AddLine(v).AddLine(string(debug.Stack()))
		sl.Logger.Println(text.String())
	} else {
		sl.Logger.Println(v)
	}
}

// ListenAndLog listens to LogChannel and logs received data
// buffer determine the buffer size of the channel. buffer = 100 is the minimum
// Make sure to use Shutdown() to stop listening and close channel
func (sl *SmartLogger) ListenAndLog(buffer int) {
	if buffer < 100 {
		buffer = 100
	}
	// create error channel with buffer size of 100
	sl.LogChannel = make(chan any, buffer)

	// start listening
	for {
		v, ok := <-sl.LogChannel
		if !ok {
			return
		} else {
			// logging data
			sl.Log(v)
		}
	}
}

// Shutdown stops ListenAndLog() and close channels
func (sl *SmartLogger) Shutdown() {
	if sl.LogChannel == nil {
		return
	}

	// close done channel safely
	sl.shutdown.Do(func() {
		close(sl.LogChannel)
	})
}
