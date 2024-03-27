package db

import (
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

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

// Unmarshal parse data into r
func (p *CreateReservationParams) Unmarshal(data map[string]string) error {
	var err error = nil
	var t time.Time

	p.Code = data["code"]
	p.FirstName = data["first_name"]
	p.LastName = data["last_name"]
	p.Email = data["email"]
	p.Phone = pgtype.Text{
		String: data["phone"],
		Valid:  data["phone"] != "",
	}

	if v, ok := data["start_date"]; ok {
		t, err = time.Parse("2006-01-02", v)
		if err != nil {
			return err
		}

		p.StartDate = pgtype.Date{
			Time:  t,
			Valid: true,
		}
	}

	if v, ok := data["end_date"]; ok {
		t, err = time.Parse("2006-01-02", v)
		if err != nil {
			return err
		}

		p.EndDate = pgtype.Date{
			Time:  t,
			Valid: true,
		}
	}

	if v, ok := data["room_id"]; ok {
		p.RoomID, err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			return err
		}
	}

	p.Notes = pgtype.Text{
		String: data["notes"],
		Valid:  data["notes"] != "",
	}

	return err
}

func (p *UpdateReservationParams) Unmarshal(data map[string]string) error {
	var err error = nil
	var t time.Time

	if v, ok := data["id"]; ok {
		p.ID, err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			return err
		}
	}

	p.Code = data["code"]
	p.FirstName = data["first_name"]
	p.LastName = data["last_name"]
	p.Email = data["email"]
	p.Phone = pgtype.Text{
		String: data["phone"],
		Valid:  data["phone"] != "",
	}

	if v, ok := data["start_date"]; ok {
		t, err = time.Parse("2006-01-02", v)
		if err != nil {
			return err
		}

		p.StartDate = pgtype.Date{
			Time:  t,
			Valid: true,
		}
	}

	if v, ok := data["end_date"]; ok {
		t, err = time.Parse("2006-01-02", v)
		if err != nil {
			return err
		}

		p.EndDate = pgtype.Date{
			Time:  t,
			Valid: true,
		}
	}

	if v, ok := data["room_id"]; ok {
		p.RoomID, err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			return err
		}
	}

	p.Notes = pgtype.Text{
		String: data["notes"],
		Valid:  data["notes"] != "",
	}

	if v, ok := data["updated_at"]; ok {
		t, err = time.Parse(time.RFC3339, v)
		if err != nil {
			return err
		}

		p.UpdatedAt = pgtype.Timestamptz{
			Time:  t,
			Valid: true,
		}
	}

	return err
}
