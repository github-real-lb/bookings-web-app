package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"github.com/github-real-lb/bookings-web-app/db"
	"github.com/github-real-lb/bookings-web-app/util"
	"github.com/github-real-lb/bookings-web-app/util/config"
	"github.com/github-real-lb/bookings-web-app/util/forms"
)

// Listing holds the data of the property for which bookings are made
type Listing struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
}

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
	Room      Room      `json:"room"`
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

const ReservationCodeLenght = 7

// GenerateReservationCode generate the reservation code.
func (r *Reservation) GenerateReservationCode() {
	// concatenating the current time with the reservation last name
	s := fmt.Sprintf("%v %s %v %v", util.RandomDatetime().Format(time.RFC3339Nano), r.LastName, r.StartDate, r.EndDate)

	// Generate SHA-256 hash of the concatenated string
	hash := sha256.New()
	hash.Write([]byte(s))

	// Generate the SHA256 checksum of the data written so far and convert to hexadecimal string
	hashString := hex.EncodeToString(hash.Sum(nil))

	// build code string
	code := make([]byte, ReservationCodeLenght)
	digitsFound := 0
	digitsMax := ReservationCodeLenght / 2
	lettersFound := 0
	lettersMax := ReservationCodeLenght - digitsMax

	for _, v := range []byte(hashString) {
		if (digitsFound < digitsMax) && (v >= 49 && v <= 57) {
			// adds digits to code if not enought digits were found and if v is a digit between 1-9
			code[digitsFound*2+1] = v
			digitsFound++
		} else if (lettersFound < lettersMax) && ((v >= 97 && v <= 104) || (v >= 106 && v <= 110) || (v >= 112 && v <= 122)) {
			// adds letters to code if not enought letters were found check if v is a letter except for 'i' or 'o'
			code[lettersFound*2] = v - 32
			lettersFound++
		}

		if digitsFound+lettersFound == ReservationCodeLenght {
			r.Code = string(code)
			return
		}
	}
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

	if v, ok := data["code"]; ok {
		r.Code = v
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

// Room holds hotel room data
type Room struct {
	ID            int64     `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	ImageFilename string    `json:"image_filename"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Rooms holds a slice of Room
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

// Marshal returns data of r
func (r *RoomRestriction) Marshal() map[string]string {
	data := make(map[string]string)
	data["id"] = fmt.Sprint(r.ID)
	data["start_date"] = r.StartDate.Format(config.DateLayout)
	data["end_date"] = r.EndDate.Format(config.DateLayout)
	data["room_id"] = fmt.Sprint(r.RoomID)
	data["reservation_id"] = fmt.Sprint(r.ReservationID)
	data["restriction"] = string(RestrictionOwnerBlock)
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

	if v, ok := data["restriction"]; ok {
		err = r.Restriction.Scan(v)
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
