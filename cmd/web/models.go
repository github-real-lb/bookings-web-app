package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/github-real-lb/bookings-web-app/util/forms"
)

// TemplateData holds data sent from handlers to templates
type TemplateData struct {
	CSRFToken string // Security Token to prevent Cross Site Request Forgery (CSRF)

	Data      map[string]any
	FloatMap  map[string]float32
	IntMap    map[string]int
	StringMap map[string]string

	Form *forms.Form

	Flash   string // Flash message
	Warning string // Warning message
	Error   string // Error message
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
	Room      Room      `json:"room"`
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Marshal returns the data of r
func (r *Reservation) Marshal() map[string]string {
	data := make(map[string]string)
	data["id"] = fmt.Sprint(r.ID)
	data["code"] = r.Code
	data["first_name"] = r.FirstName
	data["last_name"] = r.LastName
	data["email"] = r.Email
	data["phone"] = r.Phone
	data["start_date"] = r.StartDate.Format("2006-01-02")
	data["end_date"] = r.EndDate.Format("2006-01-02")
	data["room_id"] = fmt.Sprint(r.Room.ID)
	data["room_name"] = fmt.Sprint(r.Room.Name)
	data["room_description"] = fmt.Sprint(r.Room.Description)
	data["room_image_filename"] = fmt.Sprint(r.Room.ImageFilename)
	data["notes"] = r.Notes
	data["created_at"] = r.CreatedAt.Format(time.RFC3339)
	data["updated_at"] = r.StartDate.Format(time.RFC3339)
	return data
}

// Unmarshal parse data into r
func (r *Reservation) Unmarshal(data map[string]string) error {
	var err error = nil

	if _, ok := data["id"]; ok {
		r.ID, err = strconv.ParseInt(data["id"], 10, 64)
		if err != nil {
			return err
		}
	}

	if _, ok := data["first_name"]; ok {
		r.FirstName = data["first_name"]
	}

	if _, ok := data["last_name"]; ok {
		r.LastName = data["last_name"]
	}

	if _, ok := data["email"]; ok {
		r.Email = data["email"]
	}

	if _, ok := data["phone"]; ok {
		r.Phone = data["phone"]
	}

	if v, ok := data["start_date"]; ok {
		r.StartDate, err = time.Parse("2006-01-02", v[:10])
		if err != nil {
			return err
		}
	}

	if v, ok := data["end_date"]; ok {
		r.EndDate, err = time.Parse("2006-01-02", v[:10])
		if err != nil {
			return err
		}
	}

	if _, ok := data["room_id"]; ok {
		r.Room.ID, err = strconv.ParseInt(data["room_id"], 10, 64)
		if err != nil {
			return err
		}
	}

	if _, ok := data["room_name"]; ok {
		r.Room.Name = data["room_name"]
	}

	if _, ok := data["room_description"]; ok {
		r.Room.Description = data["room_description"]
	}

	if _, ok := data["room_image_filename"]; ok {
		r.Room.ImageFilename = data["room_image_filename"]
	}

	if _, ok := data["notes"]; ok {
		r.Notes = data["notes"]
	}

	if _, ok := data["created_at"]; ok {
		r.CreatedAt, err = time.Parse(time.RFC3339, data["created_at"])
		if err != nil {
			return err
		}
	}

	if _, ok := data["updated_at"]; ok {
		r.UpdatedAt, err = time.Parse(time.RFC3339, data["updated_at"])
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
func (r *Restriction) Marshal() map[string]string {
	data := make(map[string]string)
	data["id"] = fmt.Sprint(r.ID)
	data["name"] = r.Name
	data["created_at"] = r.CreatedAt.Format(time.RFC3339)
	data["updated_at"] = r.UpdatedAt.Format(time.RFC3339)
	return data
}

// Unmarshal parse data into r
func (r *Restriction) Unmarshal(data map[string]string) error {
	var err error = nil

	if _, ok := data["id"]; ok {
		r.ID, err = strconv.ParseInt(data["id"], 10, 64)
		if err != nil {
			return err
		}
	}

	if _, ok := data["name"]; ok {
		r.Name = data["name"]
	}

	if _, ok := data["created_at"]; ok {
		r.CreatedAt, err = time.Parse(time.RFC3339, data["created_at"])
		if err != nil {
			return err
		}
	}

	if _, ok := data["updated_at"]; ok {
		r.UpdatedAt, err = time.Parse(time.RFC3339, data["updated_at"])
		if err != nil {
			return err
		}
	}

	return err
}

type Room struct {
	ID            int64     `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	ImageFilename string    `json:"image_filename"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Rooms []Room

// Marshal returns data of r
func (r *Room) Marshal() map[string]string {
	data := make(map[string]string)
	data["id"] = fmt.Sprint(r.ID)
	data["name"] = r.Name
	data["description"] = r.Description
	data["image_filename"] = r.ImageFilename
	data["created_at"] = r.CreatedAt.Format(time.RFC3339)
	data["updated_at"] = r.UpdatedAt.Format(time.RFC3339)
	return data
}

// Unmarshal parse data into r
func (r *Room) Unmarshal(data map[string]string) error {
	var err error = nil

	if v, ok := data["id"]; ok {
		r.ID, err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			return err
		}
	}

	if _, ok := data["name"]; ok {
		r.Name = data["name"]
	}

	if _, ok := data["description"]; ok {
		r.Description = data["description"]
	}

	if _, ok := data["image_filename"]; ok {
		r.ImageFilename = data["image_filename"]
	}

	if _, ok := data["created_at"]; ok {
		r.CreatedAt, err = time.Parse(time.RFC3339, data["created_at"])
		if err != nil {
			return err
		}
	}

	if _, ok := data["updated_at"]; ok {
		r.UpdatedAt, err = time.Parse(time.RFC3339, data["updated_at"])
		if err != nil {
			return err
		}
	}

	return err
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
func (r *RoomRestriction) Marshal() map[string]string {
	data := make(map[string]string)
	data["id"] = fmt.Sprint(r.ID)
	data["start_date"] = r.StartDate.Format("2006-01-02")
	data["end_date"] = r.EndDate.Format("2006-01-02")
	data["room_id"] = fmt.Sprint(r.RoomID)
	data["reservation_id"] = fmt.Sprint(r.ReservationID)
	data["restrictions_id"] = fmt.Sprint(r.RestrictionsID)
	data["created_at"] = r.CreatedAt.Format(time.RFC3339)
	data["updated_at"] = r.UpdatedAt.Format(time.RFC3339)
	return data
}

// Unmarshal parse data into r
func (r *RoomRestriction) Unmarshal(data map[string]string) error {
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
