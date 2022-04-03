-- name: AddCategory :one
INSERT INTO categories (
  name
) VALUES (
  $1
) RETURNING *;

-- name: DeleteCategory :exec
DELETE FROM categories WHERE id = $1;