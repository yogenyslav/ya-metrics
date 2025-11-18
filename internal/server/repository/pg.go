package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/yogenyslav/ya-metrics/internal/model"
	"github.com/yogenyslav/ya-metrics/pkg/database"
)

// MetricPostgresRepo is a metric repository for pg.
type MetricPostgresRepo[T int64 | float64] struct {
	pg database.DB
}

// NewMetricPostgresRepo creates a new MetricPostgresRepo.
func NewMetricPostgresRepo[T int64 | float64](pg database.DB) *MetricPostgresRepo[T] {
	return &MetricPostgresRepo[T]{pg: pg}
}

const getMetric = `
	select id, mtype, delta, value
	from metrics
	where id = $1 and mtype = $2;
`

// Get retrieves a metric by name.
func (r *MetricPostgresRepo[T]) Get(ctx context.Context, name string) (*model.Metrics[T], error) {
	var m model.MetricsDto
	err := r.pg.QueryRow(ctx, &m, getMetric, name, any(m.Type))
	if err != nil {
		return nil, err
	}

	return &model.Metrics[T]{
		ID:    m.ID,
		Type:  m.Type,
		Value: T(*m.Value),
	}, nil
}

const listMetrics = `
	select id, mtype, delta, value
	from metrics;
`

// List retrieves all metrics.
func (r *MetricPostgresRepo[T]) List(ctx context.Context) ([]model.Metrics[T], error) {
	var metrics []model.Metrics[T]
	err := r.pg.QuerySlice(ctx, &metrics, listMetrics)
	if err != nil {
		return nil, err
	}
	return metrics, nil
}

const setGaugeMetric = `
	insert into metrics (id, mtype, value)
	values ($1, $2, $3)
	on conflict (id) do update set 
		value = $3,
		updated_at = current_timestamp;
`

// Set sets the value of a gauge metric.
func (r *MetricPostgresRepo[T]) Set(ctx context.Context, name string, value float64, tp string) error {
	rowsCount, err := r.pg.Exec(ctx, setGaugeMetric, name, tp, value)
	if err != nil {
		return err
	}
	if rowsCount == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

const updateCounterMetric = `
	insert into metrics (id, mtype, delta)
	values ($1, $2, $3)
	on conflict (id) do update set 
		delta = metrics.delta + $3,
		updated_at = current_timestamp;
`

// Update updates the value of a counter metric.
func (r *MetricPostgresRepo[T]) Update(ctx context.Context, name string, delta int64, tp string) error {
	rowsCount, err := r.pg.Exec(ctx, updateCounterMetric, name, tp, delta)
	if err != nil {
		return err
	}
	if rowsCount == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
