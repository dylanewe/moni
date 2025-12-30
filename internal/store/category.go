package store

import (
	"context"
	"database/sql"
)

type Category struct {
	ID   int64
	Name string
}

type CategoryStore struct {
	db *sql.DB
}

func (s *CategoryStore) Insert(context.Context, *Category) error

func (s *CategoryStore) GetAll(context.Context) ([]Category, error)
