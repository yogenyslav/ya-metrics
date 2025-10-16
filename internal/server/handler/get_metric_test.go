package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/yogenyslav/ya-metrics/internal/model"
)

func TestHandler_GetMetric(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		ms         metricService
		metricType string
		metricName string
		writer     http.ResponseWriter
		wantCode   int
	}{
		{
			name: "GetMetric with existing metric",
			ms: func() metricService {
				m := new(MockMetricService)
				m.On("GetMetric", mock.Anything, "gauge", "metric1").
					Return(model.NewGaugeMetric("metric1").ToDto(), true)
				return m
			}(),
			writer:     httptest.NewRecorder(),
			metricType: "gauge",
			metricName: "metric1",
			wantCode:   http.StatusOK,
		},
		{
			name: "GetMetric with non-existing metric",
			ms: func() metricService {
				m := new(MockMetricService)
				m.On("GetMetric", mock.Anything, "gauge", "non_existing_metric").
					Return((*model.MetricsDto)(nil), false)
				return m
			}(),
			writer:     httptest.NewRecorder(),
			metricType: "gauge",
			metricName: "non_existing_metric",
			wantCode:   http.StatusNotFound,
		},
		{
			name: "GetMetric with invalid metric type",
			ms: func() metricService {
				m := new(MockMetricService)
				m.On("GetMetric", mock.Anything, "invalid", "metric1").
					Return((*model.MetricsDto)(nil), false)
				return m
			}(),
			writer:     httptest.NewRecorder(),
			metricType: "invalid",
			metricName: "metric1",
			wantCode:   http.StatusNotFound,
		},
		{
			name: "GetMetric with missing metric name",
			ms: func() metricService {
				m := new(MockMetricService)
				m.On("GetMetric", mock.Anything, "gauge", "").
					Return((*model.MetricsDto)(nil), false)
				return m
			}(),
			writer:     httptest.NewRecorder(),
			metricType: "gauge",
			metricName: "",
			wantCode:   http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			h := NewHandler(tt.ms)

			req := httptest.NewRequest(
				http.MethodGet,
				"/value/{"+metricTypeParam+"}/{"+metricNameParam+"}",
				nil,
			)
			req.SetPathValue(metricTypeParam, tt.metricType)
			req.SetPathValue(metricNameParam, tt.metricName)

			h.GetMetric(tt.writer, req)
			assert.Equal(t, tt.wantCode, tt.writer.(*httptest.ResponseRecorder).Code)
		})
	}
}
