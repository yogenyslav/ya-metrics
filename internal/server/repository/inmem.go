package repository

import (
	models "github.com/yogenyslav/ya-metrics/internal/model"
)

// MetricInMemRepo is an in-memory repository for metrics.
type MetricInMemRepo[T int64 | float64] struct {
	storage map[string]*models.Metrics[T]
}

// NewMetricInMemRepo creates a new instance of MetricInMemRepo.
func NewMetricInMemRepo[T int64 | float64]() *MetricInMemRepo[T] {
	return &MetricInMemRepo[T]{
		storage: make(map[string]*models.Metrics[T]),
	}
}

// Get returns the value of a metric by its name and a bool flag to check if it exists.
func (r *MetricInMemRepo[T]) Get(name string) (*models.Metrics[T], bool) {
	value, exists := r.storage[name]
	return value, exists
}

// Set sets the value of a metric by its name.
func (r *MetricInMemRepo[T]) Set(name string, value T) {
	r.storage[name] = &models.Metrics[T]{Name: name, Value: value}
}

// Update updates the value of a metric by adding the delta to the current value.
func (r *MetricInMemRepo[T]) Update(name string, delta T) {
	if metric, exists := r.storage[name]; exists {
		metric.Value += delta
	} else {
		r.storage[name] = &models.Metrics[T]{Name: name, Value: delta}
	}
}
