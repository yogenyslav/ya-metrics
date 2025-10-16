package service

import (
	models "github.com/yogenyslav/ya-metrics/internal/model"
)

type metricRepo[T int64 | float64] interface {
	Get(name string) (*models.Metrics[T], bool)
}

type GaugeRepo interface {
	metricRepo[float64]
	Set(name string, value float64)
}

type CounterRepo interface {
	metricRepo[int64]
	Update(name string, delta int64)
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
