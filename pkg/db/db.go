package db

import (
	"context"
	"database/sql"

	"github.com/peteqproj/peteq/pkg/db/postgres"
)

type (
	// Database interface to use across the app
	Database interface {
		QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
		QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
		ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	}

	// Options to build db
	Options struct {
		DB *sql.DB
	}
)

// New build db from options
func New(opt Options) Database {
	return &postgres.DB{
		PG: opt.DB,
	}
}
