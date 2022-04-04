-- name: AddCategoryToProduct :one
INSERT INTO category_product_map (
  product_id,
  category_id
) VALUES (
  $1, $2
) RETURNING *;

-- name: ListProductsOfCategory :many
SELECT * FROM category_product_map 
WHERE category_id = $1
ORDER BY id 
LIMIT $2
OFFSET $3;


