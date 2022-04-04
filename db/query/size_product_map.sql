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

-- name: GetAvailableSizesOfAProduct :many
SELECT name 
FROM sizes
JOIN size_product_map ON size_product_map.size_id = sizes.id 
WHERE product_id = $1;

-- name: UpdateAvailableSizes :one
UPDATE size_product_map
SET size_id = $2
WHERE id = $1
RETURNING *;

-- name: GetASingleSizeProductMapRow :one
SELECT * FROM size_product_map
WHERE size_id = $1 AND product_id = $2;

