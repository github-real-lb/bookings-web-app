-- name: CreateRoom :one
INSERT INTO rooms (
  name
) VALUES (
  $1
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
WHERE (start_date < '2024-02-05' AND end_date > '2024-02-01')
);

-- name: ListRooms :many
SELECT * FROM rooms
ORDER BY name
LIMIT $1
OFFSET $2;

-- name: UpdateRoom :exec
UPDATE rooms
  set   name = $2,
        updated_at = $3
WHERE id = $1;