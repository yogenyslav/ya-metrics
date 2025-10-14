package service

import (
	"context"
)

// UpdateCounter updates the counter metric by the given delta.
func (s *Service) UpdateCounter(_ context.Context, name string, delta int64) error {
	s.counterRepo.Update(name, delta)
	return nil
}
