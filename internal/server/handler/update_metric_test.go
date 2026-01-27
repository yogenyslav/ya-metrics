package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/yogenyslav/ya-metrics/internal/model"
	"github.com/yogenyslav/ya-metrics/pkg"
	"github.com/yogenyslav/ya-metrics/pkg/database"
	"github.com/yogenyslav/ya-metrics/pkg/errs"
	"github.com/yogenyslav/ya-metrics/tests/mocks"
)

func TestHandler_UpdateMetric(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		ms          func() metricService
		db          database.DB
		audit       func() auditLogger
		metricType  string
		metrictName string
		metricValue string
		wantCode    int
	}{
		{
			name: "UpdateMetric gauge with valid parameters",
			ms: func() metricService {
				m := new(mocks.MockMetricService)
				m.On("UpdateMetric", mock.Anything, &model.MetricsDto{
					ID:    "metric1",
					Type:  model.Gauge,
					Value: pkg.Ptr(123.45),
				}).
					Return(nil)
				return m
			},
			db: func() *mocks.MockDB {
				m := new(mocks.MockDB)
				return m
			}(),
			audit: func() auditLogger {
				m := mocks.NewMockauditLogger(gomock.NewController(t))
				m.EXPECT().
					LogMetrics(gomock.Any(), []string{"metric1"}, gomock.Any()).
					Return(nil)
				return m
			},
			metricType:  model.Gauge,
			metrictName: "metric1",
			metricValue: "123.45",
			wantCode:    http.StatusOK,
		},
		{
			name: "UpdateMetric counter with valid parameters",
			ms: func() metricService {
				m := new(mocks.MockMetricService)
				m.On("UpdateMetric", mock.Anything, &model.MetricsDto{
					ID:    "metric1",
					Type:  model.Counter,
					Delta: pkg.Ptr(int64(123)),
				}).
					Return(nil)
				return m
			},
			db: func() *mocks.MockDB {
				m := new(mocks.MockDB)
				return m
			}(),
			audit: func() auditLogger {
				m := mocks.NewMockauditLogger(gomock.NewController(t))
				m.EXPECT().
					LogMetrics(gomock.Any(), []string{"metric1"}, gomock.Any()).
					Return(nil)
				return m
			},
			metricType:  model.Counter,
			metrictName: "metric1",
			metricValue: "123",
			wantCode:    http.StatusOK,
		},
		{
			name: "UpdateMetric with missing metric name",
			ms:   func() metricService { return new(mocks.MockMetricService) },
			audit: func() auditLogger {
				m := mocks.NewMockauditLogger(gomock.NewController(t))
				return m
			},
			metricType:  model.Gauge,
			metricValue: "123.45",
			wantCode:    http.StatusNotFound,
		},
		{
			name: "UpdateMetric with invalid metric type",
			ms: func() metricService {
				m := new(mocks.MockMetricService)
				m.On("UpdateMetric", mock.Anything, &model.MetricsDto{
					ID:   "metric1",
					Type: "invalid",
				}).
					Return(errs.ErrInvalidMetricType)
				return m
			},
			db: func() *mocks.MockDB {
				m := new(mocks.MockDB)
				return m
			}(),
			audit: func() auditLogger {
				m := mocks.NewMockauditLogger(gomock.NewController(t))
				return m
			},
			metricType:  "invalid",
			metrictName: "metric1",
			metricValue: "123.45",
			wantCode:    http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				t.Parallel()

				t.Run(tt.name+" raw request", func(t *testing.T) {
					t.Parallel()

					h := NewHandler(tt.ms(), tt.db, tt.audit())

					writer := httptest.NewRecorder()
					req := httptest.NewRequest(
						http.MethodPost,
						"/update/{"+metricTypeParam+"}/{"+metricIDParam+"}/{"+metricValueParam+"}",
						nil,
					)
					req.Header.Set("Content-Type", "text/plain")
					req.SetPathValue(metricTypeParam, tt.metricType)
					req.SetPathValue(metricIDParam, tt.metrictName)
					req.SetPathValue(metricValueParam, tt.metricValue)

					h.UpdateMetricRaw(writer, req)
					assert.Equal(t, tt.wantCode, writer.Code)
				})

				t.Run(tt.name+" json request", func(t *testing.T) {
					t.Parallel()

					h := NewHandler(tt.ms(), tt.db, tt.audit())

					data := model.MetricsDto{
						Type: tt.metricType,
						ID:   tt.metrictName,
					}
					switch tt.metricType {
					case model.Gauge:
						v, err := strconv.ParseFloat(tt.metricValue, 64)
						if tt.wantCode == http.StatusOK {
							require.NoError(t, err)
						}
						data.Value = &v
					case model.Counter:
						v, err := strconv.ParseInt(tt.metricValue, 10, 64)
						if tt.wantCode == http.StatusOK {
							require.NoError(t, err)
						}
						data.Delta = &v
					}

					body, err := json.Marshal(data)
					require.NoError(t, err)

					writer := httptest.NewRecorder()
					req := httptest.NewRequest(
						http.MethodPost,
						"/update",
						bytes.NewReader(body),
					)
					req.Header.Set("Content-Type", "application/json")

					h.UpdateMetricJSON(writer, req)
					assert.Equal(t, tt.wantCode, writer.Code)
				})
			},
		)
	}
}

func TestHandler_UpdateMetricsBatch(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		ms       func() metricService
		audit    func() auditLogger
		metrics  []model.MetricsDto
		wantCode int
	}{
		{
			name: "UpdateMetricsBatch with valid metrics",
			ms: func() metricService {
				m := new(mocks.MockMetricService)
				m.On("UpdateMetricsBatch", mock.Anything, []*model.MetricsDto{
					{
						ID:    "metric1",
						Type:  model.Gauge,
						Value: pkg.Ptr(123.45),
					},
					{
						ID:    "metric2",
						Type:  model.Counter,
						Delta: pkg.Ptr(int64(100)),
					},
				}).Return(nil)
				return m
			},
			audit: func() auditLogger {
				m := mocks.NewMockauditLogger(gomock.NewController(t))
				m.EXPECT().
					LogMetrics(gomock.Any(), []string{"metric1", "metric2"}, gomock.Any()).
					Return(nil)
				return m
			},
			metrics: []model.MetricsDto{
				{
					ID:    "metric1",
					Type:  model.Gauge,
					Value: pkg.Ptr(123.45),
				},
				{
					ID:    "metric2",
					Type:  model.Counter,
					Delta: pkg.Ptr(int64(100)),
				},
			},
			wantCode: http.StatusOK,
		},
		{
			name: "UpdateMetricsBatch with invalid metric type",
			ms: func() metricService {
				m := new(mocks.MockMetricService)
				m.On("UpdateMetricsBatch", mock.Anything, mock.Anything).
					Return(errs.ErrInvalidMetricType)
				return m
			},
			audit: func() auditLogger {
				m := mocks.NewMockauditLogger(gomock.NewController(t))
				return m
			},
			metrics: []model.MetricsDto{
				{
					ID:    "metric1",
					Type:  "invalid",
					Value: pkg.Ptr(123.45),
				},
			},
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				t.Parallel()

				h := NewHandler(tt.ms(), nil, tt.audit())

				body, err := json.Marshal(tt.metrics)
				require.NoError(t, err)

				writer := httptest.NewRecorder()
				req := httptest.NewRequest(
					http.MethodPost,
					"/updates",
					bytes.NewReader(body),
				)
				req.Header.Set("Content-Type", "application/json")

				h.UpdateMetricsBatch(writer, req)
				assert.Equal(t, tt.wantCode, writer.Code)
			},
		)
	}
}
