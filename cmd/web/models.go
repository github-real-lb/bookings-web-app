package main

import (
	"time"

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

// Reservation is used to hold reservation data
type Reservation struct {
	ID        int64     `json:"id"`
	Code      string    `json:"code"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	RoomID    int64     `json:"room_id"`
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Restriction is used to hold different type of restrictions for rooms availabilty
type Restriction struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Room struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// RoomRestriction is used to hold room restriction data
type RoomRestriction struct {
	ID             int64     `json:"id"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
	RoomID         int64     `json:"room_id"`
	ReservationID  int64     `json:"reservation_id"`
	RestrictionsID int64     `json:"restrictions_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type User struct {
	ID          int64     `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	AccessLevel int64     `json:"access_level"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
