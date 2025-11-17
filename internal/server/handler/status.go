package handler

import (
	"net/http"

	"github.com/yogenyslav/ya-metrics/pkg/errs"
)

var errStatusCodes = map[error]int{
	errs.ErrInvalidMetricType:   http.StatusBadRequest,
	errs.ErrInvalidMetricValue:  http.StatusBadRequest,
	errs.ErrNoMetricID:          http.StatusNotFound,
	errs.ErrMetricNotFound:      http.StatusNotFound,
	errs.ErrInvalidJSON:         http.StatusUnprocessableEntity,
	errs.ErrDatabaseUnavailable: http.StatusInternalServerError,
}
