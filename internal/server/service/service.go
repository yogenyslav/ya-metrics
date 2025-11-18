package service

import (
	"context"

	"github.com/yogenyslav/ya-metrics/internal/model"
)

type metricRepo[T int64 | float64] interface {
	Get(ctx context.Context, name string) (*model.Metrics[T], error)
	List(ctx context.Context) ([]model.Metrics[T], error)
}

// GaugeRepo is the interface for gauge metric repository.
type GaugeRepo interface {
	metricRepo[float64]
	Set(ctx context.Context, name string, value float64, tp string) error
}

// CounterRepo is the interface for counter metric repository.
type CounterRepo interface {
	metricRepo[int64]
	Update(ctx context.Context, name string, delta int64, tp string) error
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
