package store

import (
	"context"
	"database/sql"
)

type Transaction struct {
	ID          int64
	Description string
	Category    Category
	Amount      float64
	Date        string
}

type TransactionStore struct {
	db *sql.DB
}

func (s *TransactionStore) Insert(ctx context.Context, tx *Transaction) error {
	query := `
		INSERT INTO transactions (id, description, category, amount, date)
		VALUES ($1, $2, $3, $4, $5)
	`

	ctx, cancel := context.WithTimeout(ctx, 10)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, tx.ID, tx.Description, tx.Category, tx.Amount, tx.Date)
	if err != nil {
		return err
	}

	return nil
}

func (s *TransactionStore) GetMonthlyIncome(ctx context.Context, startDate, endDate string) (float64, error)

func (s *TransactionStore) GetMonthlyExpense(ctx context.Context, startDate, endDate string) (float64, error)
