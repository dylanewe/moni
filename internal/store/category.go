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

type categoryMap map[string]int64

func (s *CategoryStore) Insert(context.Context, *Category) error

func (s *CategoryStore) GetAll(context.Context) ([]Category, error)

func getCategoryMap(ctx context.Context, tx *sql.Tx) (categoryMap, error) {
	query := `SELECT id, name FROM categories`

	ctx, cancel := context.WithTimeout(ctx, 10)
	defer cancel()

	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categoryMap := make(categoryMap)
	for rows.Next() {
		var id int64
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			return nil, err
		}
		categoryMap[name] = id
	}

	return categoryMap, rows.Err()
}
