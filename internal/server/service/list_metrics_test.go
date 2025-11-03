package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yogenyslav/ya-metrics/internal/model"
	"github.com/yogenyslav/ya-metrics/pkg"
)

func TestService_ListMetrics(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name           string
		gaugeMetrics   []model.Metrics[float64]
		counterMetrics []model.Metrics[int64]
		want           []*model.MetricsDto
	}{
		{
			name: "List all metrics",
			gaugeMetrics: []model.Metrics[float64]{
				*model.NewGaugeMetric("mem_alloc"),
				*model.NewGaugeMetric("cpu_usage"),
			},
			counterMetrics: []model.Metrics[int64]{
				*model.NewCounterMetric("request_count"),
				*model.NewCounterMetric("error_count"),
			},
			want: []*model.MetricsDto{
				{Name: "mem_alloc", Type: model.Gauge, Value: pkg.Ptr(0.0)},
				{Name: "cpu_usage", Type: model.Gauge, Value: pkg.Ptr(0.0)},
				{Name: "request_count", Type: model.Counter, Delta: pkg.Ptr[int64](0)},
				{Name: "error_count", Type: model.Counter, Delta: pkg.Ptr[int64](0)},
			},
		},
		{
			name:           "List with no metrics",
			gaugeMetrics:   []model.Metrics[float64]{},
			counterMetrics: []model.Metrics[int64]{},
			want:           []*model.MetricsDto{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gr := &MockGaugeRepo{}
			cr := &MockCounterRepo{}

			gr.On("List").Return(tt.gaugeMetrics)
			cr.On("List").Return(tt.counterMetrics)

			s := NewService(gr, cr)
			metrics := s.ListMetrics(ctx)

			assert.ElementsMatch(t, tt.want, metrics)
		})
	}
}
