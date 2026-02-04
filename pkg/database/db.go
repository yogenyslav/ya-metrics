package database

import (
	"context"
	"database/sql"
	"errors"

	"github.com/yogenyslav/ya-metrics/pkg/errs"
)

// CtxKey is a type for context keys.
type CtxKey string

// TxKey is the context key for DB transaction.
const TxKey CtxKey = "tx"

// ErrNoTx is an error when no transaction is found in context.
var ErrNoTx = errors.New("no transaction in context")

// DB defines methods to operate with DB.
//
//go:generate mockgen -destination=../../tests/mocks/db.go -package=mocks . DB,TxDB
type DB interface {
	Exec(ctx context.Context, query string, args ...any) (int64, error)
	QueryRow(ctx context.Context, dst any, query string, args ...any) error
	QuerySlice(ctx context.Context, dst any, query string, args ...any) error
	Ping(ctx context.Context) error
	SQLDB() (*sql.DB, error)
	Close()
}

// TxDB defines methods to operate with DB transactions.
type TxDB interface {
	DB

	beginTx(ctx context.Context) (context.Context, error)
	commitTx(ctx context.Context) error
	rollbackTx(ctx context.Context) error
}

// UnitOfWork provides methods to execute operations within a transaction.
//
//go:generate mockgen -destination=../../tests/mocks/uow.go -package=mocks . UnitOfWork
type UnitOfWork interface {
	// WithTx executes the given function within a database transaction.
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error
}

type unitOfWork struct {
	db TxDB
}

// NewUnitOfWork creates a new UnitOfWork instance.
func NewUnitOfWork(db TxDB) *unitOfWork {
	return &unitOfWork{db: db}
}

// WithTx implements UnitOfWork method.
func (uow *unitOfWork) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	if uow.db == nil {
		return fn(ctx)
	}

	tx, err := uow.db.beginTx(ctx)
	if err != nil {
		return errs.Wrap(err, "begin transaction")
	}

	defer func() {
		if e := recover(); e != nil {
			uow.db.rollbackTx(tx) //nolint:errcheck // nothing we can do
			panic(e)
		}

		if err != nil {
			uow.db.rollbackTx(tx) //nolint:errcheck // nothing we can do
		} else {
			err = uow.db.commitTx(tx)
		}
	}()

	err = fn(tx)
	return err
}
