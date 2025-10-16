package handler

import (
	"net/http"

	"github.com/yogenyslav/ya-metrics/pkg/errs"
)

// UpdateMetric handles metric update requests.
func (h *Handler) UpdateMetric(w http.ResponseWriter, r *http.Request) {
	metricType := r.PathValue(metricTypeParam)
	metricName := r.PathValue(metricNameParam)
	metricValueRaw := r.PathValue(metricValueParam)

	if metricName == "" {
		h.sendError(w, errs.Wrap(errs.ErrNoMetricName))
		return
	}

	if err := h.ms.UpdateMetric(r.Context(), metricType, metricName, metricValueRaw); err != nil {
		h.sendError(w, errs.Wrap(err))
		return
	}

	w.WriteHeader(http.StatusOK)
}
