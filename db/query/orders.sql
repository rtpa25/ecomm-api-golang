-- name: AddOrder :one
INSERT INTO
  orders (
    quantity,
    user_id,
    address,
    product_id,
    selected_size
  )
VALUES
  ($1, $2, $3, $4, $5) RETURNING *;

-- name: GetOrdersForUser :many
SELECT
  *
FROM
  orders
WHERE
  user_id = $1
ORDER BY
  id;

-- name: GetSelfOrders :many
SELECT * FROM orders
JOIN products 
ON products.id=orders.product_id 
WHERE orders.user_id=$1
ORDER BY orders.id;

-- name: GetOrderById :one
SELECT * FROM orders
WHERE id=$1;

-- name: UpdateOrderForUser :one
UPDATE
  orders
SET
  quantity = $2,
  selected_size = $3,
  address = $4
WHERE
  id = $1 RETURNING *;

-- name: DeleteOrderById :exec
DELETE FROM
  orders
WHERE
  id = $1;