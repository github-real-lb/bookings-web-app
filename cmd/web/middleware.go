package main

import (
	"net/http"

	"github.com/github-real-lb/bookings-web-app/internal/config"
	"github.com/justinas/nosurf"
)

// NoSurf adds CSRF protection to all POST requests
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.AppMode == config.ProductionMode,
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}
