-- name: AddCategory :one
INSERT INTO
  categories (name)
VALUES
  ($1) RETURNING *;

-- name: ListAllCategory :many
SELECT
  *
FROM
  categories;

-- name: DeleteCategory :exec
DELETE FROM
  categories
WHERE
  id = $1;