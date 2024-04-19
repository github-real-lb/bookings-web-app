// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package db

import (
	"context"
)

type Querier interface {
	CheckRoomAvailability(ctx context.Context, arg CheckRoomAvailabilityParams) (bool, error)
	CreateReservation(ctx context.Context, arg CreateReservationParams) (Reservation, error)
	CreateRoom(ctx context.Context, arg CreateRoomParams) (Room, error)
	CreateRoomRestriction(ctx context.Context, arg CreateRoomRestrictionParams) (RoomRestriction, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteAllReservations(ctx context.Context) error
	DeleteAllRoomRestrictions(ctx context.Context) error
	DeleteAllRooms(ctx context.Context) error
	DeleteReservation(ctx context.Context, id int64) error
	DeleteRoom(ctx context.Context, id int64) error
	DeleteRoomRestriction(ctx context.Context, id int64) error
	DeleteUser(ctx context.Context, id int64) error
	GetLastRoomRestriction(ctx context.Context, roomID int64) (RoomRestriction, error)
	GetReservation(ctx context.Context, id int64) (Reservation, error)
	GetReservationByCode(ctx context.Context, code string) (Reservation, error)
	GetReservationByLastName(ctx context.Context, arg GetReservationByLastNameParams) (Reservation, error)
	GetRoom(ctx context.Context, id int64) (Room, error)
	GetRoomRestriction(ctx context.Context, id int64) (RoomRestriction, error)
	GetUser(ctx context.Context, id int64) (User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	ListAvailableRooms(ctx context.Context, arg ListAvailableRoomsParams) ([]Room, error)
	ListReservations(ctx context.Context, arg ListReservationsParams) ([]Reservation, error)
	ListReservationsByRoom(ctx context.Context, arg ListReservationsByRoomParams) ([]Reservation, error)
	ListRoomRestrictions(ctx context.Context, arg ListRoomRestrictionsParams) ([]RoomRestriction, error)
	ListRooms(ctx context.Context, arg ListRoomsParams) ([]Room, error)
	ListUsers(ctx context.Context, arg ListUsersParams) ([]User, error)
	UpdateReservation(ctx context.Context, arg UpdateReservationParams) error
	UpdateRoom(ctx context.Context, arg UpdateRoomParams) error
	UpdateRoomRestriction(ctx context.Context, arg UpdateRoomRestrictionParams) error
	UpdateUser(ctx context.Context, arg UpdateUserParams) error
}

var _ Querier = (*Queries)(nil)
