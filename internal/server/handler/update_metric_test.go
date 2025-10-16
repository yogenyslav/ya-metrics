package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
)

type MockMetricService struct {
	mock.Mock
}

func (m *MockMetricService) UpdateMetric(ctx context.Context, metricType, metricName, metricValueRaw string) error {
	args := m.Called(ctx, metricType, metricName, metricValueRaw)
	return args.Error(0)
}

func TestHandler_UpdateMetric(t *testing.T) {
	t.Parallel()

	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name     string
		args     args
		wantCode int
	}{
		{
			name: "UpdateMetric with valid parameters",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("POST", "/update/gauge/metric1/123.45", nil).WithContext(context.Background()),
			},
			wantCode: http.StatusOK,
		},
		{
			name: "UpdateMetric with missing metric name",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("POST", "/update/gauge//123.45", nil).WithContext(context.Background()),
			},
			wantCode: http.StatusNotFound,
		},
		{
			name: "UpdateMetric with invalid metric type",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("POST", "/update/invalid/metric1/123.45", nil).WithContext(context.Background()),
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "UpdateMetric with invalid metric value",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("POST", "/update/gauge/metric1/invalid", nil).WithContext(context.Background()),
			},
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				t.Parallel()

				ms := &MockMetricService{}
				h := &Handler{
					ms: ms,
				}

				if tt.wantCode == http.StatusOK {
					ms.On("UpdateMetric", mock.Anything, "gauge", "metric1", "123.45").Return(nil).Once()
				}

				h.UpdateMetric(tt.args.w, tt.args.r)
			},
		)
	}
}
