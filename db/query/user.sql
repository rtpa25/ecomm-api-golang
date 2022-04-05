-- name: CreateUser :one
INSERT INTO users (
  username,
  is_admin
) VALUES (
  $1, $2
) RETURNING *;

-- name: GetListOfUsers :many
SELECT * FROM users
ORDER BY id 
LIMIT $1
OFFSET $2;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;