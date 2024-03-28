package db

import (
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// Marshal returns data of r
func (r *Restriction) Marshal() (data map[string]string) {
	data = make(map[string]string)
	data["id"] = fmt.Sprint(r.ID)
	data["name"] = r.Name
	data["created_at"] = r.CreatedAt.Time.Format(time.RFC3339)
	data["updated_at"] = r.UpdatedAt.Time.Format(time.RFC3339)
	return
}

// Unmarshal parse data into r
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
