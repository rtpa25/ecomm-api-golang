package db

import (
	"database/sql"
)

type Store interface {
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
