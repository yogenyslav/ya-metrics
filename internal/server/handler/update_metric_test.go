package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/yogenyslav/ya-metrics/pkg/errs"
)

func TestHandler_UpdateMetric(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		ms          metricService
		metricType  string
		metrictName string
		metricValue string
		writer      http.ResponseWriter
		wantCode    int
	}{
		{
			name: "UpdateMetric with valid parameters",
			ms: func() metricService {
				m := new(MockMetricService)
				m.On("UpdateMetric", mock.Anything, "gauge", "metric1", "123.45").
					Return(nil)
				return m
			}(),
			metricType:  "gauge",
			metrictName: "metric1",
			metricValue: "123.45",
			writer:      httptest.NewRecorder(),
			wantCode:    http.StatusOK,
		},
		{
			name:        "UpdateMetric with missing metric name",
			ms:          new(MockMetricService),
			metricType:  "gauge",
			metricValue: "123.45",
			writer:      httptest.NewRecorder(),
			wantCode:    http.StatusNotFound,
		},
		{
			name: "UpdateMetric with invalid metric type",
			ms: func() metricService {
				m := new(MockMetricService)
				m.On("UpdateMetric", mock.Anything, "invalid", "metric1", "123.45").
					Return(errs.ErrInvalidMetricType)
				return m
			}(),
			metricType:  "invalid",
			metrictName: "metric1",
			metricValue: "123.45",
			writer:      httptest.NewRecorder(),
			wantCode:    http.StatusBadRequest,
		},
		{
			name: "UpdateMetric with invalid metric value",
			ms: func() metricService {
				m := new(MockMetricService)
				m.On("UpdateMetric", mock.Anything, "gauge", "metric1", "invalid").
					Return(errs.ErrInvalidMetricValue)
				return m
			}(),
			metricType:  "gauge",
			metrictName: "metric1",
			metricValue: "invalid",
			writer:      httptest.NewRecorder(),
			wantCode:    http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				t.Parallel()

				h := NewHandler(tt.ms)

				req := httptest.NewRequest(
					http.MethodPost,
					"/update/{"+metricTypeParam+"}/{"+metricNameParam+"}/{"+metricValueParam+"}",
					nil,
				)
				req.SetPathValue(metricTypeParam, tt.metricType)
				req.SetPathValue(metricNameParam, tt.metrictName)
				req.SetPathValue(metricValueParam, tt.metricValue)

				h.UpdateMetric(tt.writer, req)
				assert.Equal(t, tt.wantCode, tt.writer.(*httptest.ResponseRecorder).Code)
			},
		)
	}
}
