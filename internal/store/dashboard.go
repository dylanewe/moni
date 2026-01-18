package store

import (
	"context"
	"database/sql"
	"fmt"
)

// DashboardStore provides methods for accessing dashboard data.
type DashboardStore struct {
	db *sql.DB
}

// GetTotalIncomeAndExpense returns the total income and expense from all transactions.
func (s *DashboardStore) GetTotalIncomeAndExpense() (float64, float64, error) {
	query := `
		SELECT
			COALESCE(SUM(CASE WHEN amount > 0 THEN amount ELSE 0 END), 0) AS total_income,
			COALESCE(SUM(CASE WHEN amount < 0 THEN amount ELSE 0 END), 0) AS total_expense
		FROM transactions;
	`
	rows, err := s.db.QueryContext(context.Background(), query)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to query total income and expense: %w", err)
	}
	defer rows.Close()

	var totalIncome, totalExpense float64
	if rows.Next() {
		if err := rows.Scan(&totalIncome, &totalExpense); err != nil {
			return 0, 0, fmt.Errorf("failed to scan total income and expense: %w", err)
		}
	}

	return totalIncome, totalExpense, nil
}

// GetMonthlyIncomeAndExpense returns the income and expense for a specific month and year.
func (s *DashboardStore) GetMonthlyIncomeAndExpense(year int, month int) (float64, float64, error) {
	query := `
		SELECT
			COALESCE(SUM(CASE WHEN amount > 0 THEN amount ELSE 0 END), 0) AS monthly_income,
			COALESCE(SUM(CASE WHEN amount < 0 THEN amount ELSE 0 END), 0) AS monthly_expense
		FROM transactions
		WHERE EXTRACT(YEAR FROM date) = $1 AND EXTRACT(MONTH FROM date) = $2;
	`
	rows, err := s.db.QueryContext(context.Background(), query, year, month)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to query monthly income and expense: %w", err)
	}
	defer rows.Close()

	var monthlyIncome, monthlyExpense float64
	if rows.Next() {
		if err := rows.Scan(&monthlyIncome, &monthlyExpense); err != nil {
			return 0, 0, fmt.Errorf("failed to scan monthly income and expense: %w", err)
		}
	}

	return monthlyIncome, monthlyExpense, nil
}
