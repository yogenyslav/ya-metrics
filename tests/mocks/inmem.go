package mocks

import (
	"github.com/stretchr/testify/mock"
	"github.com/yogenyslav/ya-metrics/internal/model"
)

type MockMetricRepo[T int64 | float64] struct {
	mock.Mock
}

func (m *MockMetricRepo[T]) GetMetrics() []*model.MetricsDto {
	args := m.Called()
	return args.Get(0).([]*model.MetricsDto)
}

func (m *MockMetricRepo[T]) Get(name string) (*model.Metrics[T], bool) {
	args := m.Called(name)
	return args.Get(0).(*model.Metrics[T]), args.Bool(1)
}

func (m *MockMetricRepo[T]) List() []model.Metrics[T] {
	args := m.Called()
	return args.Get(0).([]model.Metrics[T])
}

type MockGaugeRepo struct {
	MockMetricRepo[float64]
}

func (m *MockGaugeRepo) Set(name string, value float64, tp string) {
	m.Called(name, value)
}

type MockCounterRepo struct {
	MockMetricRepo[int64]
}

func (m *MockCounterRepo) Update(name string, delta int64, tp string) {
	m.Called(name, delta)
}
