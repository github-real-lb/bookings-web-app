package db

import (
	"strconv"
)

// Unmarshal parse data into p
func (p *CreateRoomRestrictionParams) Unmarshal(data map[string]string) error {
	var err error = nil

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

	if v, ok := data["reservation_id"]; ok {
		err = p.ReservationID.Scan(v)
		if err != nil {
			return err
		}
	}

	if v, ok := data["restriction"]; ok {
		err = p.Restriction.Scan(v)
		if err != nil {
			return err
		}
	}

	return err
}

// Unmarshal parse data into p
func (p *UpdateRoomRestrictionParams) Unmarshal(data map[string]string) error {
	var err error = nil

	if v, ok := data["id"]; ok {
		p.ID, err = strconv.ParseInt(v, 10, 64)
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

	if v, ok := data["reservation_id"]; ok {
		err = p.ReservationID.Scan(v)
		if err != nil {
			return err
		}
	}

	if v, ok := data["restriction"]; ok {
		err = p.Restriction.Scan(v)
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
