package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yogenyslav/ya-metrics/internal/model"
	"github.com/yogenyslav/ya-metrics/pkg/errs"
)

func getRequestAndWriter(b *testing.B, method, url string, body []byte) (*http.Request, *httptest.ResponseRecorder) {
	b.Helper()
	req := httptest.NewRequest(
		method,
		url,
		bytes.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")
	writer := httptest.NewRecorder()
	return req, writer
}

func getGaugeData(b *testing.B, i float64) []byte {
	b.Helper()
	data := model.MetricsDto{
		ID:    "gauge_metric",
		Type:  model.Gauge,
		Value: &i,
	}
	raw, _ := json.Marshal(data)
	return raw
}

func getCounterData(b *testing.B, i int64) []byte {
	b.Helper()
	data := model.MetricsDto{
		ID:    "counter_metric",
		Type:  model.Counter,
		Delta: &i,
	}
	raw, _ := json.Marshal(data)
	return raw
}

func getBatchData(b *testing.B, gaugeVal float64, counterVal int64) []byte {
	b.Helper()
	data := []model.MetricsDto{
		{
			ID:    "gauge_metric",
			Type:  model.Gauge,
			Value: &gaugeVal,
		},
		{
			ID:    "counter_metric",
			Type:  model.Counter,
			Delta: &counterVal,
		},
	}
	raw, _ := json.Marshal(data)
	return raw
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
