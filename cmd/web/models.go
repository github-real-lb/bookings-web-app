package main

import (
	"github.com/github-real-lb/bookings-web-app/db"
	"github.com/github-real-lb/bookings-web-app/internal/forms"
)

// TemplateData holds data sent from handlers to templates
type TemplateData struct {
	StringMap map[string]string
	IntMap    map[string]int
	FloatMap  map[string]float32
	Data      map[string]any
	CSRFToken string // Security Token to prevent Cross Site Request Forgery (CSRF)
	Flash     string // Flash message
	Warning   string // Warning message
	Error     string // Error message
	Form      *forms.Form
}

type User struct {
	db.User
}

type Room struct {
	db.Room
}

// Restriction is used to hold different type of restrictions for rooms availabilty
type Restriction struct {
	db.Restriction
}

// Reservation is used to hold reservation data
type Reservation struct {
	db.Reservation
	Room Room
}

// RoomRestriction is used to hold room restriction data
type RoomRestriction struct {
	db.RoomRestriction
	Reservation Reservation
}
