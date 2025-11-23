package service

import (
	"context"

	"github.com/yogenyslav/ya-metrics/internal/model"
	"github.com/yogenyslav/ya-metrics/pkg/errs"
)

// GetMetric retrieves a metric by its type and name.
func (s *Service) GetMetric(ctx context.Context, metricType, metricID string) (*model.MetricsDto, error) {
	switch metricType {
	case model.Gauge:
		gauge, err := s.gr.Get(ctx, metricID, metricType)
		if err != nil {
			return nil, errs.Wrap(err)
		}
		return gauge.ToDto(), nil
	case model.Counter:
		counter, err := s.cr.Get(ctx, metricID, metricType)
		if err != nil {
			return nil, errs.Wrap(err)
		}
		return counter.ToDto(), nil
	default:
		return nil, errs.Wrap(errs.ErrInvalidMetricType, metricType)
	}
}
