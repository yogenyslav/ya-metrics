package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/yogenyslav/ya-metrics/internal/model"
	"github.com/yogenyslav/ya-metrics/pkg"
	"github.com/yogenyslav/ya-metrics/tests/mocks"
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
				{ID: "mem_alloc", Type: model.Gauge, Value: pkg.Ptr(0.0)},
				{ID: "cpu_usage", Type: model.Gauge, Value: pkg.Ptr(0.0)},
				{ID: "request_count", Type: model.Counter, Delta: pkg.Ptr[int64](0)},
				{ID: "error_count", Type: model.Counter, Delta: pkg.Ptr[int64](0)},
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

			gr := &mocks.MockGaugeRepo{}
			cr := &mocks.MockCounterRepo{}

			gr.On("List", mock.Anything).Return(tt.gaugeMetrics, nil)
			cr.On("List", mock.Anything).Return(tt.counterMetrics, nil)

			s := NewService(gr, cr)
			metrics, err := s.ListMetrics(ctx)
			assert.NoError(t, err)
			assert.ElementsMatch(t, tt.want, metrics)
		})
	}
}
