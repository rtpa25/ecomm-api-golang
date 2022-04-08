package db

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
)

type Store interface {
	AddProductTx(ctx context.Context, arg AddProductRequestParams) (AddProductReponseParams, error)
	Querier
}

// sqlstore provides all functions to execute all db queries and transactions
type SQLStore struct {
	db *sql.DB
	*Queries
}

// NewStore creates a new store
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transaction
func (store *SQLStore) execTxn(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)

	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx error: %v. rb error: %v", err, rbErr)
		} else {
			return err
		}
	}

	return tx.Commit()
}

type AddProductRequestParams struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	ImageUrl    string   `json:"image_url"`
	ImageID     string   `json:"image_id"`
	Price       string   `json:"price"`
	CategoryIDs []string `json:"category_ids"`
	SizeIDs     []string `json:"size_ids"`
}

type AddProductReponseParams struct {
	ID             int32    `json:"id"`
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	ImageUrl       string   `json:"image_url"`
	ImageID        string   `json:"image_id"`
	Price          string   `json:"price"`
	Categories     []string `json:"categories"`
	AvailableSizes []string `json:"available_sizes"`
}

func (store *SQLStore) AddProductTx(ctx context.Context, arg AddProductRequestParams) (AddProductReponseParams, error) {
	var result AddProductReponseParams
	err := store.execTxn(ctx, func(q *Queries) error {
		var err error
		product, err := q.AddProduct(ctx, AddProductParams{
			Name:        arg.Name,
			Description: arg.Description,
			ImageUrl:    arg.ImageUrl,
			ImageID:     arg.ImageID,
			Price:       arg.Price,
		})

		result.ID = product.ID
		result.Name = product.Name
		result.Description = product.Description
		result.Price = product.Price
		result.ImageID = product.ImageID
		result.ImageUrl = product.ImageUrl

		//added to catagory_product_map
		for _, CategoryID := range arg.CategoryIDs {
			var intCategoryId int
			intCategoryId, _ = strconv.Atoi(CategoryID)
			_, err = q.AddCategoryToProduct(ctx, AddCategoryToProductParams{
				ProductID:  product.ID,
				CategoryID: int32(intCategoryId),
			})
		}

		//added to size_product_map
		for _, SizeID := range arg.SizeIDs {
			var intSizeId int
			intSizeId, _ = strconv.Atoi(SizeID)
			_, err = q.AddSizeToProduct(ctx, AddSizeToProductParams{
				ProductID: product.ID,
				SizeID:    int32(intSizeId),
			})
		}

		availableSizes, err := q.GetAvailableSizesOfAProduct(ctx, product.ID)

		categories, err := q.GetCategoriesOfAProduct(ctx, product.ID)

		result.AvailableSizes = availableSizes
		result.Categories = categories

		return err
	})
	return result, err
}
