package db

import (
	context "context"
	sql "database/sql"

	pg "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type (
	// Database interface to use across the app
	Database interface {
		QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
		QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
		ExecContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	}

	// Options to build db
	Options struct {
		URL string
	}
)

// New build db from options
func New(opt Options) (*gorm.DB, error) {
	db, err := gorm.Open(pg.Open(opt.URL), &gorm.Config{})
	return db, err
}
