package db

import (
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
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
	data["start_date"] = r.StartDate.Time.Format("2006-01-02")
	data["end_date"] = r.EndDate.Time.Format("2006-01-02")
	data["room_id"] = fmt.Sprint(r.RoomID)
	data["notes"] = r.Notes.String
	data["created_at"] = r.CreatedAt.Time.Format(time.RFC3339)
	data["updated_at"] = r.UpdatedAt.Time.Format(time.RFC3339)
	return data
}

func (r *Reservation) Unmarshal(data map[string]string) error {
	var err error = nil
	var t time.Time

	if _, ok := data["id"]; ok {
		r.ID, err = strconv.ParseInt(data["id"], 10, 64)
		if err != nil {
			return err
		}
	}

	if _, ok := data["code"]; ok {
		r.Code = data["code"]
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
		r.Phone = pgtype.Text{
			String: data["phone"],
			Valid:  true,
		}
	}

	if _, ok := data["start_date"]; ok {
		t, err = time.Parse("2006-01-02", data["start_date"])
		if err != nil {
			return err
		}

		r.StartDate = pgtype.Date{
			Time:  t,
			Valid: true,
		}
	}

	if _, ok := data["end_date"]; ok {
		t, err = time.Parse("2006-01-02", data["end_date"])
		if err != nil {
			return err
		}

		r.EndDate = pgtype.Date{
			Time:  t,
			Valid: true,
		}
	}

	if _, ok := data["room_id"]; ok {
		r.RoomID, err = strconv.ParseInt(data["room_id"], 10, 64)
		if err != nil {
			return err
		}
	}

	if _, ok := data["notes"]; ok {
		r.Notes = pgtype.Text{
			String: data["notes"],
			Valid:  true,
		}
	}

	if _, ok := data["created_at"]; ok {
		t, err = time.Parse(time.RFC3339, data["created_at"])
		if err != nil {
			return err
		}

		r.CreatedAt = pgtype.Timestamptz{
			Time:  t,
			Valid: true,
		}
	}

	if _, ok := data["updated_at"]; ok {
		t, err = time.Parse(time.RFC3339, data["updated_at"])
		if err != nil {
			return err
		}

		r.UpdatedAt = pgtype.Timestamptz{
			Time:  t,
			Valid: true,
		}
	}

	return err
}

// Marshal returns data of r
func (r *Restriction) Marshal() map[string]string {
	data := make(map[string]string)
	data["id"] = fmt.Sprint(r.ID)
	data["name"] = r.Name
	data["created_at"] = r.CreatedAt.Time.Format(time.RFC3339)
	data["updated_at"] = r.UpdatedAt.Time.Format(time.RFC3339)
	return data
}

// Marshal returns data of r
func (r *Room) Marshal() map[string]string {
	data := make(map[string]string)
	data["id"] = fmt.Sprint(r.ID)
	data["name"] = r.Name
	data["description"] = r.Description
	data["image_filename"] = r.ImageFilename.String
	data["created_at"] = r.CreatedAt.Time.Format(time.RFC3339)
	data["updated_at"] = r.UpdatedAt.Time.Format(time.RFC3339)
	return data
}

// Marshal returns data of r
func (r *RoomRestriction) Marshal() map[string]string {
	data := make(map[string]string)
	data["id"] = fmt.Sprint(r.ID)
	data["start_date"] = r.StartDate.Time.Format("2006-01-02")
	data["end_date"] = r.EndDate.Time.Format("2006-01-02")
	data["room_id"] = fmt.Sprint(r.RoomID)
	data["reservation_id"] = fmt.Sprint(r.ReservationID.Int64)
	data["restrictions_id"] = fmt.Sprint(r.RestrictionsID)
	data["created_at"] = r.CreatedAt.Time.Format(time.RFC3339)
	data["updated_at"] = r.UpdatedAt.Time.Format(time.RFC3339)
	return data
}
