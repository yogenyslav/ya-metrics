package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/yogenyslav/ya-metrics/internal/model"
	"github.com/yogenyslav/ya-metrics/pkg/database"
	"github.com/yogenyslav/ya-metrics/pkg/errs"
)

// MetricPostgresRepo is a metric repository for pg.
type MetricPostgresRepo[T int64 | float64] struct {
	pg database.PostgresTxDB
}

// NewMetricPostgresRepo creates a new MetricPostgresRepo.
func NewMetricPostgresRepo[T int64 | float64](pg database.PostgresTxDB) *MetricPostgresRepo[T] {
	return &MetricPostgresRepo[T]{pg: pg}
}

const getMetric = `
	select id, mtype, delta, value
	from metrics
	where id = $1 and mtype = $2;
`

// Get retrieves a metric by name.
func (r *MetricPostgresRepo[T]) Get(ctx context.Context, metricName, metricType string) (*model.Metrics[T], error) {
	var m model.MetricsDto
	err := r.pg.QueryRow(ctx, &m, getMetric, metricName, metricType)
	if err != nil {
		return nil, errs.Wrap(err, "failed to query")
	}

	res := &model.Metrics[T]{
		ID:   m.ID,
		Type: m.Type,
	}

	if m.Type == model.Gauge && m.Value != nil {
		res.Value = T(*m.Value)
	} else if m.Type == model.Counter && m.Delta != nil {
		res.Value = T(*m.Delta)
	}

	return res, nil
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

func (r *MetricPostgresRepo[T]) setOrUpdateBatch(ctx context.Context, ms []*model.Metrics[T]) error {
	tx, err := r.pg.BeginTx(ctx)
	if err != nil {
		return errs.Wrap(err, "failed to begin tx")
	}
	defer tx.Rollback(ctx)

	for _, m := range ms {
		var query string

		switch any(m.Value).(type) {
		case float64:
			query = setGaugeMetric
		case int64:
			query = updateCounterMetric
		default:
			return errs.Wrap(errs.ErrInvalidMetricType, "unsupported metric type")
		}

		tag, err := tx.Exec(ctx, query, m.ID, m.Type, m.Value)
		if err != nil {
			return errs.Wrap(err, "failed to exec")
		}

		if tag.RowsAffected() == 0 {
			return errs.Wrap(pgx.ErrNoRows, "no metric found")
		}
	}

	return tx.Commit(ctx)
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

// UpdateBatch updates a batch of counter metrics.
func (r *MetricPostgresRepo[T]) UpdateBatch(ctx context.Context, ms []*model.Metrics[int64]) error {
	return r.setOrUpdateBatch(ctx, func() []*model.Metrics[T] {
		res := make([]*model.Metrics[T], 0, len(ms))
		for _, m := range ms {
			res = append(res, &model.Metrics[T]{
				ID:    m.ID,
				Type:  m.Type,
				Value: T(m.Value),
			})
		}
		return res
	}())
}

// SetBatch sets a batch of gauge metrics.
func (r *MetricPostgresRepo[T]) SetBatch(ctx context.Context, ms []*model.Metrics[float64]) error {
	return r.setOrUpdateBatch(ctx, func() []*model.Metrics[T] {
		res := make([]*model.Metrics[T], 0, len(ms))
		for _, m := range ms {
			res = append(res, &model.Metrics[T]{
				ID:    m.ID,
				Type:  m.Type,
				Value: T(m.Value),
			})
		}
		return res
	}())
}
