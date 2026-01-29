package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/yogenyslav/ya-metrics/internal/model"
	"github.com/yogenyslav/ya-metrics/internal/server/repository"
	"github.com/yogenyslav/ya-metrics/internal/server/service"
	"github.com/yogenyslav/ya-metrics/pkg"
	"github.com/yogenyslav/ya-metrics/pkg/database"
	"github.com/yogenyslav/ya-metrics/tests/mocks"
	gomock "go.uber.org/mock/gomock"
)

func TestHandler_GetMetric(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		ms         func() metricService
		db         database.DB
		audit      func() auditLogger
		metricType string
		metricID   string
		wantCode   int
	}{
		{
			name: "GetMetric gauge with existing metric",
			ms: func() metricService {
				m := new(mocks.MockMetricService)
				m.On("GetMetric", mock.Anything, model.Gauge, "metric1").
					Return(model.NewGaugeMetric("metric1").ToDto(), nil)
				return m
			},
			db: func() *mocks.MockDB {
				m := new(mocks.MockDB)
				return m
			}(),
			audit: func() auditLogger {
				m := new(mocks.MockauditLogger)
				return m
			},
			metricType: model.Gauge,
			metricID:   "metric1",
			wantCode:   http.StatusOK,
		},
		{
			name: "GetMetric counter with existing metric",
			ms: func() metricService {
				m := new(mocks.MockMetricService)
				m.On("GetMetric", mock.Anything, model.Counter, "metric1").
					Return(model.NewCounterMetric("metric1").ToDto(), nil)
				return m
			},
			db: func() *mocks.MockDB {
				m := new(mocks.MockDB)
				return m
			}(),
			audit: func() auditLogger {
				m := new(mocks.MockauditLogger)
				return m
			},
			metricType: model.Counter,
			metricID:   "metric1",
			wantCode:   http.StatusOK,
		},
		{
			name: "GetMetric with non-existing metric",
			ms: func() metricService {
				m := new(mocks.MockMetricService)
				m.On("GetMetric", mock.Anything, model.Gauge, "non_existing_metric").
					Return((*model.MetricsDto)(nil), errors.New("not found"))
				return m
			},
			db: func() *mocks.MockDB {
				m := new(mocks.MockDB)
				return m
			}(),
			audit: func() auditLogger {
				m := new(mocks.MockauditLogger)
				return m
			},
			metricType: model.Gauge,
			metricID:   "non_existing_metric",
			wantCode:   http.StatusNotFound,
		},
		{
			name: "GetMetric with invalid metric type",
			ms: func() metricService {
				m := new(mocks.MockMetricService)
				m.On("GetMetric", mock.Anything, "invalid", "metric1").
					Return((*model.MetricsDto)(nil), errors.New("invalid metric type"))
				return m
			},
			db: func() *mocks.MockDB {
				m := new(mocks.MockDB)
				return m
			}(),
			audit: func() auditLogger {
				m := new(mocks.MockauditLogger)
				return m
			},
			metricType: "invalid",
			metricID:   "metric1",
			wantCode:   http.StatusNotFound,
		},
		{
			name: "GetMetric with missing metric name",
			ms: func() metricService {
				m := new(mocks.MockMetricService)
				m.On("GetMetric", mock.Anything, model.Gauge, "").
					Return((*model.MetricsDto)(nil), errors.New("metric ID is required"))
				return m
			},
			db: func() *mocks.MockDB {
				m := new(mocks.MockDB)
				return m
			}(),
			audit: func() auditLogger {
				m := new(mocks.MockauditLogger)
				return m
			},
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

				h := NewHandler(tt.ms(), tt.db, tt.audit())
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

				h := NewHandler(tt.ms(), tt.db, tt.audit())
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

func BenchmarkHandler_GetMetric(b *testing.B) {
	b.Run("in memory repo", func(b *testing.B) {
		gaugeRepo := repository.NewMetricInMemRepo(repository.StorageState[float64]{})
		counterRepo := repository.NewMetricInMemRepo(repository.StorageState[int64]{})
		svc := service.NewService(gaugeRepo, counterRepo, nil)
		h := NewHandler(svc, nil, nil)

		svc.UpdateMetric(b.Context(), &model.MetricsDto{
			ID:    "gauge_metric",
			Type:  model.Gauge,
			Value: new(float64),
		})
		svc.UpdateMetric(b.Context(), &model.MetricsDto{
			ID:    "counter_metric",
			Type:  model.Counter,
			Delta: new(int64),
		})

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			req, writer := getRequestAndWriter(
				b,
				http.MethodGet,
				"/value/{"+metricTypeParam+"}/{"+metricIDParam+"}",
				nil,
			)
			if i%2 == 0 {
				req.SetPathValue(metricTypeParam, model.Gauge)
				req.SetPathValue(metricIDParam, "gauge_metric")
			} else {
				req.SetPathValue(metricTypeParam, model.Counter)
				req.SetPathValue(metricIDParam, "counter_metric")
			}

			h.GetMetricRaw(writer, req)
		}
	})

	b.Run("mock postgres repo", func(b *testing.B) {
		mockDB := mocks.NewMockDB(gomock.NewController(b))
		gaugeRepo := repository.NewMetricPostgresRepo[float64](mockDB)
		counterRepo := repository.NewMetricPostgresRepo[int64](mockDB)
		svc := service.NewService(gaugeRepo, counterRepo, nil)
		h := NewHandler(svc, mockDB, nil)

		mockDB.EXPECT().
			Exec(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(int64(1), nil).
			Times(2)

		mockDB.EXPECT().
			QueryRow(gomock.Any(), gomock.Any(), gomock.Any(), "gauge_metric", model.Gauge).
			DoAndReturn(func(ctx any, dest any, query string, args ...any) error {
				m := dest.(*model.MetricsDto)
				m.ID = "gauge_metric"
				m.Type = model.Gauge
				val := 0.0
				m.Value = &val
				return nil
			}).
			AnyTimes()

		mockDB.EXPECT().
			QueryRow(gomock.Any(), gomock.Any(), gomock.Any(), "counter_metric", model.Counter).
			DoAndReturn(func(ctx any, dest any, query string, args ...any) error {
				m := dest.(*model.MetricsDto)
				m.ID = "counter_metric"
				m.Type = model.Counter
				delta := int64(0)
				m.Delta = &delta
				return nil
			}).
			AnyTimes()

		svc.UpdateMetric(b.Context(), &model.MetricsDto{
			ID:    "gauge_metric",
			Type:  model.Gauge,
			Value: new(float64),
		})
		svc.UpdateMetric(b.Context(), &model.MetricsDto{
			ID:    "counter_metric",
			Type:  model.Counter,
			Delta: new(int64),
		})

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			req, writer := getRequestAndWriter(
				b,
				http.MethodGet,
				"/value/{"+metricTypeParam+"}/{"+metricIDParam+"}",
				nil,
			)
			if i%2 == 0 {
				req.SetPathValue(metricTypeParam, model.Gauge)
				req.SetPathValue(metricIDParam, "gauge_metric")
			} else {
				req.SetPathValue(metricTypeParam, model.Counter)
				req.SetPathValue(metricIDParam, "counter_metric")
			}

			h.GetMetricRaw(writer, req)
		}
	})
}
