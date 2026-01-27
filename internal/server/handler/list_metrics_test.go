package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/yogenyslav/ya-metrics/internal/model"
	"github.com/yogenyslav/ya-metrics/pkg/database"
	"github.com/yogenyslav/ya-metrics/tests/mocks"
)

func TestHandler_ListMetrics(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		ms          metricService
		db          database.DB
		audit       auditLogger
		writer      http.ResponseWriter
		wantMetrics []*model.MetricsDto
	}{
		{
			name: "ListMetrics with existing metrics",
			ms: func() metricService {
				m := new(mocks.MockMetricService)
				m.On("ListMetrics", mock.Anything).
					Return([]*model.MetricsDto{
						model.NewGaugeMetric("gauge1").ToDto(),
						model.NewCounterMetric("counter1").ToDto(),
					}, nil)
				return m
			}(),
			db: func() *mocks.MockDB {
				m := new(mocks.MockDB)
				return m
			}(),
			audit: func() auditLogger {
				m := new(mocks.MockauditLogger)
				return m
			}(),
			writer: httptest.NewRecorder(),
			wantMetrics: []*model.MetricsDto{
				model.NewGaugeMetric("gauge1").ToDto(),
				model.NewCounterMetric("counter1").ToDto(),
			},
		},
		{
			name: "ListMetrics with no metrics",
			ms: func() metricService {
				m := new(mocks.MockMetricService)
				m.On("ListMetrics", mock.Anything).
					Return([]*model.MetricsDto{}, nil)
				return m
			}(),
			db: func() *mocks.MockDB {
				m := new(mocks.MockDB)
				return m
			}(),
			audit: func() auditLogger {
				m := new(mocks.MockauditLogger)
				return m
			}(),
			writer:      httptest.NewRecorder(),
			wantMetrics: []*model.MetricsDto{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			h := NewHandler(tt.ms, tt.db, tt.audit)
			req := httptest.NewRequest(http.MethodGet, "/metrics", nil)

			h.ListMetrics(tt.writer, req)
			wantBody, err := json.Marshal(tt.wantMetrics)
			require.NoError(t, err)

			assert.ElementsMatch(t, wantBody, tt.writer.(*httptest.ResponseRecorder).Body.Bytes())
			assert.Equal(t, http.StatusOK, tt.writer.(*httptest.ResponseRecorder).Code)
		})
	}
}
