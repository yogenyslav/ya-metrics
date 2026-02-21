package repository

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"sync"

	"github.com/yogenyslav/ya-metrics/internal/model"
	"github.com/yogenyslav/ya-metrics/pkg/errs"
)

// StorageState represents the in-memory storage state for metrics.
type StorageState[T int64 | float64] map[string]*model.Metrics[T]

// RestoreMetrics restores metrics from the file.
func RestoreMetrics(
	filePath string,
) (gaugeMetrics StorageState[float64], countMetrics StorageState[int64], err error) {
	f, err := os.OpenFile(filePath, os.O_RDONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, nil, errs.Wrap(err, "open file for restore")
	}
	defer f.Close()

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, nil, errs.Wrap(err, "read data from file")
	}

	if len(data) == 0 {
		return nil, nil, nil
	}

	var v []*model.MetricsDto
	err = json.Unmarshal(data, &v)
	if err != nil {
		return nil, nil, errs.Wrap(err, "unmarshal data")
	}

	gaugeMetrics = make(StorageState[float64])
	countMetrics = make(StorageState[int64])
	for _, m := range v {
		switch m.Type {
		case model.Gauge:
			gaugeMetrics[m.ID] = m.ToGaugeMetric()
		case model.Counter:
			countMetrics[m.ID] = m.ToCounterMetric()
		}
	}

	return gaugeMetrics, countMetrics, nil
}

// MetricInMemRepo is an in-memory repository for metrics.
type MetricInMemRepo[T int64 | float64] struct {
	storage StorageState[T]
	mu      *sync.RWMutex
}

// NewMetricInMemRepo creates a new instance of MetricInMemRepo.
func NewMetricInMemRepo[T int64 | float64](state StorageState[T]) *MetricInMemRepo[T] {
	var storage StorageState[T]
	if state != nil {
		storage = state
	} else {
		storage = make(StorageState[T])
	}

	return &MetricInMemRepo[T]{
		storage: storage,
		mu:      &sync.RWMutex{},
	}
}

// GetMetrics returns all metrics in MetricsDto format.
func (r *MetricInMemRepo[T]) GetMetrics(_ context.Context) ([]*model.MetricsDto, error) {
	r.mu.RLock()
	metrics := make([]*model.MetricsDto, 0, len(r.storage))
	for _, metric := range r.storage {
		metrics = append(metrics, metric.ToDto())
	}
	r.mu.RUnlock()
	return metrics, nil
}

// Get returns the value of a metric by its name and a bool flag to check if it exists.
func (r *MetricInMemRepo[T]) Get(_ context.Context, metricName, _ string) (*model.Metrics[T], error) {
	r.mu.RLock()
	value, exists := r.storage[metricName]
	r.mu.RUnlock()
	if !exists {
		return nil, errors.New("value not found")
	}
	return value, nil
}

// Set sets the value of a metric by its name.
func (r *MetricInMemRepo[T]) Set(_ context.Context, m *model.Metrics[T]) error {
	r.mu.Lock()
	if metric, ok := r.storage[m.ID]; ok {
		metric.Value = m.Value
	} else {
		r.storage[m.ID] = m
	}
	r.mu.Unlock()
	return nil
}

// Update updates the value of a metric by adding the delta to the current value.
func (r *MetricInMemRepo[T]) Update(_ context.Context, m *model.Metrics[T]) error {
	r.mu.Lock()
	if metric, exists := r.storage[m.ID]; exists {
		metric.Value += m.Value
	} else {
		r.storage[m.ID] = m
	}
	r.mu.Unlock()
	return nil
}

// List returns a list of all metrics in the repository.
func (r *MetricInMemRepo[T]) List(_ context.Context) ([]model.Metrics[T], error) {
	r.mu.RLock()
	metrics := make([]model.Metrics[T], 0, len(r.storage))
	for _, metric := range r.storage {
		metrics = append(metrics, *metric)
	}
	r.mu.RUnlock()
	return metrics, nil
}
