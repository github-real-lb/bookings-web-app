package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/github-real-lb/bookings-web-app/util/config"
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

const ReservationCodeLenght = 8

// GenerateReservationCode generate the reservation code.
func (r *Reservation) GenerateReservationCode() error {
	// concatenating the current time with the reservation last name
	code := fmt.Sprint(time.Now().Format(config.DateLayout), r.LastName)

	// Generate SHA-256 hash of the concatenated string
	hash := sha256.New()
	_, err := hash.Write([]byte(code))
	if err != nil {
		return err
	}
	hashedBytes := hash.Sum(nil)

	// Convert the hashed bytes to hexadecimal string
	code = hex.EncodeToString(hashedBytes)[:ReservationCodeLenght]
	r.Code = strings.ToUpper(code)

	return nil
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
	data["start_date"] = r.StartDate.Format(config.DateLayout)
	data["end_date"] = r.EndDate.Format(config.DateLayout)
	data["room_id"] = fmt.Sprint(r.Room.ID)
	data["room_name"] = fmt.Sprint(r.Room.Name)
	data["room_description"] = fmt.Sprint(r.Room.Description)
	data["room_image_filename"] = fmt.Sprint(r.Room.ImageFilename)
	data["notes"] = r.Notes
	data["created_at"] = r.CreatedAt.Format(config.DateTimeLayout)
	data["updated_at"] = r.UpdatedAt.Format(config.DateTimeLayout)
	return data
}

// Unmarshal parse data into r
func (r *Reservation) Unmarshal(data map[string]string) error {
	var err error = nil

	if v, ok := data["id"]; ok {
		r.ID, err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			return err
		}
	}

	if v, ok := data["first_name"]; ok {
		r.FirstName = v
	}

	if v, ok := data["last_name"]; ok {
		r.LastName = v
	}

	if v, ok := data["email"]; ok {
		r.Email = v
	}

	if v, ok := data["phone"]; ok {
		r.Phone = v
	}

	if v, ok := data["start_date"]; ok {
		r.StartDate, err = time.Parse(config.DateLayout, v[:10])
		if err != nil {
			return err
		}
	}

	if v, ok := data["end_date"]; ok {
		r.EndDate, err = time.Parse(config.DateLayout, v[:10])
		if err != nil {
			return err
		}
	}

	if v, ok := data["room_id"]; ok {
		r.Room.ID, err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			return err
		}
	}

	if v, ok := data["room_name"]; ok {
		r.Room.Name = v
	}

	if v, ok := data["room_description"]; ok {
		r.Room.Description = v
	}

	if v, ok := data["room_image_filename"]; ok {
		r.Room.ImageFilename = v
	}

	if v, ok := data["notes"]; ok {
		r.Notes = v
	}

	if v, ok := data["created_at"]; ok {
		r.CreatedAt, err = time.Parse(config.DateTimeLayout, v)
		if err != nil {
			return err
		}
	}

	if v, ok := data["updated_at"]; ok {
		r.UpdatedAt, err = time.Parse(config.DateTimeLayout, v)
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
	data["created_at"] = r.CreatedAt.Format(config.DateTimeLayout)
	data["updated_at"] = r.UpdatedAt.Format(config.DateTimeLayout)
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

	if v, ok := data["name"]; ok {
		r.Name = v
	}

	if v, ok := data["description"]; ok {
		r.Description = v
	}

	if v, ok := data["image_filename"]; ok {
		r.ImageFilename = v
	}

	if v, ok := data["created_at"]; ok {
		r.CreatedAt, err = time.Parse(config.DateTimeLayout, v)
		if err != nil {
			return err
		}
	}

	if v, ok := data["updated_at"]; ok {
		r.UpdatedAt, err = time.Parse(config.DateTimeLayout, v)
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
	data["start_date"] = r.StartDate.Format(config.DateLayout)
	data["end_date"] = r.EndDate.Format(config.DateLayout)
	data["room_id"] = fmt.Sprint(r.RoomID)
	data["reservation_id"] = fmt.Sprint(r.ReservationID)
	data["restrictions_id"] = fmt.Sprint(r.RestrictionsID)
	data["created_at"] = r.CreatedAt.Format(config.DateTimeLayout)
	data["updated_at"] = r.UpdatedAt.Format(config.DateTimeLayout)
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
		r.StartDate, err = time.Parse(config.DateLayout, v)
		if err != nil {
			return err
		}
	}

	if v, ok := data["end_date"]; ok {
		r.EndDate, err = time.Parse(config.DateLayout, v)
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
		r.CreatedAt, err = time.Parse(config.DateTimeLayout, v)
		if err != nil {
			return err
		}
	}

	if v, ok := data["updated_at"]; ok {
		r.UpdatedAt, err = time.Parse(config.DateTimeLayout, v)
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
