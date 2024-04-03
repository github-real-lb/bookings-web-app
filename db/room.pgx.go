package db

import (
	"strconv"
	"time"

	"github.com/github-real-lb/bookings-web-app/util/config"
	"github.com/jackc/pgx/v5/pgtype"
)

// Unmarshal parse data into p
func (p *CheckRoomAvailabilityParams) Unmarshal(data map[string]string) error {
	var err error = nil
	var t time.Time

	if v, ok := data["room_id"]; ok {
		p.RoomID, err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			return err
		}
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

	return err
}

// Unmarshal parse data into p
func (p *CreateRoomParams) Unmarshal(data map[string]string) {
	if v, ok := data["name"]; ok {
		p.Name = v
	}

	if v, ok := data["description"]; ok {
		p.Description = v
	}

	if v, ok := data["image_filename"]; ok {
		p.ImageFilename = v
	}
}

// Unmarshal parse data into p
func (p *ListAvailableRoomsParams) Unmarshal(data map[string]string) error {
	var err error = nil
	var t time.Time

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

	return err
}
