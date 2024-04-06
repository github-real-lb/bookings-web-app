-- name: CreateRoomRestriction :one
INSERT INTO room_restrictions (
  start_date, end_date, room_id, reservation_id, restriction
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: DeleteRoomRestriction :exec
DELETE FROM room_restrictions
WHERE id = $1;

-- name: GetLastRoomRestriction :one
SELECT * FROM room_restrictions
WHERE room_id = $1 
ORDER BY created_at DESC
LIMIT 1;

-- name: GetRoomRestriction :one
SELECT * FROM room_restrictions
WHERE id = $1 LIMIT 1;

-- name: ListRoomRestrictions :many
SELECT * FROM room_restrictions
ORDER BY room_id, start_date
LIMIT $1
OFFSET $2;

-- name: UpdateRoomRestriction :exec
UPDATE room_restrictions
  set   start_date = $2,
        end_date = $3, 
        room_id = $4,
        reservation_id = $5, 
        restriction =  $6,
        updated_at = $7
WHERE id = $1;