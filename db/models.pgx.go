package db

import (
	"fmt"
	"time"
)

// Marshal returns data of r
func (r *Reservation) Marshal() map[string]string {
	data := make(map[string]string)
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
	data["updated_at"] = r.UpdatedAt.Time.Format(time.RFC3339)
	return data
}

// Marshal returns data of r
func (r *Restriction) Marshal() map[string]string {
	data := make(map[string]string)
	data["id"] = fmt.Sprint(r.ID)
	data["name"] = r.Name
	data["created_at"] = r.CreatedAt.Time.Format(time.RFC3339)
	data["updated_at"] = r.UpdatedAt.Time.Format(time.RFC3339)
	return data
}

// Marshal returns data of r
func (r *RoomRestriction) Marshal() map[string]string {
	data := make(map[string]string)
	data["id"] = fmt.Sprint(r.ID)
	data["start_date"] = r.StartDate.Time.Format("2006-01-02")
	data["end_date"] = r.EndDate.Time.Format("2006-01-02")
	data["room_id"] = fmt.Sprint(r.RoomID)
	data["reservation_id"] = fmt.Sprint(r.ReservationID.Int64)
	data["restrictions_id"] = fmt.Sprint(r.RestrictionsID)
	data["created_at"] = r.CreatedAt.Time.Format(time.RFC3339)
	data["updated_at"] = r.UpdatedAt.Time.Format(time.RFC3339)
	return data
}
