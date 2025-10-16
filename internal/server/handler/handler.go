package handler

import (
	"context"
	"errors"
	"net/http"
)

const (
	metricTypeParam  = "metricType"
	metricNameParam  = "metricName"
	metricValueParam = "metricValue"
)

type metricService interface {
	UpdateMetric(ctx context.Context, metricType, name, rawValue string) error
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
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /update/{metricType}/{metricName}/{metricValue}", h.UpdateMetric)
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
