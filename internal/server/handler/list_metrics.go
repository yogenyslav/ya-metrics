package handler

import (
	"encoding/json"
	"net/http"

	"github.com/yogenyslav/ya-metrics/pkg/errs"
)

// ListMetrics handles requests to list all metrics.
func (h *Handler) ListMetrics(w http.ResponseWriter, r *http.Request) {
	metrics, err := h.ms.ListMetrics(r.Context())
	if err != nil {
		h.sendError(w, errs.Wrap(err))
		return
	}

	resp, err := json.Marshal(metrics)
	if err != nil {
		h.sendError(w, errs.Wrap(err))
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
