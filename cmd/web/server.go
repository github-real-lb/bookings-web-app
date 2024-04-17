package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/github-real-lb/bookings-web-app/db"
	"github.com/github-real-lb/bookings-web-app/util/loggers"
	"github.com/github-real-lb/bookings-web-app/util/mailers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	LoggerBufferSize = 100
	MailerBufferSize = 100
)

// Server handles all routing and provides all database functions
type Server struct {
	Router        *http.Server
	DatabaseStore db.DatabaseStore
	ErrorLogger   loggers.Loggerer
	InfoLogger    loggers.Loggerer
	Mailer        mailers.Mailerer
}

// NewServer returns a new Server with Router and Database Store
func NewServer(store db.DatabaseStore, errLogger loggers.Loggerer, infoLogger loggers.Loggerer, mailer mailers.Mailerer) *Server {
	mux := chi.NewRouter()

	s := Server{
		Router: &http.Server{
			Addr:    app.ServerAddress,
			Handler: mux,
		},
		DatabaseStore: store,
		ErrorLogger:   errLogger,
		InfoLogger:    infoLogger,
		Mailer:        mailer,
	}

	// add middleware that recover from panics
	mux.Use(middleware.Recoverer)

	// add middleware that loads and saves and session on every request
	mux.Use(app.Session.LoadAndSave)

	// add middleware that provides CSRF protection to all POST requests
	if !app.InTestingMode() {
		mux.Use(NoSurf)
	}

	// add middleware that logs incoming requests and their responses
	mux.Use(s.LogRequestsAndResponse)

	// setting routes
	mux.Get("/", s.HomeHandler)
	mux.Get("/about", s.AboutHandler)

	mux.Get("/rooms/{index}", s.RoomsHandler)
	mux.Get("/rooms/room/{name}", s.RoomHandler)
	mux.Post("/search-room-availability", s.PostSearchRoomAvailabilityHandler)

	mux.Get("/contact", s.ContactHandler)

	mux.Get("/available-rooms-search", s.AvailableRoomsSearchHandler)
	mux.Post("/available-rooms-search", s.PostAvailableRoomsSearchHandler)
	mux.Get("/available-rooms/{index}", s.AvailableRoomsListHandler)

	mux.Get("/make-reservation", s.MakeReservationHandler)
	mux.Post("/make-reservation", s.PostMakeReservationHandler)

	mux.Get("/reservation-summary", s.ReservationSummaryHandler)

	// setting file server
	fileServer := http.FileServer(http.Dir(app.StaticPath))
	mux.Handle("/"+app.StaticDirectoryName+"/*", http.StripPrefix("/"+app.StaticDirectoryName, fileServer))

	return &s
}

// Start calls the http.Server ListenAndServer method
func (s *Server) Start() {
	fmt.Printf("Starting http server on %s ... \n", app.ServerAddress)

	// start listening to errors
	go s.ErrorLogger.ListenAndLog(LoggerBufferSize)

	// start listening to info
	go s.InfoLogger.ListenAndLog(LoggerBufferSize)

	// start listening to mail data
	go s.Mailer.ListenAndMail(s.ErrorLogger.MyLogChannel(), MailerBufferSize)

	// start listening to http requests
	err := s.Router.ListenAndServe()

	if err != nil && err.Error() != "http: Server closed" {
		s.ErrorLogger.Log(ServerError{
			Prompt: "error starting http server",
			Err:    err,
		})
	}
}

// Stop calls the http.Server Shutdown method
func (s *Server) Stop() {
	// create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Print("Shutting down http server... ")

	// inform the server to stop accepting new info
	s.InfoLogger.Shutdown()

	// inform the server to stop accepting new errors
	s.ErrorLogger.Shutdown()

	// inform the server to stop accepting new mail data
	s.Mailer.Shutdown()

	// inform the server to stop accepting new requests
	err := s.Router.Shutdown(ctx)

	// wait for existing connections to finish processing before returning from this function
	<-ctx.Done()

	if err != nil {
		fmt.Println("Error")

		s.ErrorLogger.Log(ServerError{
			Prompt: "error shuting down http server",
			Err:    err,
		})
	} else {
		fmt.Println("Success")
	}
}

// LogError logs err using the ErrorLogger
func (s *Server) LogInfo(info string) {
	var infoChan = s.InfoLogger.MyLogChannel()

	// if log channel is nil logging directly with logger
	if infoChan == nil {
		s.InfoLogger.Log(info)
		return
	}

	// logging through channel
	infoChan <- info
}

// LogError logs err using the ErrorLogger
func (s *Server) LogError(err error) {
	var errChan = s.ErrorLogger.MyLogChannel()

	// if log channel is nil logging directly with logger
	if errChan == nil {
		s.ErrorLogger.Log(err)
		return
	}

	// logging through channel
	errChan <- err
}

// LogErrorAndRedirect logs err, puts err's prompt in session, and redirects to url
func (s *Server) LogErrorAndRedirect(w http.ResponseWriter, r *http.Request, err ServerError, url string) {
	app.Session.Put(r.Context(), "error", err.Prompt)
	s.LogError(err)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// LogInternalServerError logs err, and sends StatusInternalServerError response
func (s *Server) LogInternalServerError(w http.ResponseWriter, err error) {
	s.LogError(err)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// ResponseJSON write v to w as json response.
// Errors are loggied by the server and also returned
func (s *Server) ResponseJSON(w http.ResponseWriter, r *http.Request, v any) error {
	// set Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// marshal v to json
	bs, err := json.Marshal(v)
	if err != nil {
		e := ServerError{
			Prompt: "unable to marshal json response",
			Err:    err,
		}
		s.LogError(e)
		return e
	}

	// write json data to the response body
	_, err = w.Write(bs)
	if err != nil {
		e := ServerError{
			Prompt: "unable to write json response",
			Err:    err,
		}

		s.LogError(e)
		return e
	}

	return nil
}

// SendMail sends email using the Mailer
func (s *Server) SendMail(data mailers.MailData) {
	var err error
	var mailChan = s.Mailer.MyMailChannel()

	// if mail channel is nil sending email directly with Mailer
	if mailChan == nil {
		err = s.Mailer.SendMail(data)
		if err != nil {
			s.LogError(err)
		}
		return
	}
	// sending email through the channel
	mailChan <- data
}
