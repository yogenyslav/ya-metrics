package database

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/yogenyslav/ya-metrics/pkg/retry"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Postgres is a wrapper around PostgreSQL using pgx connection pool.
type Postgres struct {
	pool     *pgxpool.Pool
	retryCfg *retry.Config
}

// NewPostgres creates a new pg instance.
func NewPostgres(ctx context.Context, dsn string, retryCfg *retry.Config) (*Postgres, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	return &Postgres{
		pool:     pool,
		retryCfg: retryCfg,
	}, nil
}

// SQLDB return a sql.DB format database conn.
func (p *Postgres) SQLDB() (*sql.DB, error) {
	db, err := sql.Open("pgx", p.pool.Config().ConnString())
	if err != nil {
		return nil, err
	}
	return db, nil
}

// Close the underlying pool.
func (p *Postgres) Close() {
	p.pool.Close()
}

// Ping the database.
func (p *Postgres) Ping(ctx context.Context) error {
	return retry.WithLinearBackoffRetry(p.retryCfg, func() error {
		err := p.pool.Ping(ctx)

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.ConnectionException, pgerrcode.ConnectionDoesNotExist, pgerrcode.ConnectionFailure,
				pgerrcode.SQLClientUnableToEstablishSQLConnection, pgerrcode.SQLServerRejectedEstablishmentOfSQLConnection,
				pgerrcode.TransactionResolutionUnknown, pgerrcode.ProtocolViolation:
				return err
			}
		}

		return retry.ErrUnretriable
	})
}

// Exec executes a DML query.
func (p *Postgres) Exec(ctx context.Context, query string, args ...any) (int64, error) {
	var (
		tag pgconn.CommandTag
		err error
	)

	err = retry.WithLinearBackoffRetry(p.retryCfg, func() error {
		tag, err = p.pool.Exec(ctx, query, args...)
		return err
	})
	if err != nil {
		return 0, err
	}

	return tag.RowsAffected(), nil
}

// QueryRow executes a DQL query that must return at most one row.
func (p *Postgres) QueryRow(ctx context.Context, dst any, query string, args ...any) error {
	return retry.WithLinearBackoffRetry(p.retryCfg, func() error {
		return pgxscan.Get(ctx, p.pool, dst, query, args...)
	})
}

// QuerySlice executes a DQL query that returns multiple rows.
func (p *Postgres) QuerySlice(ctx context.Context, dst any, query string, args ...any) error {
	return retry.WithLinearBackoffRetry(p.retryCfg, func() error {
		return pgxscan.Select(ctx, p.pool, dst, query, args...)
	})
}

// BeginTx starts a new transaction.
func (p *Postgres) BeginTx(ctx context.Context) (pgx.Tx, error) {
	var (
		err error
		tx  pgx.Tx
	)

	err = retry.WithLinearBackoffRetry(p.retryCfg, func() error {
		tx, err = p.pool.BeginTx(ctx, pgx.TxOptions{})
		return err
	})
	if err != nil {
		return nil, err
	}

	return tx, nil
}
