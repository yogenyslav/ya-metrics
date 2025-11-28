package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/yogenyslav/ya-metrics/internal/model"
	"github.com/yogenyslav/ya-metrics/pkg"
	"github.com/yogenyslav/ya-metrics/pkg/errs"
	"github.com/yogenyslav/ya-metrics/tests/mocks"
)

func TestService_GetMetric(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		metricType string
		metricID   string
		wantMetric *model.MetricsDto
		wantErr    bool
	}{
		{
			name:       "Get existing gauge metric",
			metricType: model.Gauge,
			metricID:   "mem_alloc",
			wantMetric: &model.MetricsDto{ID: "mem_alloc", Type: model.Gauge, Value: pkg.Ptr(0.0)},
			wantErr:    false,
		},
		{
			name:       "Get non-existing gauge metric",
			metricType: model.Gauge,
			metricID:   "non_existing_gauge",
			wantMetric: nil,
			wantErr:    true,
		},
		{
			name:       "Get existing counter metric",
			metricType: model.Counter,
			metricID:   "request_count",
			wantMetric: &model.MetricsDto{ID: "request_count", Type: model.Counter, Delta: pkg.Ptr[int64](0)},
			wantErr:    false,
		},
		{
			name:       "Get non-existing counter metric",
			metricType: model.Counter,
			metricID:   "non_existing_counter",
			wantMetric: nil,
			wantErr:    true,
		},
		{
			name:       "Get metric with invalid type",
			metricType: "invalid_type",
			metricID:   "some_metric",
			wantMetric: nil,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gr := &mocks.MockGaugeRepo{}
			cr := &mocks.MockCounterRepo{}

			if tt.wantMetric != nil {
				switch tt.metricType {
				case model.Gauge:
					gr.On("Get", mock.Anything, tt.metricID, model.Gauge).Return(
						model.NewGaugeMetric(tt.metricID), nil,
					)
				case model.Counter:
					cr.On("Get", mock.Anything, tt.metricID, model.Counter).Return(
						model.NewCounterMetric(tt.metricID), nil,
					)
				}
			} else {
				switch tt.metricType {
				case model.Gauge:
					gr.On("Get", mock.Anything, tt.metricID, model.Gauge).Return(
						&model.Metrics[float64]{}, errs.ErrMetricNotFound,
					)
				case model.Counter:
					cr.On("Get", mock.Anything, tt.metricID, model.Counter).Return(
						&model.Metrics[int64]{}, errs.ErrMetricNotFound,
					)
				}
			}

			s := NewService(gr, cr)
			metric, err := s.GetMetric(context.Background(), tt.metricType, tt.metricID)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantMetric, metric)
			}
		})
	}
}
