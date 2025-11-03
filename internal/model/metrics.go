package model

// Metric types.
const (
	Counter = "counter"
	Gauge   = "gauge"
)

// Metrics represents a metric with its properties.
type Metrics[T int64 | float64] struct {
	ID    string `json:"id"`
	Type  string `json:"type"`
	Value T      `json:"value"`
}

// NewGaugeMetric creates a new gauge metric.
func NewGaugeMetric(id string) *Metrics[float64] {
	return &Metrics[float64]{
		ID:   id,
		Type: Gauge,
	}
}

// NewCounterMetric creates a new counter metric.
func NewCounterMetric(id string) *Metrics[int64] {
	return &Metrics[int64]{
		ID:   id,
		Type: Counter,
	}
}

// ToDto converts Metrics to MetricsDto.
func (m *Metrics[T]) ToDto() *MetricsDto {
	metric := &MetricsDto{
		ID:   m.ID,
		Type: m.Type,
	}

	switch v := any(m.Value).(type) {
	case int64:
		metric.Delta = &v
	case float64:
		metric.Value = &v
	}

	return metric
}

// MetricsDto is a struct for transferring metric data.
type MetricsDto struct {
	ID    string   `json:"id"`
	Type  string   `json:"type"`
	Value *float64 `json:"value,omitempty"`
	Delta *int64   `json:"delta,omitempty"`
}
