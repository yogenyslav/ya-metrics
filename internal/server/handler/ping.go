package handler

import (
	"net/http"

	"github.com/yogenyslav/ya-metrics/pkg/errs"
)

// Ping checks the database connectivity.
func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	err := h.db.Ping(r.Context())
	if err != nil {
		h.sendError(w, errs.Wrap(errs.ErrDatabaseUnavailable, err.Error()))
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
}
