package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/yogenyslav/ya-metrics/internal/config"
	"github.com/yogenyslav/ya-metrics/internal/model"
	"github.com/yogenyslav/ya-metrics/internal/server/audit"
	"github.com/yogenyslav/ya-metrics/internal/server/repository"
	"github.com/yogenyslav/ya-metrics/internal/server/service"
	"github.com/yogenyslav/ya-metrics/pkg"
	"github.com/yogenyslav/ya-metrics/pkg/errs"
	"github.com/yogenyslav/ya-metrics/tests/mocks"
	"go.uber.org/mock/gomock"
)

func ExampleHandler() {
	// Setup dependencies.
	uow := mocks.NewMockUnitOfWork(gomock.NewController(nil))

	gaugeRepo := repository.NewMetricInMemRepo[float64](nil)
	counterRepo := repository.NewMetricInMemRepo[int64](nil)
	metricService := service.NewService(gaugeRepo, counterRepo, uow)

	auditCfg := &config.AuditConfig{File: "audit.log"}
	audit := audit.New(auditCfg)

	h := NewHandler(metricService, nil, audit)

	router := chi.NewRouter()
	h.RegisterRoutes(router)

	// HTTP requests to endpoints.
	go http.ListenAndServe(":8080", router)

	updateMetricsReq := model.MetricsDto{
		ID:    "example_gauge",
		Type:  model.Gauge,
		Value: pkg.Ptr(1.23),
	}
	updateBody, _ := json.Marshal(updateMetricsReq)

	resp, err := http.Post("http://localhost:8080/update/", "application/json", bytes.NewReader(updateBody))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Printf("Update metric status: %s\n", resp.Status)

	getResp, err := http.Get("http://localhost:8080/value/gauge/example_gauge")
	if err != nil {
		panic(err)
	}
	defer getResp.Body.Close()
	fmt.Printf("Get metric status: %s\n", getResp.Status)

	// Output:
	// Update metric status: 200 OK
	// Get metric status: 200 OK
}

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
