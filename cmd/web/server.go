package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/github-real-lb/bookings-web-app/db"
	"github.com/github-real-lb/bookings-web-app/util"
	"github.com/github-real-lb/bookings-web-app/util/loggers"
	"github.com/github-real-lb/bookings-web-app/util/mailers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Server handles all routing and provides all database functions
type Server struct {
	Router        *http.Server
	DatabaseStore db.DatabaseStore
	ErrorLogger   loggers.Loggerer
	Mailer        mailers.Mailerer
}

// NewServer returns a new Server with Router and Database Store
func NewServer(store db.DatabaseStore, logger loggers.Loggerer, mailer mailers.Mailerer) *Server {
	mux := chi.NewRouter()

	server := Server{
		Router: &http.Server{
			Addr:    app.ServerAddress,
			Handler: mux,
		},
		DatabaseStore: store,
		ErrorLogger:   logger,
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

	// setting routes
	mux.Get("/", server.HomeHandler)
	mux.Get("/about", server.AboutHandler)

	mux.Get("/rooms/{index}", server.RoomsHandler)
	mux.Get("/rooms/room/{name}", server.RoomHandler)
	mux.Post("/search-room-availability", server.PostSearchRoomAvailabilityHandler)

	mux.Get("/contact", server.ContactHandler)

	mux.Get("/available-rooms-search", server.AvailableRoomsSearchHandler)
	mux.Post("/available-rooms-search", server.PostAvailableRoomsSearchHandler)
	mux.Get("/available-rooms/{index}", server.AvailableRoomsListHandler)

	mux.Get("/make-reservation", server.MakeReservationHandler)
	mux.Post("/make-reservation", server.PostMakeReservationHandler)

	mux.Get("/reservation-summary", server.ReservationSummaryHandler)

	// setting file server
	fileServer := http.FileServer(http.Dir(app.StaticPath))
	mux.Handle("/"+app.StaticDirectoryName+"/*", http.StripPrefix("/"+app.StaticDirectoryName, fileServer))

	return &server
}

// Start calls the http.Server ListenAndServer method
func (s *Server) Start() {
	fmt.Printf("Starting http server on %s ... \n", app.ServerAddress)

	// start listening to errors
	go s.ErrorLogger.ListenAndLog(100)

	// start listening to mail data
	go s.Mailer.ListenAndMail(s.ErrorLogger.MyLogChannel(), 100)

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

// LogError logs err
func (s *Server) LogError(err error) {
	// if log channel is nil logging directly with logger
	if s.ErrorLogger.MyLogChannel() == nil {
		s.ErrorLogger.Log(err)
		return
	}

	// logging through channel
	s.ErrorLogger.MyLogChannel() <- err
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

// LogRenderErrorAndRedirect logs err, puts error message with template name in session, and redirects to url
func (s *Server) LogRenderErrorAndRedirect(w http.ResponseWriter, r *http.Request, template string, err error, url string) {
	e := ServerError{
		Prompt: fmt.Sprintf(`unable to render "%s" template`, template),
		URL:    r.URL.Path,
		Err:    err,
	}

	app.Session.Put(r.Context(), "error", e.Prompt)
	s.LogError(err)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
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

type ServerError struct {
	Prompt string
	URL    string
	Err    error
}

func (e ServerError) Error() string {
	text := util.NewText()

	if e.Err != nil {
		text.AddLineIndent(e.Err.Error(), "\t")
	}

	if e.Prompt != "" {
		text.AddLineIndent(fmt.Sprint("PROMPT: ", e.Prompt), "\t")
	}
	if e.URL != "" {
		text.AddLineIndent(fmt.Sprint("URL: ", e.URL), "\t")
	}

	return text.String()
}

func (s *Server) Mail(data mailers.MailData) {
	var err error
	if s.Mailer.MyMailChannel() == nil {
		err = s.Mailer.SendMail(data)
		if err != nil {
			s.LogError(err)
		}
	}

	s.Mailer.MyMailChannel() <- data
}
