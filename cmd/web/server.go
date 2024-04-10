package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/github-real-lb/bookings-web-app/db"
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

// LogError logs error with message as a prefix
func (s *Server) LogError(r *http.Request, message string, err error) {
	message = fmt.Sprintf("messege: %s\nurl: %s", message, r.URL.Path)
	app.LogError(message, err)
}

// LogErrorAndRedirect logs error, put message in session, and redirect to url
func (s *Server) LogErrorAndRedirect(w http.ResponseWriter, r *http.Request, message string, err error, url string) {
	app.Session.Put(r.Context(), "error", message)

	message = fmt.Sprintf("PROMPT: %s\nURL: %s", message, r.URL.Path)
	app.LogError(message, err)
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
