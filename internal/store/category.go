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

func (s *CategoryStore) Insert(ctx context.Context, cat *Category) error {
	query := `INSERT INTO categories (id, name) VALUES ($1, $2)`

	ctx, cancel := context.WithTimeout(ctx, 10)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, cat.ID, cat.Name)
	if err != nil {
		return err
	}

	return nil
}

func (s *CategoryStore) GetAll(ctx context.Context) ([]Category, error) {
	query := `SELECT id, name FROM categories`

	ctx, cancel := context.WithTimeout(ctx, 10)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var cat Category
		if err := rows.Scan(&cat.ID, &cat.Name); err != nil {
			return nil, err
		}

		categories = append(categories, cat)
	}

	return categories, rows.Err()
}

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
