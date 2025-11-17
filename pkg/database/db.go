package database

import "context"

// DB defines methods to operate with DB.
//
//go:generate mockgen -destination=../../tests/mocks/db.go -package=mocks . DB
type DB interface {
	Ping(ctx context.Context) error
}
