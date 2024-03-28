package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/github-real-lb/bookings-web-app/util/forms"
)

type StringMap map[string]string

// TemplateData holds data sent from handlers to templates
type TemplateData struct {
	StringMap StringMap
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
func (r *Reservation) Marshal() (data StringMap) {
	data = make(StringMap)
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
func (r *Reservation) Unmarshal(data StringMap) error {
	var err error = nil

	if v, ok := data["id"]; ok {
		r.ID, err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			return err
		}
	}

	r.FirstName = data["first_name"]
	r.LastName = data["last_name"]
	r.Email = data["email"]
	r.Phone = data["phone"]

	if v, ok := data["start_date"]; ok {
		r.StartDate, err = time.Parse("2006-01-02", v)
		if err != nil {
			return err
		}
	}

	if v, ok := data["end_date"]; ok {
		r.EndDate, err = time.Parse("2006-01-02", v)
		if err != nil {
			return err
		}
	}

	if v, ok := data["room_id"]; ok {
		r.RoomID, err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			return err
		}
	}

	r.Notes = data["notes"]

	if v, ok := data["created_at"]; ok {
		r.CreatedAt, err = time.Parse(time.RFC3339, v)
		if err != nil {
			return err
		}
	}

	if v, ok := data["updated_at"]; ok {
		r.UpdatedAt, err = time.Parse(time.RFC3339, v)
		if err != nil {
			return err
		}
	}

	return err
}

// Restriction is used to hold different type of restrictions for rooms availabilty
type Restriction struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Marshal returns data of r
func (r *Restriction) Marshal() (data StringMap) {
	data = make(StringMap)
	data["id"] = fmt.Sprint(r.ID)
	data["name"] = r.Name
	data["created_at"] = r.CreatedAt.Format(time.RFC3339)
	data["updated_at"] = r.UpdatedAt.Format(time.RFC3339)
	return
}

// Unmarshal parse data into r
func (r *Restriction) Unmarshal(data StringMap) error {
	var err error = nil

	if v, ok := data["id"]; ok {
		r.ID, err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			return err
		}
	}

	r.Name = data["name"]

	if v, ok := data["created_at"]; ok {
		r.CreatedAt, err = time.Parse(time.RFC3339, v)
		if err != nil {
			return err
		}
	}

	if v, ok := data["updated_at"]; ok {
		r.UpdatedAt, err = time.Parse(time.RFC3339, v)
		if err != nil {
			return err
		}
	}

	return err
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

// Marshal returns data of r
func (r *RoomRestriction) Marshal() (data StringMap) {
	data = make(StringMap)
	data["id"] = fmt.Sprint(r.ID)
	data["start_date"] = r.StartDate.Format("2006-01-02")
	data["end_date"] = r.EndDate.Format("2006-01-02")
	data["room_id"] = fmt.Sprint(r.RoomID)
	data["reservation_id"] = fmt.Sprint(r.ReservationID)
	data["restrictions_id"] = fmt.Sprint(r.RestrictionsID)
	data["created_at"] = r.CreatedAt.Format(time.RFC3339)
	data["updated_at"] = r.UpdatedAt.Format(time.RFC3339)
	return
}

// Unmarshal parse data into r
func (r *RoomRestriction) Unmarshal(data StringMap) error {
	var err error = nil

	if v, ok := data["id"]; ok {
		r.ID, err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			return err
		}
	}

	if v, ok := data["start_date"]; ok {
		r.StartDate, err = time.Parse("2006-01-02", v)
		if err != nil {
			return err
		}
	}

	if v, ok := data["end_date"]; ok {
		r.EndDate, err = time.Parse("2006-01-02", v)
		if err != nil {
			return err
		}
	}

	if v, ok := data["room_id"]; ok {
		r.RoomID, err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			return err
		}
	}

	if v, ok := data["reservation_id"]; ok {
		r.ReservationID, err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			return err
		}
	}

	if v, ok := data["restrictions_id"]; ok {
		r.RestrictionsID, err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			return err
		}
	}

	if v, ok := data["created_at"]; ok {
		r.CreatedAt, err = time.Parse(time.RFC3339, v)
		if err != nil {
			return err
		}
	}

	if v, ok := data["updated_at"]; ok {
		r.UpdatedAt, err = time.Parse(time.RFC3339, v)
		if err != nil {
			return err
		}
	}

	return err
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
