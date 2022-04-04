// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0
// source: category.sql

package db

import (
	"context"
)

const addCategory = `-- name: AddCategory :one
INSERT INTO categories (
  name
) VALUES (
  $1
) RETURNING id, name
`

func (q *Queries) AddCategory(ctx context.Context, name string) (Category, error) {
	row := q.db.QueryRowContext(ctx, addCategory, name)
	var i Category
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const deleteCategory = `-- name: DeleteCategory :exec
DELETE FROM categories WHERE id = $1
`

func (q *Queries) DeleteCategory(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, deleteCategory, id)
	return err
}

const listAllCategory = `-- name: ListAllCategory :many
SELECT id, name FROM categories
`

func (q *Queries) ListAllCategory(ctx context.Context) ([]Category, error) {
	rows, err := q.db.QueryContext(ctx, listAllCategory)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Category{}
	for rows.Next() {
		var i Category
		if err := rows.Scan(&i.ID, &i.Name); err != nil {
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
