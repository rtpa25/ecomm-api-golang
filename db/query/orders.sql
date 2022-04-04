-- name: AddOrder :one
INSERT INTO orders (
  amount,
  user_id,
  status,
  address,
  prodcut_id
) VALUES (
  $1, $2, $3, $4,$5
) RETURNING *;

-- name: DeleteOrder :exec
DELETE FROM orders WHERE id = $1;

-- name: UpdateOrder :one
UPDATE orders
SET amount = $2, status = $3, address = $4
WHERE id = $1
RETURNING *;

-- name: ListOrders :many
SELECT * FROM orders
ORDER BY id 
LIMIT $1
OFFSET $2; 

-- name: GetACertainOrder :one
SELECT * FROM orders
WHERE id = $1; 

-- name: GetOrderForCertainUser :one
SELECT * FROM orders
WHERE user_id = $1
ORDER BY id;

