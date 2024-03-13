package main

import (
	"net/http"

	"github.com/github-real-lb/bookings-web-app/pkg/config"
	"github.com/github-real-lb/bookings-web-app/pkg/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	// add middleware that recover from panics
	mux.Use(middleware.Recoverer)

	// add middleware that loads and saves and session on every request
	mux.Use(app.Session.LoadAndSave)

	// add middleware that provides CSRF protection to all POST requests
	mux.Use(NoSurf)

	// setting routes
	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/make-reservation", handlers.Repo.Reservation)

	mux.Get("/generals-quarters", handlers.Repo.Generals)
	mux.Post("/generals-quarters", handlers.Repo.PostAvailability)

	mux.Get("/majors-suite", handlers.Repo.Majors)
	mux.Post("/majors-suite", handlers.Repo.PostAvailability)

	mux.Get("/contact", handlers.Repo.Contact)

	mux.Get("/search-availability", handlers.Repo.Availability)
	mux.Post("/search-availability", handlers.Repo.PostAvailability)
	mux.Get("/search-availability-json", handlers.Repo.AvailabilityJSON)

	// setting file server
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
