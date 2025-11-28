package mocks

import (
	context "context"

	"github.com/stretchr/testify/mock"
	"github.com/yogenyslav/ya-metrics/internal/model"
)

type MockMetricService struct {
	mock.Mock
}

func (m *MockMetricService) UpdateMetric(ctx context.Context, metric *model.MetricsDto) error {
	args := m.Called(ctx, metric)
	m.ExpectedCalls = m.ExpectedCalls[1:]
	return args.Error(0)
}

func (m *MockMetricService) UpdateMetricsBatch(ctx context.Context, metrics []*model.MetricsDto) error {
	args := m.Called(ctx, metrics)
	m.ExpectedCalls = m.ExpectedCalls[1:]
	return args.Error(0)
}

func (m *MockMetricService) GetMetric(ctx context.Context, metricType, metricID string) (*model.MetricsDto, error) {
	args := m.Called(ctx, metricType, metricID)
	m.ExpectedCalls = m.ExpectedCalls[1:]
	return args.Get(0).(*model.MetricsDto), args.Error(1)
}

func (m *MockMetricService) ListMetrics(ctx context.Context) ([]*model.MetricsDto, error) {
	args := m.Called(ctx)
	m.ExpectedCalls = m.ExpectedCalls[1:]
	return args.Get(0).([]*model.MetricsDto), args.Error(1)
}
