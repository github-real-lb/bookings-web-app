package main

import (
	"fmt"
	"time"

	"github.com/github-real-lb/bookings-web-app/db"
	"github.com/github-real-lb/bookings-web-app/util/forms"
)

// Listing holds the data of the property for which bookings are made
type Listing struct {
	Title   string `json:"title"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
}

// TemplateData holds data sent from handlers to templates
type TemplateData struct {
	CSRFToken string // Security Token to prevent Cross Site Request Forgery (CSRF)
	Data      map[string]any

	Form *forms.Form

	IsAuthenticated bool // Determines if a user is logged in

	Listing Listing // Data of the property

	Error   string // Error message
	Flash   string // Success message
	Warning string // Warning message
}

// Reservation holds reservation data
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
	Room      Room      `json:"room"`
}

// Room holds hotel room data
type Room struct {
	ID            int64     `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	ImageFilename string    `json:"image_filename"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Restriction is the database restriction enum
type Restriction db.Restriction

const (
	RestrictionReservation Restriction = Restriction(db.RestrictionReservation)
	RestrictionOwnerBlock  Restriction = Restriction(db.RestrictionOwnerBlock)
)

func (r *Restriction) Scan(src any) error {
	switch s := src.(type) {
	case []byte:
		*r = Restriction(s)
	case string:
		*r = Restriction(s)
	default:
		return fmt.Errorf("unsupported scan type for Restriction: %T", src)
	}
	return nil
}

// RoomRestriction holds room restriction data
type RoomRestriction struct {
	ID            int64       `json:"id"`
	StartDate     time.Time   `json:"start_date"`
	EndDate       time.Time   `json:"end_date"`
	RoomID        int64       `json:"room_id"`
	ReservationID int64       `json:"reservation_id"`
	Restriction   Restriction `json:"restriction"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
}

// User holds user data
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
