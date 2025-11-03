package repository

import "github.com/yogenyslav/ya-metrics/internal/model"

// MetricInMemRepo is an in-memory repository for metrics.
type MetricInMemRepo[T int64 | float64] struct {
	storage map[string]*model.Metrics[T]
}

// NewMetricInMemRepo creates a new instance of MetricInMemRepo.
func NewMetricInMemRepo[T int64 | float64]() *MetricInMemRepo[T] {
	return &MetricInMemRepo[T]{
		storage: make(map[string]*model.Metrics[T]),
	}
}

// Get returns the value of a metric by its name and a bool flag to check if it exists.
func (r *MetricInMemRepo[T]) Get(name string) (*model.Metrics[T], bool) {
	value, exists := r.storage[name]
	return value, exists
}

// Set sets the value of a metric by its name.
func (r *MetricInMemRepo[T]) Set(name string, value T, tp string) {
	if metric, ok := r.storage[name]; ok {
		metric.Value = value
		return
	}
	r.storage[name] = &model.Metrics[T]{ID: name, Type: tp, Value: value}
}

// Update updates the value of a metric by adding the delta to the current value.
func (r *MetricInMemRepo[T]) Update(name string, delta T, tp string) {
	if metric, exists := r.storage[name]; exists {
		metric.Value += delta
	} else {
		r.storage[name] = &model.Metrics[T]{ID: name, Type: tp, Value: delta}
	}
}

// List returns a list of all metrics in the repository.
func (r *MetricInMemRepo[T]) List() []model.Metrics[T] {
	metrics := make([]model.Metrics[T], 0, len(r.storage))
	for _, metric := range r.storage {
		metrics = append(metrics, *metric)
	}
	return metrics
}
