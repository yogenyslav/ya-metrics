package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yogenyslav/ya-metrics/internal/model"
	"github.com/yogenyslav/ya-metrics/pkg"
)

func TestService_GetMetric(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		metricType string
		metricName string
		wantMetric *model.MetricsDto
		wantExists bool
	}{
		{
			name:       "Get existing gauge metric",
			metricType: model.Gauge,
			metricName: "mem_alloc",
			wantMetric: &model.MetricsDto{Name: "mem_alloc", Type: model.Gauge, Value: pkg.Ptr(0.0)},
			wantExists: true,
		},
		{
			name:       "Get non-existing gauge metric",
			metricType: model.Gauge,
			metricName: "non_existing_gauge",
			wantMetric: nil,
			wantExists: false,
		},
		{
			name:       "Get existing counter metric",
			metricType: model.Counter,
			metricName: "request_count",
			wantMetric: &model.MetricsDto{Name: "request_count", Type: model.Counter, Delta: pkg.Ptr[int64](0)},
			wantExists: true,
		},
		{
			name:       "Get non-existing counter metric",
			metricType: model.Counter,
			metricName: "non_existing_counter",
			wantMetric: nil,
			wantExists: false,
		},
		{
			name:       "Get metric with invalid type",
			metricType: "invalid_type",
			metricName: "some_metric",
			wantMetric: nil,
			wantExists: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gr := &MockGaugeRepo{}
			cr := &MockCounterRepo{}

			if tt.wantMetric != nil {
				gr.On("Get", tt.metricName).Return(model.NewGaugeMetric(tt.wantMetric.Name), tt.wantExists)
				cr.On("Get", tt.metricName).Return(model.NewCounterMetric(tt.wantMetric.Name), tt.wantExists)
			} else {
				gr.On("Get", tt.metricName).Return((*model.Metrics[float64])(nil), tt.wantExists)
				cr.On("Get", tt.metricName).Return((*model.Metrics[int64])(nil), tt.wantExists)
			}

			s := NewService(gr, cr)
			metric, exists := s.GetMetric(context.Background(), tt.metricType, tt.metricName)
			assert.Equal(t, tt.wantMetric, metric)
			assert.Equal(t, tt.wantExists, exists)
		})
	}
}
