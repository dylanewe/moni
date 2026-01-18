package store

import (
	"context"
	"database/sql"
)

type Store struct {
	Transactions interface {
		Insert(context.Context, []Transaction) error
		GetIncomeByDate(ctx context.Context, startDate, endDate string) (float64, error)
		GetExpenseByDate(ctx context.Context, startDate, endDate string) (float64, error)
	}
	Categories interface {
		Insert(context.Context, *Category) error
		GetAll(context.Context) ([]Category, error)
	}
	Dashboard interface {
		GetTotalIncomeAndExpense() (float64, float64, error)
		GetMonthlyIncomeAndExpense(year int, month int) (float64, float64, error)
	}
}

func NewStore(db *sql.DB) Store {
	return Store{
		Transactions: &TransactionStore{db},
		Categories:   &CategoryStore{db},
		Dashboard:    &DashboardStore{db},
	}
}

func withTx(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
