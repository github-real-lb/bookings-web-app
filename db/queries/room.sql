-- name: CheckRoomAvailabilty :one
SELECT count(*) = 0 as availabe
FROM room_restrictions
WHERE room_id = $1 AND (start_date < @end_date::date AND end_date > @start_date::date);

-- name: CreateRoom :one
INSERT INTO rooms (
  name, description, image_filename
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: DeleteRoom :exec
DELETE FROM rooms
WHERE id = $1;

-- name: GetRoom :one
SELECT * FROM rooms
WHERE id = $1 LIMIT 1;

-- name: ListAvailableRooms :many
SELECT *
FROM rooms
WHERE id NOT IN (
SELECT room_id
FROM room_restrictions
WHERE (start_date < @end_date::date AND end_date > @start_date::date)
);

-- name: ListRooms :many
SELECT * FROM rooms
ORDER BY name
LIMIT $1
OFFSET $2;

-- name: UpdateRoom :exec
UPDATE rooms
  set   name = $2,
        description = $3,
        updated_at = $4
WHERE id = $1;