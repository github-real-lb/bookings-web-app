package loggers

import (
	"io"
	"log"
	"os"
	"runtime/debug"
	"sync"

	"github.com/github-real-lb/bookings-web-app/util"
)

const BufferSize = 100

// AppLogger holds the app loggers
type AppLogger struct {
	// ErrorChannel is achannel to pass errors to the error logger.
	// Use after calling ListenAndLogErrors()
	ErrorChannel chan error

	ErrorLog *log.Logger // this is the error loggers
	InfoLog  *log.Logger // this is the information loggers

	LogDebugStack bool //determines if error logger logs the debug.stack() information

	shutdown sync.Once     // ensures the shutdown is only performed once
	done     chan struct{} // used to shutdown the logger go routine
}

// NewAppLogger returns an initialized AppLogger with configurable output
func NewAppLogger(output io.Writer) AppLogger {
	if output == nil {
		output = os.Stdout
	}

	return AppLogger{
		ErrorLog:      log.New(output, "ERROR\t", log.Ldate|log.LstdFlags),
		InfoLog:       log.New(output, "INFO\t", log.Ldate|log.LstdFlags),
		LogDebugStack: false,
	}
}

// ListenAndLogErrors runs a go routine that listens to error channel and logs the errors data
// Make sure to use done channel to stop go routine and close channel
func (al *AppLogger) ListenAndLogErrors() {
	// create error channel with buffer size of 100
	al.ErrorChannel = make(chan error, BufferSize)

	// create the done channel to shutdown the go routine
	al.done = make(chan struct{})

	// start go routine to log errors
	go func() {
		for {
			select {
			case err := <-al.ErrorChannel:
				al.LogError(err)
			case <-al.done:
				for err := range al.ErrorChannel {
					al.LogError(err)
				}
				close(al.ErrorChannel)
				return
			}
		}
	}()
}

// LogError logs server side errors.
func (al *AppLogger) LogError(err error) {
	if al.LogDebugStack {
		text := util.NewText().AddLine(err).AddLine(string(debug.Stack()))
		al.ErrorLog.Println(text.String())
	} else {
		al.ErrorLog.Println(err.Error())
	}
}

// StopListen stop the go routine that listens for errors
func (al *AppLogger) Shutdown() {
	if al.done == nil {
		return
	}

	// close done channel safely
	al.shutdown.Do(func() {
		close(al.done)
	})
}
