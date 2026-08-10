// Package postgres implements the repository and query-service interfaces on
// top of PostgreSQL with raw SQL. Repositories accept a DBTX (*sql.DB or
// *sql.Tx) so they can join application-level transactions; query services
// implement the read-model interfaces declared in internal/query.
package postgres

import (
	"context"
	"database/sql"
)

// DBTX is implemented by *sql.DB and *sql.Tx
type DBTX interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
}
