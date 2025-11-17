package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/yogenyslav/ya-metrics/internal/model"
	"github.com/yogenyslav/ya-metrics/pkg"
	"github.com/yogenyslav/ya-metrics/pkg/database"
	"github.com/yogenyslav/ya-metrics/tests/mocks"
)

func TestHandler_GetMetric(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		ms         func() metricService
		db         database.DB
		metricType string
		metricID   string
		wantCode   int
	}{
		{
			name: "GetMetric gauge with existing metric",
			ms: func() metricService {
				m := new(MockMetricService)
				m.On("GetMetric", mock.Anything, model.Gauge, "metric1").
					Return(model.NewGaugeMetric("metric1").ToDto(), true)
				return m
			},
			db: func() *mocks.MockDB {
				m := new(mocks.MockDB)
				return m
			}(),
			metricType: model.Gauge,
			metricID:   "metric1",
			wantCode:   http.StatusOK,
		},
		{
			name: "GetMetric counter with existing metric",
			ms: func() metricService {
				m := new(MockMetricService)
				m.On("GetMetric", mock.Anything, model.Counter, "metric1").
					Return(model.NewCounterMetric("metric1").ToDto(), true)
				return m
			},
			db: func() *mocks.MockDB {
				m := new(mocks.MockDB)
				return m
			}(),
			metricType: model.Counter,
			metricID:   "metric1",
			wantCode:   http.StatusOK,
		},
		{
			name: "GetMetric with non-existing metric",
			ms: func() metricService {
				m := new(MockMetricService)
				m.On("GetMetric", mock.Anything, model.Gauge, "non_existing_metric").
					Return((*model.MetricsDto)(nil), false)
				return m
			},
			db: func() *mocks.MockDB {
				m := new(mocks.MockDB)
				return m
			}(),
			metricType: model.Gauge,
			metricID:   "non_existing_metric",
			wantCode:   http.StatusNotFound,
		},
		{
			name: "GetMetric with invalid metric type",
			ms: func() metricService {
				m := new(MockMetricService)
				m.On("GetMetric", mock.Anything, "invalid", "metric1").
					Return((*model.MetricsDto)(nil), false)
				return m
			},
			db: func() *mocks.MockDB {
				m := new(mocks.MockDB)
				return m
			}(),
			metricType: "invalid",
			metricID:   "metric1",
			wantCode:   http.StatusNotFound,
		},
		{
			name: "GetMetric with missing metric name",
			ms: func() metricService {
				m := new(MockMetricService)
				m.On("GetMetric", mock.Anything, model.Gauge, "").
					Return((*model.MetricsDto)(nil), false)
				return m
			},
			db: func() *mocks.MockDB {
				m := new(mocks.MockDB)
				return m
			}(),
			metricType: model.Gauge,
			metricID:   "",
			wantCode:   http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			t.Run(tt.name+" raw request", func(t *testing.T) {
				t.Parallel()

				h := NewHandler(tt.ms(), tt.db)
				writer := httptest.NewRecorder()

				req := httptest.NewRequest(
					http.MethodGet,
					"/value/{"+metricTypeParam+"}/{"+metricIDParam+"}",
					nil,
				)
				req.SetPathValue(metricTypeParam, tt.metricType)
				req.SetPathValue(metricIDParam, tt.metricID)

				h.GetMetricRaw(writer, req)
				assert.Equal(t, tt.wantCode, writer.Code)

				if tt.wantCode == http.StatusOK {
					assert.Equal(t, "0", writer.Body.String())
				}
			})

			t.Run(tt.name+"json request", func(t *testing.T) {
				t.Parallel()

				h := NewHandler(tt.ms(), tt.db)
				writer := httptest.NewRecorder()

				data := model.MetricsDto{
					Type: tt.metricType,
					ID:   tt.metricID,
				}
				body, err := json.Marshal(data)
				require.NoError(t, err)

				req := httptest.NewRequest(
					http.MethodPost,
					"/value",
					bytes.NewReader(body),
				)

				h.GetMetricJSON(writer, req)
				assert.Equal(t, tt.wantCode, writer.Code)

				var want model.MetricsDto
				switch tt.metricType {
				case model.Gauge:
					want = model.MetricsDto{
						ID:    tt.metricID,
						Type:  model.Gauge,
						Value: pkg.Ptr(0.0),
					}
				case model.Counter:
					want = model.MetricsDto{
						ID:    tt.metricID,
						Type:  model.Counter,
						Delta: pkg.Ptr[int64](0),
					}
				}

				wantStr, err := json.Marshal(want)
				require.NoError(t, err)

				if tt.wantCode == http.StatusOK {
					assert.Equal(t, string(wantStr), writer.Body.String())
				}
			})
		})
	}
}
