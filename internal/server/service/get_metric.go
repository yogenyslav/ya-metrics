package service

import (
	"context"

	"github.com/yogenyslav/ya-metrics/internal/model"
)

// GetMetric retrieves a metric by its type and name.
func (s *Service) GetMetric(ctx context.Context, metricType, metricName string) (*model.MetricsDto, bool) {
	switch metricType {
	case model.Gauge:
		gauge, found := s.gr.Get(metricName)
		if !found {
			return nil, false
		}
		return gauge.ToDto(), true
	case model.Counter:
		counter, found := s.cr.Get(metricName)
		if !found {
			return nil, false
		}
		return counter.ToDto(), true
	default:
		return nil, false
	}
}
