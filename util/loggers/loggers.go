package loggers

import (
	"log"
	"os"
	"runtime/debug"

	"github.com/github-real-lb/bookings-web-app/util"
)

// AppLogger holds the app loggers
type AppLogger struct {
	ErrorChannel  chan error
	ErrorLog      *log.Logger
	InfoLog       *log.Logger
	LogDebugStack bool
}

// NewAppLogger returns and initialized AppLogger
func NewAppLogger() AppLogger {
	return AppLogger{
		ErrorLog:      log.New(os.Stdout, "ERROR\t", log.Ldate|log.LstdFlags),
		InfoLog:       log.New(os.Stdout, "INFO\t", log.Ldate|log.LstdFlags),
		LogDebugStack: false,
	}
}

// ListenAndLogErrors runs a go routine that listens to error channel and logs the errors data
// Make sure to use done channel to stop go routine and close channel
func (al *AppLogger) ListenAndLogErrors(done chan struct{}) {
	// create error channel
	al.ErrorChannel = make(chan error)

	// start go routine
	go func() {
		for {
			select {
			case err := <-al.ErrorChannel:
				al.LogError(err)
			case <-done:
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
func (al *AppLogger) Shutdown(done chan struct{}) {
	if done == nil {
		return
	}

	done <- struct{}{}
}
