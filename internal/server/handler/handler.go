package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/yogenyslav/ya-metrics/internal/model"
	"github.com/yogenyslav/ya-metrics/pkg/database"
)

const (
	metricTypeParam  = "metricType"
	metricIDParam    = "metricID"
	metricValueParam = "metricValue"
)

type metricService interface {
	UpdateMetric(ctx context.Context, metricType, metricID, rawValue string) error
	GetMetric(ctx context.Context, metricType, metricID string) (*model.MetricsDto, error)
	ListMetrics(ctx context.Context) ([]*model.MetricsDto, error)
}

// Handler serves HTTP requests.
type Handler struct {
	ms metricService
	db database.DB
}

// NewHandler creates new HTTP handler.
func NewHandler(ms metricService, db database.DB) *Handler {
	return &Handler{
		ms: ms,
		db: db,
	}
}

// RegisterRoutes registers HTTP routes.
func (h *Handler) RegisterRoutes(router chi.Router) {
	router.Get("/", h.ListMetrics)
	router.Get("/ping", h.Ping)
	router.Post("/value/", h.GetMetricJSON)
	router.Get("/value/{metricType}/{metricID}", h.GetMetricRaw)
	router.Post(
		"/update/",
		h.UpdateMetricJSON,
	)
	router.Post(
		"/update/{metricType}/{metricID}/{metricValue}",
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
