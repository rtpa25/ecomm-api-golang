// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0
// source: size_product_map.sql

package db

import (
	"context"
)

const addSizeToProduct = `-- name: AddSizeToProduct :one
INSERT INTO
  size_product_map (product_id, size_id)
VALUES
  ($1, $2) RETURNING id, product_id, size_id
`

type AddSizeToProductParams struct {
	ProductID int32 `json:"product_id"`
	SizeID    int32 `json:"size_id"`
}

func (q *Queries) AddSizeToProduct(ctx context.Context, arg AddSizeToProductParams) (SizeProductMap, error) {
	row := q.db.QueryRowContext(ctx, addSizeToProduct, arg.ProductID, arg.SizeID)
	var i SizeProductMap
	err := row.Scan(&i.ID, &i.ProductID, &i.SizeID)
	return i, err
}

const getASingleSizeProductMapRow = `-- name: GetASingleSizeProductMapRow :one
SELECT
  id, product_id, size_id
FROM
  size_product_map
WHERE
  size_id = $1
  AND product_id = $2
`

type GetASingleSizeProductMapRowParams struct {
	SizeID    int32 `json:"size_id"`
	ProductID int32 `json:"product_id"`
}

func (q *Queries) GetASingleSizeProductMapRow(ctx context.Context, arg GetASingleSizeProductMapRowParams) (SizeProductMap, error) {
	row := q.db.QueryRowContext(ctx, getASingleSizeProductMapRow, arg.SizeID, arg.ProductID)
	var i SizeProductMap
	err := row.Scan(&i.ID, &i.ProductID, &i.SizeID)
	return i, err
}

const getAvailableSizesOfAProduct = `-- name: GetAvailableSizesOfAProduct :many
SELECT
  name
FROM
  sizes
  JOIN size_product_map ON size_product_map.size_id = sizes.id
WHERE
  product_id = $1
`

func (q *Queries) GetAvailableSizesOfAProduct(ctx context.Context, productID int32) ([]string, error) {
	rows, err := q.db.QueryContext(ctx, getAvailableSizesOfAProduct, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []string{}
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		items = append(items, name)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listProductsOfSize = `-- name: ListProductsOfSize :many
SELECT
  id, product_id, size_id
FROM
  size_product_map
WHERE
  size_id = $1
ORDER BY
  id
LIMIT
  $2 OFFSET $3
`

type ListProductsOfSizeParams struct {
	SizeID int32 `json:"size_id"`
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListProductsOfSize(ctx context.Context, arg ListProductsOfSizeParams) ([]SizeProductMap, error) {
	rows, err := q.db.QueryContext(ctx, listProductsOfSize, arg.SizeID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []SizeProductMap{}
	for rows.Next() {
		var i SizeProductMap
		if err := rows.Scan(&i.ID, &i.ProductID, &i.SizeID); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateAvailableSizes = `-- name: UpdateAvailableSizes :one
UPDATE
  size_product_map
SET
  size_id = $2
WHERE
  id = $1 RETURNING id, product_id, size_id
`

type UpdateAvailableSizesParams struct {
	ID     int32 `json:"id"`
	SizeID int32 `json:"size_id"`
}

func (q *Queries) UpdateAvailableSizes(ctx context.Context, arg UpdateAvailableSizesParams) (SizeProductMap, error) {
	row := q.db.QueryRowContext(ctx, updateAvailableSizes, arg.ID, arg.SizeID)
	var i SizeProductMap
	err := row.Scan(&i.ID, &i.ProductID, &i.SizeID)
	return i, err
}
