-- name: AddSize :one
INSERT INTO
  sizes (name)
VALUES
  ($1) RETURNING *;

-- name: ListAllSizes :many
SELECT
  *
FROM
  sizes;

-- name: DeleteSize :exec
DELETE FROM
  sizes
WHERE
  id = $1;