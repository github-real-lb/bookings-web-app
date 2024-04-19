-- name: CreateUser :one
INSERT INTO users (
  first_name, last_name, email, password, access_level
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY first_name, last_name
LIMIT $1
OFFSET $2;

-- name: UpdateUser :exec
UPDATE users
  set   first_name = $2,
        last_name = $3, 
        email = $4,
        password = $5, 
        access_level =  $6,
        updated_at = $7
WHERE id = $1;