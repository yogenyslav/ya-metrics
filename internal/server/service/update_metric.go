package service

import (
	"context"
	"strconv"

	"github.com/yogenyslav/ya-metrics/internal/model"
	"github.com/yogenyslav/ya-metrics/pkg/errs"
)

// UpdateMetric updates a metric with the given type, name, and raw value.
func (s *Service) UpdateMetric(ctx context.Context, metricType, name, rawValue string) error {
	switch metricType {
	case model.Counter:
		metricValue, err := strconv.ParseInt(rawValue, 10, 64)
		if err != nil {
			return errs.Wrap(errs.ErrInvalidMetricValue, err.Error())
		}
		s.cr.Update(name, metricValue, metricType)
		return nil
	case model.Gauge:
		metricValue, err := strconv.ParseFloat(rawValue, 64)
		if err != nil {
			return errs.Wrap(errs.ErrInvalidMetricValue, err.Error())
		}
		s.gr.Set(name, metricValue, metricType)
		return nil
	}
	return errs.Wrap(errs.ErrInvalidMetricType)
}
