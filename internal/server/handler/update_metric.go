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

	m := &model.MetricsDto{
		ID:   metricID,
		Type: metricType,
	}

	switch metricType {
	case model.Gauge:
		gaugeValue, err := strconv.ParseFloat(metricValueRaw, 64)
		if err != nil {
			h.sendError(w, errs.Wrap(errs.ErrInvalidMetricValue, err.Error()))
			return
		}
		m.Value = &gaugeValue
	case model.Counter:
		counterValue, err := strconv.ParseInt(metricValueRaw, 10, 64)
		if err != nil {
			h.sendError(w, errs.Wrap(errs.ErrInvalidMetricValue, err.Error()))
			return
		}
		m.Delta = &counterValue
	}

	if err := h.ms.UpdateMetric(r.Context(), m); err != nil {
		h.sendError(w, errs.Wrap(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	err := h.audit.LogMetrics(r.Context(), []string{m.ID}, r.RemoteAddr)
	if err != nil {
		h.sendError(w, errs.Wrap(err))
		return
	}
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

	if err := h.ms.UpdateMetric(r.Context(), &req); err != nil {
		h.sendError(w, errs.Wrap(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	err = h.audit.LogMetrics(r.Context(), []string{req.ID}, r.RemoteAddr)
	if err != nil {
		h.sendError(w, errs.Wrap(err))
		return
	}
}

// UpdateMetricsBatch handles batch metric update requests.
func (h *Handler) UpdateMetricsBatch(w http.ResponseWriter, r *http.Request) {
	var req []*model.MetricsDto

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.sendError(w, errs.Wrap(err))
		return
	}

	if err := json.Unmarshal(body, &req); err != nil {
		h.sendError(w, errs.Wrap(errs.ErrInvalidJSON, err.Error()))
		return
	}

	metricsNames := make([]string, 0, len(req))
	for _, m := range req {
		if m.ID == "" {
			h.sendError(w, errs.Wrap(errs.ErrNoMetricID))
			return
		}
		metricsNames = append(metricsNames, m.ID)
	}

	err = h.ms.UpdateMetricsBatch(r.Context(), req)
	if err != nil {
		h.sendError(w, errs.Wrap(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	err = h.audit.LogMetrics(r.Context(), metricsNames, r.RemoteAddr)
	if err != nil {
		h.sendError(w, errs.Wrap(err))
		return
	}
}
