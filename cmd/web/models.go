package main

import (
	"fmt"
	"strconv"
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

// Marshal returns the data of r
func (r *Reservation) Marshal() (data map[string]string) {
	data = make(map[string]string)
	data["id"] = fmt.Sprint(r.ID)
	data["code"] = r.Code
	data["first_name"] = r.FirstName
	data["last_name"] = r.LastName
	data["email"] = r.Email
	data["phone"] = r.Phone
	data["start_date"] = r.StartDate.Format("2006-01-02")
	data["end_date"] = r.EndDate.Format("2006-01-02")
	data["room_id"] = fmt.Sprint(r.RoomID)
	data["notes"] = r.Notes
	data["created_at"] = r.CreatedAt.Format(time.RFC3339)
	data["updated_at"] = r.StartDate.Format(time.RFC3339)
	return
}

// Unmarshal parse data into r
func (r *Reservation) Unmarshal(data map[string]string) error {
	var err error = nil

	if value, exist := data["id"]; exist {
		r.ID, err = strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
	}

	r.FirstName = data["first_name"]
	r.LastName = data["last_name"]
	r.Email = data["email"]
	r.Phone = data["phone"]

	if value, exist := data["start_date"]; exist {
		r.StartDate, err = time.Parse("2006-01-02", value)
		if err != nil {
			return err
		}
	}

	if value, exist := data["end_date"]; exist {
		r.EndDate, err = time.Parse("2006-01-02", value)
		if err != nil {
			return err
		}
	}

	if value, exist := data["room_id"]; exist {
		r.RoomID, err = strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
	}

	r.Notes = data["notes"]

	return err
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
