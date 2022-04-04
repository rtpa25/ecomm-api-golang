-- name: AddSizeToProduct :one
INSERT INTO size_product_map (
  product_id,
  size_id
) VALUES (
  $1, $2
) RETURNING *;

-- name: ListProductsOfSize :many
SELECT * FROM size_product_map 
WHERE size_id = $1
ORDER BY id 
LIMIT $2
OFFSET $3;


