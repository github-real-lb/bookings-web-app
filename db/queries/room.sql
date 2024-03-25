-- name: CreateRoom :one
INSERT INTO rooms (
  name
) VALUES (
  $1
)
RETURNING *;

-- name: GetRoom :one
SELECT * FROM rooms
WHERE id = $1 LIMIT 1;

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

-- name: DeleteRoom :exec
DELETE FROM rooms
WHERE id = $1;