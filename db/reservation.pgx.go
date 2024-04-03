package db

import (
	"strconv"
	"time"

	"github.com/github-real-lb/bookings-web-app/util/config"
	"github.com/jackc/pgx/v5/pgtype"
)

// Unmarshal parse data into p
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
		t, err = time.Parse(config.DateLayout, v)
		if err != nil {
			return err
		}

		p.StartDate = pgtype.Date{
			Time:  t,
			Valid: true,
		}
	}

	if v, ok := data["end_date"]; ok {
		t, err = time.Parse(config.DateLayout, v)
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

// Unmarshal parse data into p
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
		t, err = time.Parse(config.DateLayout, v)
		if err != nil {
			return err
		}

		p.StartDate = pgtype.Date{
			Time:  t,
			Valid: true,
		}
	}

	if v, ok := data["end_date"]; ok {
		t, err = time.Parse(config.DateLayout, v)
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
