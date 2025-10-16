package models

// Metric types
const (
	Counter = "counter"
	Gauge   = "gauge"
)

// Metrics represents a metric with its properties.
type Metrics[T int64 | float64] struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value T      `json:"value,omitempty"`
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
