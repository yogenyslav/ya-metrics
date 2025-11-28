package database

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5"
)

// DB defines methods to operate with DB.
//
//go:generate mockgen -destination=../../tests/mocks/db.go -package=mocks . DB
type DB interface {
	Exec(ctx context.Context, query string, args ...any) (int64, error)
	QueryRow(ctx context.Context, dsy any, query string, args ...any) error
	QuerySlice(ctx context.Context, dst any, query string, args ...any) error
	Ping(ctx context.Context) error
	SQLDB() (*sql.DB, error)
	Close()
}

// PostgresTxDB defines transactional methods for pg.
//
//go:generate mockgen -destination=../../tests/mocks/pg_tx_db.go -package=mocks . PostgresTxDB
type PostgresTxDB interface {
	DB
	BeginTx(ctx context.Context) (pgx.Tx, error)
}
