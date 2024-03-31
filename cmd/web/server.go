package main

import (
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
	mux.Get("/", server.Home)
	mux.Get("/about", server.AboutHandler)

	mux.Get("/rooms/{index}", server.RoomsHandler)
	mux.Get("/rooms/room/{name}", server.RoomHandler)
	mux.Post("/search-room-availability", server.PostSearchRoomAvailabilityHandler)

	mux.Get("/contact", server.ContactHandler)

	mux.Get("/search-availability", server.SearchAvailabilityHandler)
	mux.Post("/search-availability", server.PostSearchAvailabilityHandler)

	mux.Get("/make-reservation", server.MakeReservationHandler)
	mux.Post("/make-reservation", server.PostMakeReservationHandler)

	mux.Get("/choose-room/{index}", server.ChooseRoomHandler)

	mux.Get("/reservation-summary", server.ReservationSummaryHandler)

	// setting file server
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return &server
}
