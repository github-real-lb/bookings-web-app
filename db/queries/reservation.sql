-- name: CreateReservation :one
INSERT INTO reservations (
  code, first_name, last_name, email, phone, start_date, end_date, room_id, notes
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING *;

-- name: DeleteAllReservations :exec
DELETE FROM reservations;

-- name: DeleteReservation :exec
DELETE FROM reservations
WHERE id = $1;

-- name: GetReservation :one
SELECT * FROM reservations
WHERE id = $1 LIMIT 1;

-- name: GetReservationByLastName :one
SELECT * FROM reservations
WHERE code = $1 AND last_name = $2 LIMIT 1;

-- name: ListReservations :many
SELECT * FROM reservations 
ORDER BY start_date, end_date ASC
LIMIT $1
OFFSET $2;

-- name: ListReservationsAndRooms :many
SELECT sqlc.embed(reservations), sqlc.embed(rooms) 
FROM reservations
LEFT JOIN rooms ON (reservations.room_id = rooms.id)
ORDER BY reservations.start_date, rooms.name ASC
LIMIT $1
OFFSET $2;


-- name: UpdateReservation :exec
UPDATE reservations
  set   code = $2,
        first_name = $3,
        last_name = $4, 
        email = $5,
        phone = $6, 
        start_date =  $7,
        end_date = $8,
        room_id = $9,
        notes = $10,
        updated_at = $11
WHERE id = $1;