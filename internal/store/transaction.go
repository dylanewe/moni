package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type Transaction struct {
	ID           int64   `json:"-"`
	Description  string  `json:"description"`
	CategoryName string  `json:"category"`
	Amount       float64 `json:"amount"`
	Date         string  `json:"date"`
}

type TransactionStore struct {
	db *sql.DB
}

func (s *TransactionStore) Insert(ctx context.Context, transactions []Transaction) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		categoryMap, err := getCategoryMap(ctx, tx)
		if err != nil {
			return err
		}

		valueStrings := make([]string, 0, len(transactions))
		valueArgs := make([]any, 0, len(transactions)*5)

		for i, t := range transactions {
			categoryID, exists := categoryMap[t.CategoryName]
			if !exists {
				return fmt.Errorf("category not found: %s", t.CategoryName)
			}

			valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d",
				i*5+1, i*5+2, i*5+3, i*5+4, i*5+5))
			valueArgs = append(valueArgs, t.ID, t.Description, categoryID, t.Amount, t.Date)
		}

		query := fmt.Sprintf("INSERT INTO transactions (id, description, category_id, amount, date) VALUES %s",
			strings.Join(valueStrings, ","))

		ctx, cancel := context.WithTimeout(ctx, 10)
		defer cancel()

		_, err = s.db.ExecContext(ctx, query, valueArgs...)
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *TransactionStore) GetIncomeByDate(ctx context.Context, startDate, endDate string) (float64, error) {
	var income sql.NullFloat64

	query := `
		SELECT COALESCE(SUM(amount), 0)
		FROM transactions
		WHERE date::date >= $1::date
			AND date::date <= $2::date
			AND amount >= 0
	`

	err := s.db.QueryRow(query, startDate, endDate).Scan(&income)
	if err != nil {
		return 0, fmt.Errorf("failed to get income: %w", err)
	}

	if !income.Valid {
		return 0, nil
	}

	return income.Float64, nil
}

func (s *TransactionStore) GetExpenseByDate(ctx context.Context, startDate, endDate string) (float64, error) {
	var expense sql.NullFloat64

	query := `
		SELECT COALESCE(SUM(amount), 0)
		FROM transactions
		WHERE date::date >= $1::date
			AND date::date <= $2::date
			AND amount < 0
	`

	err := s.db.QueryRow(query, startDate, endDate).Scan(&expense)
	if err != nil {
		return 0, fmt.Errorf("failed to get income: %w", err)
	}

	if !expense.Valid {
		return 0, nil
	}

	return expense.Float64, nil
}
