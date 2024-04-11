package db

import (
	"strconv"
)

// Unmarshal parse data into p
func (p *CreateReservationParams) Unmarshal(data map[string]string) error {
	var err error = nil

	p.Code = data["code"]
	p.FirstName = data["first_name"]
	p.LastName = data["last_name"]
	p.Email = data["email"]

	if v, ok := data["phone"]; ok {
		err = p.Phone.Scan(v)
		if err != nil {
			return err
		}
	}

	if v, ok := data["start_date"]; ok {
		err = p.StartDate.Scan(v)
		if err != nil {
			return err
		}
	}

	if v, ok := data["end_date"]; ok {
		err = p.EndDate.Scan(v)
		if err != nil {
			return err
		}
	}

	if v, ok := data["room_id"]; ok {
		p.RoomID, err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			return err
		}
	}

	if v, ok := data["notes"]; ok {
		err = p.Notes.Scan(v)
		if err != nil {
			return err
		}
	}

	return err
}

// Unmarshal parse data into p
func (p *UpdateReservationParams) Unmarshal(data map[string]string) error {
	var err error = nil

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

	if v, ok := data["phone"]; ok {
		err = p.Phone.Scan(v)
		if err != nil {
			return err
		}
	}

	if v, ok := data["start_date"]; ok {
		err = p.StartDate.Scan(v)
		if err != nil {
			return err
		}
	}

	if v, ok := data["end_date"]; ok {
		err = p.EndDate.Scan(v)
		if err != nil {
			return err
		}
	}

	if v, ok := data["room_id"]; ok {
		p.RoomID, err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			return err
		}
	}

	if v, ok := data["notes"]; ok {
		err = p.Notes.Scan(v)
		if err != nil {
			return err
		}
	}

	if v, ok := data["updated_at"]; ok {
		err = p.UpdatedAt.Scan(v)
		if err != nil {
			return err
		}
	}

	return err
}
