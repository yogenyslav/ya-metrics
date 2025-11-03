package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/yogenyslav/ya-metrics/internal/model"
	"github.com/yogenyslav/ya-metrics/pkg/errs"
)

// UpdateMetricRaw handles raw metric update requests.
func (h *Handler) UpdateMetricRaw(w http.ResponseWriter, r *http.Request) {
	metricType := r.PathValue(metricTypeParam)
	metricID := r.PathValue(metricIDParam)
	metricValueRaw := r.PathValue(metricValueParam)

	if metricID == "" {
		h.sendError(w, errs.Wrap(errs.ErrNoMetricID))
		return
	}

	if err := h.ms.UpdateMetric(r.Context(), metricType, metricID, metricValueRaw); err != nil {
		h.sendError(w, errs.Wrap(err))
		return
	}

	w.WriteHeader(http.StatusOK)
}

// UpdateMetricJSON handles JSON metric update requests.
func (h *Handler) UpdateMetricJSON(w http.ResponseWriter, r *http.Request) {
	var req model.MetricsDto

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.sendError(w, errs.Wrap(err))
		return
	}

	if err := json.Unmarshal(body, &req); err != nil {
		h.sendError(w, errs.Wrap(errs.ErrInvalidJSON, err.Error()))
		return
	}

	if req.ID == "" {
		h.sendError(w, errs.Wrap(errs.ErrNoMetricID))
		return
	}

	var metricRawValue string
	switch {
	case req.Type == model.Gauge && req.Value != nil:
		metricRawValue = strconv.FormatFloat(*req.Value, 'f', -1, 64)
	case req.Type == model.Counter && req.Delta != nil:
		metricRawValue = strconv.FormatInt(*req.Delta, 10)
	}

	if err := h.ms.UpdateMetric(r.Context(), req.Type, req.ID, metricRawValue); err != nil {
		h.sendError(w, errs.Wrap(err))
		return
	}

	w.WriteHeader(http.StatusOK)
}
