package service

import (
	"context"
)

// UpdateGauge updates the gauge metric to the given value.
func (s *Service) UpdateGauge(_ context.Context, name string, value float64) error {
	s.gr.Set(name, value)
	return nil
}
