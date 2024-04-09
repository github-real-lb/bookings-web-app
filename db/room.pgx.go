package db

import (
	"strconv"
)

// Unmarshal parse data into p
func (p *CheckRoomAvailabilityParams) Unmarshal(data map[string]string) error {
	var err error = nil

	if v, ok := data["room_id"]; ok {
		p.RoomID, err = strconv.ParseInt(v, 10, 64)
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

	if v, ok := data["limit"]; ok {
		i, err := strconv.Atoi(v)
		if err != nil {
			return err
		}

		p.Limit = int32(i)
	}

	if v, ok := data["offset"]; ok {
		i, err := strconv.Atoi(v)
		if err != nil {
			return err
		}

		p.Offset = int32(i)
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

	return err
}
