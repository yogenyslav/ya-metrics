package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/yogenyslav/ya-metrics/internal/model"
	"github.com/yogenyslav/ya-metrics/pkg/errs"
)

// GetMetricRaw handles raw metric retrieval requests.
func (h *Handler) GetMetricRaw(w http.ResponseWriter, r *http.Request) {
	metricType := r.PathValue(metricTypeParam)
	metricID := r.PathValue(metricIDParam)

	metric, found := h.ms.GetMetric(r.Context(), metricType, metricID)
	if !found {
		h.sendError(w, errs.Wrap(errs.ErrMetricNotFound))
		return
	}

	var value string
	if metricType == model.Gauge {
		value = strconv.FormatFloat(*metric.Value, 'f', -1, 64)
	} else {
		value = strconv.FormatInt(*metric.Delta, 10)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(value))
}

// GetMetricJSON handles JSON metric retrieval requests.
func (h *Handler) GetMetricJSON(w http.ResponseWriter, r *http.Request) {
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

	metric, found := h.ms.GetMetric(r.Context(), req.Type, req.ID)
	if !found {
		h.sendError(w, errs.Wrap(errs.ErrMetricNotFound))
		return
	}

	resp := model.MetricsDto{
		ID:    metric.ID,
		Type:  metric.Type,
		Value: metric.Value,
		Delta: metric.Delta,
	}

	respBody, err := json.Marshal(resp)
	if err != nil {
		h.sendError(w, errs.Wrap(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respBody)
}
