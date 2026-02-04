package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/yogenyslav/ya-metrics/internal/model"
	"github.com/yogenyslav/ya-metrics/internal/server/repository"
	"github.com/yogenyslav/ya-metrics/internal/server/service"
	"github.com/yogenyslav/ya-metrics/pkg/database"
	"github.com/yogenyslav/ya-metrics/tests/mocks"
	gomock "go.uber.org/mock/gomock"
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

func BenchmarkHandler_ListMetrics(b *testing.B) {
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
			req, writer := getRequestAndWriter(b, http.MethodGet, "/", nil)
			h.ListMetrics(writer, req)
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
			QuerySlice(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, dest interface{}, query string, args ...any) error {
				metrics := []model.Metrics[float64]{
					{
						ID:    "gauge_metric",
						Type:  model.Gauge,
						Value: 0,
					},
				}
				metricsRaw, _ := json.Marshal(metrics)
				return json.Unmarshal(metricsRaw, dest)
			}).
			AnyTimes()

		mockDB.EXPECT().
			QuerySlice(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, dest interface{}, query string, args ...any) error {
				metrics := []model.Metrics[int64]{
					{
						ID:    "counter_metric",
						Type:  model.Counter,
						Value: 0,
					},
				}
				metricsRaw, _ := json.Marshal(metrics)
				return json.Unmarshal(metricsRaw, dest)
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
			req, writer := getRequestAndWriter(b, http.MethodGet, "/", nil)
			h.ListMetrics(writer, req)
		}
	})
}
