package loggers

import (
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
)

// AppLogger holds the app loggers
type AppLogger struct {
	ErrorChannel  ErrorChannel
	ErrorLog      *log.Logger
	InfoLog       *log.Logger
	LogDebugStack bool
}

// ErrorData holds error data
type ErrorData struct {
	Prefix string
	Error  error
}

// ErrorChannel is a channel to pass errors data
type ErrorChannel chan ErrorData

// NewAppLogger returns and initialized AppLogger
func NewAppLogger() AppLogger {
	return AppLogger{
		ErrorChannel:  make(ErrorChannel),
		ErrorLog:      log.New(os.Stdout, "ERROR\t", log.Ldate|log.LstdFlags|log.Lshortfile),
		InfoLog:       log.New(os.Stdout, "INFO\t", log.Ldate|log.LstdFlags),
		LogDebugStack: false,
	}
}

// ListenForErrors runs a go routine that listens to error channel and logs the errors data
func (al *AppLogger) ListenForErrors() {
	if al.ErrorChannel == nil {
		return
	}

	go func() {
		for {
			e := <-al.ErrorChannel
			if e.Prefix == "$$$SHUTDOWN$$$" {
				return
			}

			al.LogError(e)
		}
	}()
}

// StopListen stop the go routine that listens for errors
func (al *AppLogger) StopListen() {
	if al.ErrorChannel == nil {
		return
	}

	al.ErrorChannel <- ErrorData{
		Prefix: "$$$SHUTDOWN$$$",
	}
}

// LogServerError handles server side error of the web app, loggin the error and
// writing StatusInternalServerError to the http.ResponseWriter
func (al *AppLogger) LogServerError(w http.ResponseWriter, e ErrorData) {
	al.LogError(e)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// LogError logs server side errors.
func (al *AppLogger) LogError(e ErrorData) {
	if al.LogDebugStack {
		al.ErrorLog.Println(e.String(), "\n", string(debug.Stack()))
	} else {
		al.ErrorLog.Println(e.String())
	}

}

func (e *ErrorData) String() string {
	var sb strings.Builder

	if e.Prefix != "" {
		sb.WriteString(e.Prefix)
		sb.WriteString("\n")
	}

	if e.Error != nil {
		sb.WriteString("\t")
		sb.WriteString(e.Error.Error())
		sb.WriteString("\n")
	}
	return sb.String()
}
