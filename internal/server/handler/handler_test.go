package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/yogenyslav/ya-metrics/internal/model"
	"github.com/yogenyslav/ya-metrics/pkg/errs"
)

type MockMetricService struct {
	mock.Mock
}

func (m *MockMetricService) UpdateMetric(ctx context.Context, metricType, metricName, metricValueRaw string) error {
	args := m.Called(ctx, metricType, metricName, metricValueRaw)
	return args.Error(0)
}

func (m *MockMetricService) GetMetric(ctx context.Context, metricType, metricName string) (*model.MetricsDto, bool) {
	args := m.Called(ctx, metricType, metricName)
	return args.Get(0).(*model.MetricsDto), args.Bool(1)
}

func (m *MockMetricService) ListMetrics(ctx context.Context) []*model.MetricsDto {
	args := m.Called(ctx)
	return args.Get(0).([]*model.MetricsDto)
}

func TestHandler_sendError(t *testing.T) {
	t.Parallel()

	type fields struct {
		ms metricService
	}
	type args struct {
		w          http.ResponseWriter
		wrappedErr error
	}

	tests := []struct {
		name     string
		fields   fields
		args     args
		wantCode int
	}{
		{
			name: "sendError with mapped error",
			args: args{
				w:          httptest.NewRecorder(),
				wrappedErr: errs.ErrInvalidMetricType,
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "sendError with unmapped error",
			args: args{
				w:          httptest.NewRecorder(),
				wrappedErr: errors.New("something went wrong"),
			},
			wantCode: http.StatusInternalServerError,
		},
		{
			name: "sendError with nil error",
			args: args{
				w:          httptest.NewRecorder(),
				wrappedErr: nil,
			},
			wantCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				t.Parallel()

				h := &Handler{
					ms: tt.fields.ms,
				}

				h.sendError(tt.args.w, tt.args.wrappedErr)
				if tt.args.wrappedErr != nil {
					recorder := tt.args.w.(*httptest.ResponseRecorder)

					assert.Equal(t, tt.wantCode, recorder.Code)
					assert.Contains(t, recorder.Body.String(), tt.args.wrappedErr.Error())
				}
			},
		)
	}
}
