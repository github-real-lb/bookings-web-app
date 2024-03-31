package db

import (
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// Unmarshal parse data into p
func (p *CheckRoomAvailabiltyParams) Unmarshal(data map[string]string) error {
	var err error = nil
	var t time.Time

	if v, ok := data["room_id"]; ok {
		p.RoomID, err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			return err
		}
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

	return err
}

// Unmarshal parse data into p
func (p *CreateRoomParams) Unmarshal(data map[string]string) {
	p.Name = data["name"]
	p.Description = data["description"]

	if v, ok := data["image_filename"]; ok {
		p.ImageFilename = pgtype.Text{
			String: v,
			Valid:  v != "",
		}
	}
}

// Unmarshal parse data into p
func (p *ListAvailableRoomsParams) Unmarshal(data map[string]string) error {
	var err error = nil
	var t time.Time

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

	return err
}
