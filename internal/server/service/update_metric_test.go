package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/yogenyslav/ya-metrics/internal/model"
	"github.com/yogenyslav/ya-metrics/pkg"
	"github.com/yogenyslav/ya-metrics/pkg/errs"
	"github.com/yogenyslav/ya-metrics/tests/mocks"
)

func TestService_UpdateMetric(t *testing.T) {
	t.Parallel()

	type args struct {
		req *model.MetricsDto
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
				req: &model.MetricsDto{
					ID:    "mem_alloc",
					Type:  model.Gauge,
					Value: pkg.Ptr(123.45),
					Delta: (*int64)(nil),
				},
			},
			wantErr: false,
		},
		{
			name: "Update non-existing gauge metric",
			args: args{
				req: &model.MetricsDto{
					ID:    "non_existing_gauge",
					Type:  model.Gauge,
					Value: pkg.Ptr(67.89),
					Delta: (*int64)(nil),
				},
			},
			wantErr: false,
		},
		{
			name: "Update existing counter metric",
			args: args{
				req: &model.MetricsDto{
					ID:    "request_count",
					Type:  model.Counter,
					Value: (*float64)(nil),
					Delta: pkg.Ptr(int64(10)),
				},
			},
			wantErr: false,
		},
		{
			name: "Update non-existing counter metric",
			args: args{
				req: &model.MetricsDto{
					ID:    "non_existing_counter",
					Type:  model.Counter,
					Value: (*float64)(nil),
					Delta: pkg.Ptr(int64(5)),
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid metric type",
			args: args{
				req: &model.MetricsDto{
					ID:    "some_metric",
					Type:  "invalid_type",
					Value: new(float64),
					Delta: new(int64),
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				t.Parallel()

				gr := &mocks.MockGaugeRepo{}
				cr := &mocks.MockCounterRepo{}

				s := Service{
					gr: gr,
					cr: cr,
				}

				switch tt.args.req.Type {
				case model.Gauge:
					if !tt.wantErr {
						gr.On("Set", mock.Anything, &model.Metrics[float64]{
							ID:    tt.args.req.ID,
							Type:  model.Gauge,
							Value: *tt.args.req.Value,
						}).Return(nil)
					} else {
						gr.On("Set", mock.Anything, &model.Metrics[float64]{
							ID:   tt.args.req.ID,
							Type: model.Gauge,
						}).Return(errs.ErrInvalidMetricValue)
					}
				case model.Counter:
					if !tt.wantErr {
						cr.On("Update", mock.Anything, &model.Metrics[int64]{
							ID:    tt.args.req.ID,
							Type:  model.Counter,
							Value: *tt.args.req.Delta,
						}).Return(nil)
					} else {
						cr.On("Update", mock.Anything, &model.Metrics[int64]{
							ID:   tt.args.req.ID,
							Type: model.Counter,
						}).Return(errs.ErrInvalidMetricValue)
					}
				}

				err := s.UpdateMetric(ctx, tt.args.req)
				if tt.wantErr {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
				}
			},
		)
	}
}
