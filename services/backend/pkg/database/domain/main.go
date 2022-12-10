package domain

import (
	"context"
	"database/sql"
)

// BoilerConfig provides a sqlboiler database configuration information.
type BoilerConfig interface {
	// GetDebug returns true if debug information should be outputted to the DebugWriter handler.
	GetDebug() bool
}

// Executor can perform SQL queries with or without a context.
type Executor interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row

	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

// TxBeginner can begin sql transactions.
type TxBeginner interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}
