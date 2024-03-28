package db

import (
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// Marshal returns data of r
func (r *RoomRestriction) Marshal() (data map[string]string) {
	data = make(map[string]string)
	data["id"] = fmt.Sprint(r.ID)
	data["start_date"] = r.StartDate.Time.Format("2006-01-02")
	data["end_date"] = r.EndDate.Time.Format("2006-01-02")
	data["room_id"] = fmt.Sprint(r.RoomID)
	data["reservation_id"] = fmt.Sprint(r.ReservationID.Int64)
	data["restrictions_id"] = fmt.Sprint(r.RestrictionsID)
	data["created_at"] = r.CreatedAt.Time.Format(time.RFC3339)
	data["updated_at"] = r.UpdatedAt.Time.Format(time.RFC3339)
	return
}

// Unmarshal parse data into r
func (p *CreateRoomRestrictionParams) Unmarshal(data map[string]string) error {
	var err error = nil
	var t time.Time
	var i int64

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

	if v, ok := data["reservation_id"]; ok {
		i, err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			return err
		}

		p.ReservationID = pgtype.Int8{
			Int64: i,
			Valid: true,
		}
	}

	if v, ok := data["restrictions_id"]; ok {
		p.RestrictionsID, err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			return err
		}
	}

	return err
}

// Unmarshal parse data into r
func (p *UpdateRoomRestrictionParams) Unmarshal(data map[string]string) error {
	var err error = nil
	var t time.Time
	var i int64

	if v, ok := data["id"]; ok {
		p.ID, err = strconv.ParseInt(v, 10, 64)
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

	if v, ok := data["room_id"]; ok {
		p.RoomID, err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			return err
		}
	}

	if v, ok := data["reservation_id"]; ok {
		i, err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			return err
		}

		p.ReservationID = pgtype.Int8{
			Int64: i,
			Valid: true,
		}
	}

	if v, ok := data["restrictions_id"]; ok {
		p.RestrictionsID, err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			return err
		}
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
