package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	http.Server
}

// NewServer returns New http.Server to serves all HTTP requests
func NewServer(addr string) *Server {
	mux := chi.NewRouter()
	server := Server{
		Server: http.Server{
			Addr:    addr,
			Handler: mux,
		},
	}

	// add middleware that recover from panics
	mux.Use(middleware.Recoverer)

	// add middleware that loads and saves and session on every request
	mux.Use(app.Session.LoadAndSave)

	// add middleware that provides CSRF protection to all POST requests
	mux.Use(NoSurf)

	// setting routes
	mux.Get("/", server.Home)
	mux.Get("/about", server.About)

	mux.Get("/generals-quarters", server.Generals)
	mux.Post("/generals-quarters", server.PostAvailability)

	mux.Get("/majors-suite", server.Majors)
	mux.Post("/majors-suite", server.PostAvailability)

	mux.Get("/contact", server.Contact)

	mux.Get("/search-availability", server.Availability)
	mux.Post("/search-availability", server.PostAvailability)
	mux.Post("/search-availability-json", server.AvailabilityJSON)

	mux.Get("/make-reservation", server.Reservation)
	mux.Post("/make-reservation", server.PostReservation)
	mux.Get("/reservation-summary", server.ReservationSummary)

	// setting file server
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return &server
}
