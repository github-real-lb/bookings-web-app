// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: room.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createRoom = `-- name: CreateRoom :one
INSERT INTO rooms (
  name
) VALUES (
  $1
)
RETURNING id, name, created_at, updated_at
`

func (q *Queries) CreateRoom(ctx context.Context, name string) (Room, error) {
	row := q.db.QueryRow(ctx, createRoom, name)
	var i Room
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteRoom = `-- name: DeleteRoom :exec
DELETE FROM rooms
WHERE id = $1
`

func (q *Queries) DeleteRoom(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deleteRoom, id)
	return err
}

const getRoom = `-- name: GetRoom :one
SELECT id, name, created_at, updated_at FROM rooms
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetRoom(ctx context.Context, id int64) (Room, error) {
	row := q.db.QueryRow(ctx, getRoom, id)
	var i Room
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listRooms = `-- name: ListRooms :many
SELECT id, name, created_at, updated_at FROM rooms
ORDER BY name
LIMIT $1
OFFSET $2
`

type ListRoomsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListRooms(ctx context.Context, arg ListRoomsParams) ([]Room, error) {
	rows, err := q.db.Query(ctx, listRooms, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Room{}
	for rows.Next() {
		var i Room
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateRoom = `-- name: UpdateRoom :exec
UPDATE rooms
  set   name = $2,
        updated_at = $3
WHERE id = $1
`

type UpdateRoomParams struct {
	ID        int64              `json:"id"`
	Name      string             `json:"name"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
}

func (q *Queries) UpdateRoom(ctx context.Context, arg UpdateRoomParams) error {
	_, err := q.db.Exec(ctx, updateRoom, arg.ID, arg.Name, arg.UpdatedAt)
	return err
}