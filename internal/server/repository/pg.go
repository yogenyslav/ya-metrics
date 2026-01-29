package repository

import (
	"context"
	"reflect"

	"github.com/jackc/pgx/v5"
	"github.com/yogenyslav/ya-metrics/internal/model"
	"github.com/yogenyslav/ya-metrics/pkg/database"
	"github.com/yogenyslav/ya-metrics/pkg/errs"
)

// MetricPostgresRepo is a metric repository for pg.
type MetricPostgresRepo[T int64 | float64] struct {
	pg database.DB
}

// NewMetricPostgresRepo creates a new MetricPostgresRepo.
func NewMetricPostgresRepo[T int64 | float64](pg database.DB) *MetricPostgresRepo[T] {
	return &MetricPostgresRepo[T]{pg: pg}
}

const getGaugeMetric = `
	select id, mtype, coalesce(delta::double precision, value::double precision) as value
	from metrics
	where id = $1 and mtype = $2;
`

const getCounterMetric = `
	select id, mtype, coalesce(delta::double precision, value::double precision)::bigint as value
	from metrics
	where id = $1 and mtype = $2;
`

// Get retrieves a metric by name.
func (r *MetricPostgresRepo[T]) Get(ctx context.Context, metricName, metricType string) (*model.Metrics[T], error) {
	var (
		m     model.Metrics[T]
		query string
	)

	if reflect.TypeFor[T]().Kind() == reflect.Float64 {
		query = getGaugeMetric
	} else {
		query = getCounterMetric
	}

	err := r.pg.QueryRow(ctx, &m, query, metricName, metricType)
	if err != nil {
		return nil, errs.Wrap(err, "failed to query")
	}
	return &m, nil
}

const listGaugeMetrics = `
	select id, mtype, coalesce(delta::double precision, value::double precision) as value
	from metrics;
`

const listCounterMetrics = `
	select id, mtype, coalesce(delta::double precision, value::double precision)::bigint as value
	from metrics;
`

// List retrieves all metrics.
func (r *MetricPostgresRepo[T]) List(ctx context.Context) ([]model.Metrics[T], error) {
	var (
		metrics []model.Metrics[T]
		query   string
	)

	if reflect.TypeFor[T]().Kind() == reflect.Float64 {
		query = listGaugeMetrics
	} else {
		query = listCounterMetrics
	}

	err := r.pg.QuerySlice(ctx, &metrics, query)
	if err != nil {
		return nil, errs.Wrap(err, "failed to query")
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

const updateCounterMetric = `
	insert into metrics (id, mtype, delta)
	values ($1, $2, $3)
	on conflict (id) do update set 
		delta = metrics.delta + $3,
		updated_at = current_timestamp;
`

func (r *MetricPostgresRepo[T]) setOrUpdate(ctx context.Context, m *model.Metrics[T]) error {
	var query string

	switch any(m.Value).(type) {
	case float64:
		query = setGaugeMetric
	case int64:
		query = updateCounterMetric
	default:
		return errs.Wrap(errs.ErrInvalidMetricType, "unsupported metric type")
	}

	rowsCount, err := r.pg.Exec(ctx, query, m.ID, m.Type, m.Value)
	if err != nil {
		return errs.Wrap(err, "failed to exec")
	}

	if rowsCount == 0 {
		return errs.Wrap(pgx.ErrNoRows, "no metric found")
	}

	return nil
}

// Set sets the value of a gauge metric.
func (r *MetricPostgresRepo[T]) Set(ctx context.Context, m *model.Metrics[float64]) error {
	return r.setOrUpdate(ctx, &model.Metrics[T]{
		ID:    m.ID,
		Type:  m.Type,
		Value: T(m.Value),
	})
}

// Update updates the value of a counter metric.
func (r *MetricPostgresRepo[T]) Update(ctx context.Context, m *model.Metrics[int64]) error {
	return r.setOrUpdate(ctx, &model.Metrics[T]{
		ID:    m.ID,
		Type:  m.Type,
		Value: T(m.Value),
	})
}
