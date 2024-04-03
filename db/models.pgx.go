package db

import (
	"fmt"
	"strconv"

	"github.com/github-real-lb/bookings-web-app/util/config"
)

// Marshal returns data of r
func (r *Reservation) Marshal() map[string]string {
	data := make(map[string]string)
	data["id"] = fmt.Sprint(r.ID)
	data["code"] = r.Code
	data["first_name"] = r.FirstName
	data["last_name"] = r.LastName
	data["email"] = r.Email
	data["phone"] = r.Phone.String
	data["start_date"] = r.StartDate.Time.Format(config.DateLayout)
	data["end_date"] = r.EndDate.Time.Format(config.DateLayout)
	data["room_id"] = fmt.Sprint(r.RoomID)
	data["notes"] = r.Notes.String
	data["created_at"] = r.CreatedAt.Time.Format(config.DateTimeLayout)
	data["updated_at"] = r.UpdatedAt.Time.Format(config.DateTimeLayout)
	return data
}

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
		err = r.Phone.Scan(v)
		if err != nil {
			return err
		}
	}

	if v, ok := data["start_date"]; ok {
		err = r.StartDate.Scan(v)
		if err != nil {
			return err
		}
	}

	if v, ok := data["end_date"]; ok {
		err = r.EndDate.Scan(v)
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

	if v, ok := data["notes"]; ok {
		err = r.Notes.Scan(v)
		if err != nil {
			return err
		}
	}

	if v, ok := data["created_at"]; ok {
		err = r.CreatedAt.Scan(v)
		if err != nil {
			return err
		}
	}

	if v, ok := data["updated_at"]; ok {
		err = r.UpdatedAt.Scan(v)
		if err != nil {
			return err
		}
	}

	return err
}

// Marshal returns data of r
func (r *Room) Marshal() map[string]string {
	data := make(map[string]string)
	data["id"] = fmt.Sprint(r.ID)
	data["name"] = r.Name
	data["description"] = r.Description
	data["image_filename"] = r.ImageFilename
	data["created_at"] = r.CreatedAt.Time.Format(config.DateTimeLayout)
	data["updated_at"] = r.UpdatedAt.Time.Format(config.DateTimeLayout)
	return data
}

// Marshal returns data of r
func (r *RoomRestriction) Marshal() map[string]string {
	data := make(map[string]string)
	data["id"] = fmt.Sprint(r.ID)
	data["start_date"] = r.StartDate.Time.Format(config.DateLayout)
	data["end_date"] = r.EndDate.Time.Format(config.DateLayout)
	data["room_id"] = fmt.Sprint(r.RoomID)
	data["reservation_id"] = fmt.Sprint(r.ReservationID.Int64)
	data["restriction"] = string(r.Restriction)
	data["created_at"] = r.CreatedAt.Time.Format(config.DateTimeLayout)
	data["updated_at"] = r.UpdatedAt.Time.Format(config.DateTimeLayout)
	return data
}
