package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/github-real-lb/bookings-web-app/db"
	"github.com/github-real-lb/bookings-web-app/util/loggers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Server handles all routing and provides all database functions
type Server struct {
	Router        *http.Server
	DatabaseStore db.DatabaseStore
}

// NewServer returns a new Server with Router and Database Store
func NewServer(store db.DatabaseStore) *Server {
	mux := chi.NewRouter()

	server := Server{
		Router: &http.Server{
			Addr:    app.ServerAddress,
			Handler: mux,
		},
		DatabaseStore: store,
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
	fmt.Printf("Starting http server on %s... \n", app.ServerAddress)
	err := s.Router.ListenAndServe()
	if err != nil && err.Error() != "http: Server closed" {
		app.Logger.LogError(loggers.ErrorData{
			Prefix: "error starting http server",
			Error:  err,
		})
	}
}

// Stop calls the http.Server Shutdown method
func (s *Server) Stop() {
	// create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Print("Shutting down http server... ")
	err := s.Router.Shutdown(ctx)
	if err != nil {
		app.Logger.LogError(loggers.ErrorData{
			Prefix: "error shuting down http server",
			Error:  err,
		})
	} else {
		fmt.Println("Success")
	}

}

// LogError logs error with message as a prefix
func (s *Server) LogError(r *http.Request, message string, err error) {
	message = fmt.Sprintf("\n\tPROMPT: %s\n\tURL: %s", message, r.URL.Path)
	app.Logger.ErrorChannel <- loggers.ErrorData{
		Prefix: message,
		Error:  err,
	}
}

// LogErrorAndRedirect logs error, puts message in session, and redirects to url
func (s *Server) LogErrorAndRedirect(w http.ResponseWriter, r *http.Request, message string, err error, url string) {
	app.Session.Put(r.Context(), "error", message)

	message = fmt.Sprintf("\n\tPROMPT: %s\n\tURL: %s", message, r.URL.Path)
	app.Logger.ErrorChannel <- loggers.ErrorData{
		Prefix: message,
		Error:  err,
	}

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// LogRenderErrorAndRedirect logs error with template rendering.
// It puts error message in session with template name, and redirects to url
func (s *Server) LogRenderErrorAndRedirect(w http.ResponseWriter, r *http.Request, template string, err error, url string) {
	message := fmt.Sprintf(`unable to render "%s" template`, template)
	app.Session.Put(r.Context(), "error", message)

	message = fmt.Sprintf("\n\tPROMPT: %s\n\tURL: %s", message, r.URL.Path)

	app.Logger.ErrorChannel <- loggers.ErrorData{
		Prefix: message,
		Error:  err,
	}

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// ResponseJSON write v to w as json response
func (s *Server) ResponseJSON(w http.ResponseWriter, r *http.Request, v any) error {
	bs, err := json.Marshal(v)
	if err != nil {
		s.LogError(r, "unable to marshal json response", err)
		return err
	}

	_, err = w.Write(bs)
	if err != nil {
		s.LogError(r, "unable to write json response", err)
		return err
	}

	return nil
}
