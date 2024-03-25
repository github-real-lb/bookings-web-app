package main

import (
	"net/http"

	"github.com/github-real-lb/bookings-web-app/db"
)

// Server handles all routing and provides all database functions
type Server struct {
	Router        *http.Server
	DatabaseStore db.DatabaseStore
}

// NewServer returns a new Server
func NewServer(store db.DatabaseStore) *Server {
	server := Server{
		Router: &http.Server{
			Addr:    app.ServerAddress,
			Handler: NewHandler(HandlerStore{}),
		},
		DatabaseStore: store,
	}

	return &server
}
