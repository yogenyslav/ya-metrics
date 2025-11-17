package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Postgres is a wrapper around PostgreSQL using pgx connection pool.
type Postgres struct {
	pool *pgxpool.Pool
}

// NewPostgres creates a new pg instance.
func NewPostgres(ctx context.Context, dsn string) (*Postgres, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	return &Postgres{pool: pool}, nil
}

// Close the underlying pool.
func (p *Postgres) Close() {
	p.pool.Close()
}

// Ping the database.
func (p *Postgres) Ping(ctx context.Context) error {
	return p.pool.Ping(ctx)
}
