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
	err := s.uow.WithTx(ctx, func(ctx context.Context) error {
		var err error
		for _, req := range reqs {
			err = s.UpdateMetric(ctx, req)
			if err != nil {
				return errs.Wrap(err, "update metric in tx")
			}
		}
		return nil
	})
	return errs.Wrap(err, "update metrics batch")
}
