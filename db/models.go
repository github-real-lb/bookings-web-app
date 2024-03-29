// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package db

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Reservation struct {
	ID        int64              `json:"id"`
	Code      string             `json:"code"`
	FirstName string             `json:"first_name"`
	LastName  string             `json:"last_name"`
	Email     string             `json:"email"`
	Phone     pgtype.Text        `json:"phone"`
	StartDate pgtype.Date        `json:"start_date"`
	EndDate   pgtype.Date        `json:"end_date"`
	RoomID    int64              `json:"room_id"`
	Notes     pgtype.Text        `json:"notes"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
}

type Restriction struct {
	ID        int64              `json:"id"`
	Name      string             `json:"name"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
}

type Room struct {
	ID        int64              `json:"id"`
	Name      string             `json:"name"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
}

type RoomRestriction struct {
	ID             int64              `json:"id"`
	StartDate      pgtype.Date        `json:"start_date"`
	EndDate        pgtype.Date        `json:"end_date"`
	RoomID         int64              `json:"room_id"`
	ReservationID  pgtype.Int8        `json:"reservation_id"`
	RestrictionsID int64              `json:"restrictions_id"`
	CreatedAt      pgtype.Timestamptz `json:"created_at"`
	UpdatedAt      pgtype.Timestamptz `json:"updated_at"`
}

type User struct {
	ID          int64              `json:"id"`
	FirstName   string             `json:"first_name"`
	LastName    string             `json:"last_name"`
	Email       string             `json:"email"`
	Password    string             `json:"password"`
	AccessLevel int64              `json:"access_level"`
	CreatedAt   pgtype.Timestamptz `json:"created_at"`
	UpdatedAt   pgtype.Timestamptz `json:"updated_at"`
}
