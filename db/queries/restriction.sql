-- name: CreateRestriction :one
INSERT INTO restrictions (
  name
) VALUES (
  $1
)
RETURNING *;

-- name: GetRestriction :one
SELECT * FROM restrictions
WHERE id = $1 LIMIT 1;

-- name: ListRestrictions :many
SELECT * FROM restrictions
ORDER BY name
LIMIT $1
OFFSET $2;

-- name: UpdateRestriction :exec
UPDATE restrictions
  set   name = $2,
        updated_at = $3
WHERE id = $1;

-- name: DeleteRestriction :exec
DELETE FROM restrictions
WHERE id = $1;