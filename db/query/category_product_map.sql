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

-- name: GetCategoriesOfAProduct :many
SELECT name 
FROM categories
JOIN category_product_map ON category_product_map.category_id = categories.id 
WHERE product_id = $1;

-- name: UpdateCatagoriesOfACertainProduct :one
UPDATE category_product_map
SET category_id = $2
WHERE id = $1
RETURNING *;

-- name: GetASingleCategoryProductMapRow :one
SELECT * FROM category_product_map
WHERE category_id = $1 AND product_id = $2;