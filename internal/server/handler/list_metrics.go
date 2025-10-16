package handler

import (
	"encoding/json"
	"net/http"

	"github.com/yogenyslav/ya-metrics/pkg/errs"
)

// ListMetrics handles requests to list all metrics.
func (h *Handler) ListMetrics(w http.ResponseWriter, r *http.Request) {
	metrics := h.ms.ListMetrics(r.Context())

	resp, err := json.Marshal(metrics)
	if err != nil {
		h.sendError(w, errs.Wrap(err))
		return
	}

	w.Write(resp)
}
