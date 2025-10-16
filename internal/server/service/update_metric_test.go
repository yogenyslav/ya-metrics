package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	models "github.com/yogenyslav/ya-metrics/internal/model"
)

type MockGaugeRepo struct {
	mock.Mock
}

func (m *MockGaugeRepo) Get(name string) (*models.Metrics[float64], bool) {
	args := m.Called(name)
	return args.Get(0).(*models.Metrics[float64]), args.Bool(1)
}

func (m *MockGaugeRepo) Set(name string, value float64) {
	m.Called(name, value)
}

type MockCounterRepo struct {
	mock.Mock
}

func (m *MockCounterRepo) Get(name string) (*models.Metrics[int64], bool) {
	args := m.Called(name)
	return args.Get(0).(*models.Metrics[int64]), args.Bool(1)
}

func (m *MockCounterRepo) Update(name string, delta int64) {
	m.Called(name, delta)
}

func TestService_UpdateMetric(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx        context.Context
		metricType string
		name       string
		rawValue   string
	}

	ctx := context.Background()

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Update existing gauge metric",
			args: args{
				metricType: models.Gauge,
				name:       "mem_alloc",
				rawValue:   "123.45",
			},
			wantErr: false,
		},
		{
			name: "Update non-existing gauge metric",
			args: args{
				metricType: models.Gauge,
				name:       "non_existing_gauge",
				rawValue:   "67.89",
			},
			wantErr: false,
		},
		{
			name: "Update existing counter metric",
			args: args{
				metricType: models.Counter,
				name:       "request_count",
				rawValue:   "10",
			},
			wantErr: false,
		},
		{
			name: "Update non-existing counter metric",
			args: args{
				metricType: models.Counter,
				name:       "non_existing_counter",
				rawValue:   "5",
			},
			wantErr: false,
		},
		{
			name: "Invalid gauge value",
			args: args{
				metricType: models.Gauge,
				name:       "mem_alloc",
				rawValue:   "invalid_float",
			},
			wantErr: true,
		},
		{
			name: "Invalid counter value",
			args: args{
				metricType: models.Counter,
				name:       "request_count",
				rawValue:   "invalid_int",
			},
			wantErr: true,
		},
		{
			name: "Invalid metric type",
			args: args{
				metricType: "invalid_type",
				name:       "some_metric",
				rawValue:   "100",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				t.Parallel()

				gr := &MockGaugeRepo{}
				cr := &MockCounterRepo{}

				s := Service{
					gr: gr,
					cr: cr,
				}

				switch tt.args.metricType {
				case models.Gauge:
					gr.On("Set", tt.args.name, mock.AnythingOfType("float64")).Return()
				case models.Counter:
					cr.On("Update", tt.args.name, mock.AnythingOfType("int64")).Return()
				}

				err := s.UpdateMetric(ctx, tt.args.metricType, tt.args.name, tt.args.rawValue)
				if tt.wantErr {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
				}
			},
		)
	}
}
