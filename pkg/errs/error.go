package errs

import (
	"errors"
)

// 400.
var (
	// ErrInvalidMetricType is an error when invalid metric type is provided.
	ErrInvalidMetricType = errors.New("invalid metric type")
	// ErrInvalidMetricValue is an error when failed to parse metric value.
	ErrInvalidMetricValue = errors.New("invalid metric value")
)

// 404.
var (
	// ErrNoMetricName is an error when no metric name is provided.
	ErrNoMetricName = errors.New("no metric name provided")
	// ErrMetricNotFound is an error when the requested metric is not found.
	ErrMetricNotFound = errors.New("metric not found")
)

// 422.
var (
	// ErrInvalidJSON is an error when the provided JSON is invalid.
	ErrInvalidJSON = errors.New("invalid JSON")
)
