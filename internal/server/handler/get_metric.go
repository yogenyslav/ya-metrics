package handler

import (
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

	w.Write([]byte(metric.Value))
}
