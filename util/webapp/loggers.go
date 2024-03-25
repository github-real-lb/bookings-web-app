package webapp

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/debug"
)

type AppLogger struct {
	ErrorLog *log.Logger
	InfoLog  *log.Logger
}

// NewAppLogger returns and initialized AppLogger
func NewAppLogger() AppLogger {
	return AppLogger{
		ErrorLog: log.New(os.Stdout, "ERROR\t", log.Ldate|log.LstdFlags|log.Lshortfile),
		InfoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.LstdFlags),
	}
}

// LogClientError handles client side error of the web app, loggin the error and
// writing the StatusCode to the http.ResponseWriter
func (al *AppLogger) LogClientError(w http.ResponseWriter, code int) {
	al.InfoLog.Println("Client error with status", code)
	http.Error(w, http.StatusText(code), code)
}

// LogServerError handles server side error of the web app, loggin the error and
// writing StatusInternalServerError to the http.ResponseWriter
func (al *AppLogger) LogServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	al.ErrorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// LogError logs server side errors
func (al *AppLogger) LogError(err error) {
	s := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	al.ErrorLog.Println(s)
}
