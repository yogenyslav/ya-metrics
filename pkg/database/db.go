package database

import "context"

// DB defines methods to operate with DB.
//
//go:generate mockgen -destination=../../tests/mocks/db.go -package=mocks . DB
type DB interface {
	Exec(ctx context.Context, query string, args ...any) (int64, error)
	QueryRow(ctx context.Context, dsy any, query string, args ...any) error
	QuerySlice(ctx context.Context, dst any, query string, args ...any) error
	Ping(ctx context.Context) error
	Close()
}
