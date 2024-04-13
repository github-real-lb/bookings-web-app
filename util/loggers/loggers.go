package loggers

import (
	"errors"
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
		ErrorChannel:  make(chan error),
		ErrorLog:      log.New(os.Stdout, "ERROR\t", log.Ldate|log.LstdFlags),
		InfoLog:       log.New(os.Stdout, "INFO\t", log.Ldate|log.LstdFlags),
		LogDebugStack: false,
	}
}

// ListenAndLogErrors runs a go routine that listens to error channel and logs the errors data
func (al *AppLogger) ListenAndLogErrors() {
	if al.ErrorChannel == nil {
		return
	}

	go func() {
		for {
			err := <-al.ErrorChannel
			if err.Error() == "$$$SHUTDOWN$$$" {
				return
			}

			al.LogError(err)
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
	if al.ErrorChannel == nil {
		return
	}

	al.ErrorChannel <- errors.New("$$$SHUTDOWN$$$")
}
