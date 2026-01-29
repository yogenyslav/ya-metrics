package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/ya-metrics/internal/model"
	"github.com/yogenyslav/ya-metrics/pkg/database"
)

const (
	metricTypeParam  = "metricType"
	metricIDParam    = "metricID"
	metricValueParam = "metricValue"
)

type metricService interface {
	UpdateMetric(ctx context.Context, metric *model.MetricsDto) error
	UpdateMetricsBatch(ctx context.Context, metrics []*model.MetricsDto) error
	GetMetric(ctx context.Context, metricType, metricID string) (*model.MetricsDto, error)
	ListMetrics(ctx context.Context) ([]*model.MetricsDto, error)
}

//go:generate mockgen -destination=../../../tests/mocks/audit.go -package=mocks . auditLogger
type auditLogger interface {
	LogMetrics(ctx context.Context, metrics []string, ipAddr string) error
}

// Handler serves HTTP requests.
type Handler struct {
	ms    metricService
	db    database.DB
	audit auditLogger
}

// NewHandler creates new HTTP handler.
func NewHandler(ms metricService, db database.DB, audit auditLogger) *Handler {
	return &Handler{
		ms:    ms,
		db:    db,
		audit: audit,
	}
}

// RegisterRoutes registers HTTP routes.
func (h *Handler) RegisterRoutes(router chi.Router) {
	router.Get("/", h.ListMetrics)
	router.Get("/ping", h.Ping)
	router.Post("/value/", h.GetMetricJSON)
	router.Get("/value/{metricType}/{metricID}", h.GetMetricRaw)
	router.Post("/updates/", h.UpdateMetricsBatch)
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
	log.Debug().Err(wrappedErr).Msg("handler error")

	for e, code := range errStatusCodes {
		if errors.Is(wrappedErr, e) {
			statusCode = code
			err = e.Error()
			break
		}
	}

	http.Error(w, err, statusCode)
}
