package postgres

import (
	"context"
	"database/sql"

	// use postgres
	_ "github.com/lib/pq"
)

type (
	// DB postgres
	DB struct {
		PG *sql.DB
	}
)

// QueryContext runs query on the db
func (d *DB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return d.PG.QueryContext(ctx, query, args...)
}

// ExecContext executes query on the db
func (d *DB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return d.PG.ExecContext(ctx, query, args...)
}

// QueryRowContext run row query on the db
func (d *DB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return d.PG.QueryRowContext(ctx, query, args...)
}

// Connect opens connection to db
func Connect(url string) (*sql.DB, error) {
	db, err := sql.Open("postgres", url)
	return db, err
}
