package db

import (
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// Unmarshal parse data into p
func (p *UpdateRestrictionParams) Unmarshal(data map[string]string) error {
	var err error = nil
	var t time.Time

	if v, ok := data["id"]; ok {
		p.ID, err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			return err
		}
	}

	p.Name = data["name"]

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
