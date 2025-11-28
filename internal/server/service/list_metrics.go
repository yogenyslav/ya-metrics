package service

import (
	"context"

	"github.com/yogenyslav/ya-metrics/internal/model"
	"github.com/yogenyslav/ya-metrics/pkg/errs"
)

// ListMetrics retrieves all gauge and counter metrics.
func (s *Service) ListMetrics(ctx context.Context) ([]*model.MetricsDto, error) {
	gauges, err := s.gr.List(ctx)
	if err != nil {
		return nil, errs.Wrap(err, "list gauge metrics")
	}
	counters, err := s.cr.List(ctx)
	if err != nil {
		return nil, errs.Wrap(err, "list counter metrics")
	}

	result := make([]*model.MetricsDto, 0, len(gauges)+len(counters))

	for _, g := range gauges {
		result = append(result, g.ToDto())
	}
	for _, c := range counters {
		result = append(result, c.ToDto())
	}

	return result, nil
}
