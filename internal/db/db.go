package db

import (
	"context"
	"database/sql"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dylanewe/moni/internal/service"
	"github.com/dylanewe/moni/internal/store"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DBConnectionMsg struct {
	DB    *sql.DB
	Store *store.Store
	Err   error
}

func Init(addr string, categories []string) tea.Cmd {
	return func() tea.Msg {
		ctx := context.TODO()
		db, err := sql.Open("pgx", addr)
		if err != nil {
			return DBConnectionMsg{Err: err}
		}
		if err = db.PingContext(ctx); err != nil {
			return DBConnectionMsg{Err: fmt.Errorf("failed pinging db: %v", err)}
		}

		store := store.NewStore(db)

		return DBConnectionMsg{
			DB:    db,
			Store: &store,
			Err:   nil,
		}
	}
}

type ExtractStatementMsg struct {
	Transactions []store.Transaction
	Err          error
}

func ExtractStatement(llmParser service.StatementParser, cat []string, file string) tea.Cmd {
	return func() tea.Msg {
		tx, err := llmParser.ParseStatement(context.TODO(), cat, file)
		if err != nil {
			return ExtractStatementMsg{Err: fmt.Errorf("failed to extract transactions: %v", err)}
		}

		return ExtractStatementMsg{
			Transactions: tx,
			Err:          nil,
		}
	}
}

type AddStatementMsg struct {
	Err error
}

func AddStatement(txStore *store.Store, tx []store.Transaction) tea.Cmd {
	return func() tea.Msg {
		if err := txStore.Transactions.Insert(context.TODO(), tx); err != nil {
			return AddStatementMsg{Err: fmt.Errorf("failed to insert transactions: %v", err)}
		}
		return AddStatementMsg{Err: nil}
	}
}
