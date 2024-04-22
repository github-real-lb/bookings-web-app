package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/justinas/nosurf"
)

// NoSurf is a middleware that adds CSRF protection to all POST requests
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InDevelopmentMode(),
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}

// LogRequestsAndResponse is a middleware that logs requests received and their response time
func (s *Server) LogRequestsAndResponse(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Time the execution
		start := time.Now()
		next.ServeHTTP(w, r) // Pass control to the next handler
		duration := time.Since(start)

		s.LogInfo(fmt.Sprintf("%s %s received from %s and handled in %v", r.Method, r.URL.Path, r.RemoteAddr, duration))
	})
}

// Auth is a middleware that restrict access to authenticated users
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !IsAuthenticated(r) {
			app.Session.Put(r.Context(), "error", "Pleasae log first!")
			http.Redirect(w, r, "/user/login", http.StatusTemporaryRedirect)
			return
		}

		next.ServeHTTP(w, r)
	})
}
