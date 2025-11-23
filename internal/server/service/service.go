package service

import (
	"context"

	"github.com/yogenyslav/ya-metrics/internal/model"
)

type metricRepo[T int64 | float64] interface {
	Get(ctx context.Context, metricID, metricType string) (*model.Metrics[T], error)
	List(ctx context.Context) ([]model.Metrics[T], error)
}

// GaugeRepo is the interface for gauge metric repository.
type GaugeRepo interface {
	metricRepo[float64]
	Set(ctx context.Context, m *model.Metrics[float64]) error
	SetBatch(ctx context.Context, ms []*model.Metrics[float64]) error
}

// CounterRepo is the interface for counter metric repository.
type CounterRepo interface {
	metricRepo[int64]
	Update(ctx context.Context, m *model.Metrics[int64]) error
	UpdateBatch(ctx context.Context, ms []*model.Metrics[int64]) error
}

// Service provides metric-related operations.
type Service struct {
	gr GaugeRepo
	cr CounterRepo
}

// NewService creates a new Service instance.
func NewService(gr GaugeRepo, cr CounterRepo) *Service {
	return &Service{
		gr: gr,
		cr: cr,
	}
}
