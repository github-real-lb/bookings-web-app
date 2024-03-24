package main

import "github.com/github-real-lb/bookings-web-app/internal/forms"

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

// Reservation is used to hold reservation data
type Reservation struct {
	FirstName string
	LastName  string
	Email     string
	Phone     string
}
