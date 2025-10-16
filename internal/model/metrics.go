package model

import "strconv"

// Metric types.
const (
	Counter = "counter"
	Gauge   = "gauge"
)

// Metrics represents a metric with its properties.
type Metrics[T int64 | float64] struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value T      `json:"value"`
}

// ToDto converts Metrics to MetricsDto.
func (m *Metrics[T]) ToDto() *MetricsDto {
	var valueStr string
	switch v := any(m.Value).(type) {
	case int64:
		valueStr = strconv.FormatInt(v, 10)
	case float64:
		valueStr = strconv.FormatFloat(v, 'f', -1, 64)
	}
	return &MetricsDto{
		Name:  m.Name,
		Type:  m.Type,
		Value: valueStr,
	}
}

// MetricsDto is a struct for transferring metric data.
type MetricsDto struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

// NewGaugeMetric creates a new gauge metric.
func NewGaugeMetric(name string) *Metrics[float64] {
	return &Metrics[float64]{
		Name: name,
		Type: Gauge,
	}
}

// NewCounterMetric creates a new counter metric.
func NewCounterMetric(name string) *Metrics[int64] {
	return &Metrics[int64]{
		Name: name,
		Type: Counter,
	}
}
