package service

import (
	"context"

	"github.com/yogenyslav/ya-metrics/internal/model"
)

// ListMetrics retrieves all gauge and counter metrics.
func (s *Service) ListMetrics(ctx context.Context) []*model.MetricsDto {
	gauges := s.gr.List()
	counters := s.cr.List()

	result := make([]*model.MetricsDto, 0, len(gauges)+len(counters))

	for _, g := range gauges {
		result = append(result, g.ToDto())
	}
	for _, c := range counters {
		result = append(result, c.ToDto())
	}

	return result
}
