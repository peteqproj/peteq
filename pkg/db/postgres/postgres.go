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
		pg *sql.Conn
	}

	// ReadOptions to query db
	ReadOptions struct {
		Query string
	}
)

func (d *DB) Read(ctx context.Context, opt ReadOptions) (*sql.Rows, error) {
	return d.pg.QueryContext(ctx, opt.Query)
}
