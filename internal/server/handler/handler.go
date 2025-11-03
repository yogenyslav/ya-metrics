package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/yogenyslav/ya-metrics/internal/model"
)

const (
	metricTypeParam  = "metricType"
	metricNameParam  = "metricName"
	metricValueParam = "metricValue"
)

type metricService interface {
	UpdateMetric(ctx context.Context, metricType, name, rawValue string) error
	GetMetric(ctx context.Context, metricType, metricName string) (*model.MetricsDto, bool)
	ListMetrics(ctx context.Context) []*model.MetricsDto
}

// Handler serves HTTP requests.
type Handler struct {
	ms metricService
}

// NewHandler creates new HTTP handler.
func NewHandler(ms metricService) *Handler {
	return &Handler{
		ms: ms,
	}
}

// RegisterRoutes registers HTTP routes.
func (h *Handler) RegisterRoutes(router chi.Router) {
	router.HandleFunc("GET /", h.ListMetrics)
	router.HandleFunc("POST /value", h.GetMetricJSON)
	router.HandleFunc("GET /value/{metricType}/{metricName}", h.GetMetricRaw)
	router.HandleFunc(
		"POST /update",
		h.UpdateMetricJSON,
	)
	router.HandleFunc(
		"POST /update/{metricType}/{metricName}/{metricValue}",
		h.UpdateMetricRaw,
	)
}

func (h *Handler) sendError(w http.ResponseWriter, wrappedErr error) {
	if wrappedErr == nil {
		return
	}

	statusCode := http.StatusInternalServerError
	err := wrappedErr.Error()

	for e, code := range errStatusCodes {
		if errors.Is(wrappedErr, e) {
			statusCode = code
			err = e.Error()
			break
		}
	}

	http.Error(w, err, statusCode)
}
