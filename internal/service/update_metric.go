package service

import (
	"context"
	"strconv"

	models "github.com/yogenyslav/ya-metrics/internal/model"
	"github.com/yogenyslav/ya-metrics/pkg/errs"
)

// UpdateMetric updates a metric with the given type, name, and raw value.
func (s *Service) UpdateMetric(ctx context.Context, metricType, name, rawValue string) error {
	switch metricType {
	case models.Counter:
		metricValue, err := strconv.ParseInt(rawValue, 10, 64)
		if err != nil {
			return errs.Wrap(errs.ErrInvalidMetricValue, err.Error())
		}
		return s.UpdateCounter(ctx, name, metricValue)
	case models.Gauge:
		metricValue, err := strconv.ParseFloat(rawValue, 64)
		if err != nil {
			return errs.Wrap(errs.ErrInvalidMetricValue, err.Error())
		}
		return s.UpdateGauge(ctx, name, metricValue)
	}
	return errs.Wrap(errs.ErrInvalidMetricType)
}
