package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/yogenyslav/ya-metrics/internal/model"
)

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
				metricType: model.Gauge,
				name:       "mem_alloc",
				rawValue:   "123.45",
			},
			wantErr: false,
		},
		{
			name: "Update non-existing gauge metric",
			args: args{
				metricType: model.Gauge,
				name:       "non_existing_gauge",
				rawValue:   "67.89",
			},
			wantErr: false,
		},
		{
			name: "Update existing counter metric",
			args: args{
				metricType: model.Counter,
				name:       "request_count",
				rawValue:   "10",
			},
			wantErr: false,
		},
		{
			name: "Update non-existing counter metric",
			args: args{
				metricType: model.Counter,
				name:       "non_existing_counter",
				rawValue:   "5",
			},
			wantErr: false,
		},
		{
			name: "Invalid gauge value",
			args: args{
				metricType: model.Gauge,
				name:       "mem_alloc",
				rawValue:   "invalid_float",
			},
			wantErr: true,
		},
		{
			name: "Invalid counter value",
			args: args{
				metricType: model.Counter,
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
				case model.Gauge:
					gr.On("Set", tt.args.name, mock.AnythingOfType("float64")).Return()
				case model.Counter:
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
