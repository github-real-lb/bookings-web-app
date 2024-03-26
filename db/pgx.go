package db

import (
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// Unmarshal parse data into r
func (r *CreateReservationParams) Unmarshal(data map[string]string) error {
	var err error = nil

	r.Code = data["code"]
	r.FirstName = data["first_name"]
	r.LastName = data["last_name"]
	r.Email = data["email"]
	r.Phone = pgtype.Text{
		String: data["phone"],
		Valid:  data["phone"] != "",
	}

	if value, exist := data["start_date"]; exist {
		date, err := time.Parse("2006-01-02", value)
		if err != nil {
			return err
		}

		r.StartDate = pgtype.Date{
			Time:  date,
			Valid: true,
		}
	}

	if value, exist := data["end_date"]; exist {
		date, err := time.Parse("2006-01-02", value)
		if err != nil {
			return err
		}

		r.EndDate = pgtype.Date{
			Time:  date,
			Valid: true,
		}
	}

	if value, exist := data["room_id"]; exist {
		r.RoomID, err = strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
	}

	r.Notes = pgtype.Text{
		String: data["notes"],
		Valid:  data["notes"] != "",
	}

	return err
}

// Marshal returns data of r
func (r *Reservation) Marshal() (data map[string]string) {
	data = make(map[string]string)
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
	data["updated_at"] = r.StartDate.Time.Format(time.RFC3339)
	return
}
