package service

import (
	"context"

	"github.com/yogenyslav/ya-metrics/internal/model"
	"github.com/yogenyslav/ya-metrics/pkg/errs"
)

// UpdateMetric updates a metric with the given type, name, and raw value.
func (s *Service) UpdateMetric(ctx context.Context, req *model.MetricsDto) error {
	switch req.Type {
	case model.Counter:
		return s.cr.Update(ctx, &model.Metrics[int64]{
			ID:    req.ID,
			Type:  model.Counter,
			Value: *req.Delta,
		})
	case model.Gauge:
		return s.gr.Set(ctx, &model.Metrics[float64]{
			ID:    req.ID,
			Type:  model.Gauge,
			Value: *req.Value,
		})
	}
	return errs.Wrap(errs.ErrInvalidMetricType)
}

// UpdateMetricsBatch updates a batch of metrics.
func (s *Service) UpdateMetricsBatch(ctx context.Context, reqs []*model.MetricsDto) error {
	var (
		gauges   []*model.Metrics[float64]
		counters []*model.Metrics[int64]
	)

	for _, req := range reqs {
		switch req.Type {
		case model.Gauge:
			gauges = append(gauges, &model.Metrics[float64]{
				ID:    req.ID,
				Type:  model.Gauge,
				Value: *req.Value,
			})
		case model.Counter:
			counters = append(counters, &model.Metrics[int64]{
				ID:    req.ID,
				Type:  model.Counter,
				Value: *req.Delta,
			})
		default:
			return errs.Wrap(errs.ErrInvalidMetricType)
		}
	}

	if err := s.gr.SetBatch(ctx, gauges); err != nil {
		return errs.Wrap(err, "batch set gauge metrics")
	}

	if err := s.cr.UpdateBatch(ctx, counters); err != nil {
		return errs.Wrap(err, "batch update counter metrics")
	}

	return nil
}
