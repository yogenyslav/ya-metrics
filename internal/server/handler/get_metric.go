package handler

import (
	"encoding/json"
	"net/http"

	"github.com/yogenyslav/ya-metrics/pkg/errs"
)

// GetMetric handles metric retrieval requests.
func (h *Handler) GetMetric(w http.ResponseWriter, r *http.Request) {
	metricType := r.PathValue(metricTypeParam)
	metricName := r.PathValue(metricNameParam)

	metric, found := h.ms.GetMetric(r.Context(), metricType, metricName)
	if !found {
		h.sendError(w, errs.Wrap(errs.ErrMetricNotFound))
		return
	}

	resp, err := json.Marshal(metric)
	if err != nil {
		h.sendError(w, errs.Wrap(err))
		return
	}

	w.Write(resp)
}
