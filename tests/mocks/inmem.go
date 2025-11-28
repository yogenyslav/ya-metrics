package mocks

import (
	context "context"

	"github.com/stretchr/testify/mock"
	"github.com/yogenyslav/ya-metrics/internal/model"
)

type MockMetricRepo[T int64 | float64] struct {
	mock.Mock
}

func (m *MockMetricRepo[T]) GetMetrics(ctx context.Context) ([]*model.MetricsDto, error) {
	args := m.Called(ctx)
	m.ExpectedCalls = m.ExpectedCalls[1:]
	return args.Get(0).([]*model.MetricsDto), args.Error(1)
}

func (m *MockMetricRepo[T]) Get(ctx context.Context, metricName, metricType string) (*model.Metrics[T], error) {
	args := m.Called(ctx, metricName, metricType)
	m.ExpectedCalls = m.ExpectedCalls[1:]
	return args.Get(0).(*model.Metrics[T]), args.Error(1)
}

func (m *MockMetricRepo[T]) List(ctx context.Context) ([]model.Metrics[T], error) {
	args := m.Called(ctx)
	m.ExpectedCalls = m.ExpectedCalls[1:]
	return args.Get(0).([]model.Metrics[T]), args.Error(1)
}

type MockGaugeRepo struct {
	MockMetricRepo[float64]
}

func (m *MockGaugeRepo) Set(ctx context.Context, metric *model.Metrics[float64]) error {
	args := m.Called(ctx, metric)
	m.ExpectedCalls = m.ExpectedCalls[1:]
	return args.Error(0)
}

func (m *MockGaugeRepo) SetBatch(ctx context.Context, metrics []*model.Metrics[float64]) error {
	args := m.Called(ctx, metrics)
	m.ExpectedCalls = m.ExpectedCalls[1:]
	return args.Error(0)
}

type MockCounterRepo struct {
	MockMetricRepo[int64]
}

func (m *MockCounterRepo) Update(ctx context.Context, metric *model.Metrics[int64]) error {
	args := m.Called(ctx, metric)
	m.ExpectedCalls = m.ExpectedCalls[1:]
	return args.Error(0)
}

func (m *MockCounterRepo) UpdateBatch(ctx context.Context, metrics []*model.Metrics[int64]) error {
	args := m.Called(ctx, metrics)
	m.ExpectedCalls = m.ExpectedCalls[1:]
	return args.Error(0)
}
