package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yogenyslav/ya-metrics/internal/model"
	"github.com/yogenyslav/ya-metrics/pkg"
	"github.com/yogenyslav/ya-metrics/tests/mocks"
	gomock "go.uber.org/mock/gomock"
)

func TestMetricPostgresRepo_Get(t *testing.T) {
	t.Parallel()

	type testCase[T int64 | float64] struct {
		name       string
		db         func() *mocks.MockDB
		metricName string
		wantMetric *model.Metrics[T]
		wantErr    bool
	}

	t.Run("int64", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		tests := []testCase[int64]{
			{
				name: "Get existing int64 metric",
				db: func() *mocks.MockDB {
					mockDB := mocks.NewMockDB(ctrl)
					mockDB.EXPECT().
						QueryRow(gomock.Any(), gomock.Any(), getCounterMetric, gomock.Any(), gomock.Any()).
						DoAndReturn(func(ctx context.Context, dest any, query string, args ...any) error {
							d := dest.(*model.MetricsDto)
							d.ID = "metric1"
							d.Type = model.Counter
							d.Delta = pkg.Ptr(int64(10))
							return nil
						})
					return mockDB
				},
				metricName: "metric1",
				wantMetric: &model.Metrics[int64]{ID: "metric1", Type: model.Counter, Value: 10},
				wantErr:    false,
			},
			{
				name: "Get non-existing int64 metric",
				db: func() *mocks.MockDB {
					mockDB := mocks.NewMockDB(ctrl)
					mockDB.EXPECT().
						QueryRow(gomock.Any(), gomock.Any(), getCounterMetric, gomock.Any(), gomock.Any()).
						Return(errors.New("not found"))
					return mockDB
				},
				metricName: "metric2",
				wantMetric: nil,
				wantErr:    true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				repo := NewMetricPostgresRepo[int64](tt.db())
				got, err := repo.Get(context.Background(), tt.metricName, model.Counter)

				if tt.wantErr {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
					assert.Equal(t, tt.wantMetric, got)
				}
			})
		}
	})

	t.Run("float64", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		tests := []testCase[float64]{
			{
				name: "Get existing float64 metric",
				db: func() *mocks.MockDB {
					mockDB := mocks.NewMockDB(ctrl)
					mockDB.EXPECT().
						QueryRow(gomock.Any(), gomock.Any(), getGaugeMetric, gomock.Any(), gomock.Any()).
						DoAndReturn(func(ctx context.Context, dest any, query string, args ...any) error {
							d := dest.(*model.MetricsDto)
							d.ID = "metric1"
							d.Type = model.Gauge
							val := 123.45
							d.Value = &val
							return nil
						})
					return mockDB
				},
				metricName: "metric1",
				wantMetric: &model.Metrics[float64]{ID: "metric1", Type: model.Gauge, Value: 123.45},
				wantErr:    false,
			},
			{
				name: "Get non-existing float64 metric",
				db: func() *mocks.MockDB {
					mockDB := mocks.NewMockDB(ctrl)
					mockDB.EXPECT().
						QueryRow(gomock.Any(), gomock.Any(), getGaugeMetric, gomock.Any(), gomock.Any()).
						Return(errors.New("not found"))
					return mockDB
				},
				metricName: "metric2",
				wantMetric: nil,
				wantErr:    true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				repo := NewMetricPostgresRepo[float64](tt.db())
				got, err := repo.Get(context.Background(), tt.metricName, model.Gauge)

				if tt.wantErr {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
					assert.Equal(t, tt.wantMetric, got)
				}
			})
		}
	})
}

func TestMetricPostgresRepo_List(t *testing.T) {
	t.Parallel()

	type testCase[T int64 | float64] struct {
		name    string
		db      func() *mocks.MockDB
		want    []model.Metrics[T]
		wantErr bool
	}

	t.Run("int64", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		tests := []testCase[int64]{
			{
				name: "List int64 metrics",
				db: func() *mocks.MockDB {
					mockDB := mocks.NewMockDB(ctrl)
					mockDB.EXPECT().
						QuerySlice(gomock.Any(), gomock.Any(), listCounterMetrics, gomock.Any()).
						DoAndReturn(func(ctx context.Context, dest any, query string, args ...any) error {
							d := dest.(*[]model.Metrics[int64])
							*d = []model.Metrics[int64]{
								{ID: "metric1", Type: model.Counter, Value: 10},
								{ID: "metric2", Type: model.Counter, Value: 20},
							}
							return nil
						})
					return mockDB
				},
				want: []model.Metrics[int64]{
					{ID: "metric1", Type: model.Counter, Value: 10},
					{ID: "metric2", Type: model.Counter, Value: 20},
				},
				wantErr: false,
			},
			{
				name: "List from empty int64 repo",
				db: func() *mocks.MockDB {
					mockDB := mocks.NewMockDB(ctrl)
					mockDB.EXPECT().
						QuerySlice(gomock.Any(), gomock.Any(), listCounterMetrics, gomock.Any()).
						DoAndReturn(func(ctx context.Context, dest any, query string, args ...any) error {
							d := dest.(*[]model.Metrics[int64])
							*d = []model.Metrics[int64]{}
							return nil
						})
					return mockDB
				},
				want:    []model.Metrics[int64]{},
				wantErr: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				repo := NewMetricPostgresRepo[int64](tt.db())
				got, err := repo.List(context.Background())

				if tt.wantErr {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
					assert.Equal(t, tt.want, got)
				}
			})
		}
	})

	t.Run("float64", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		tests := []testCase[float64]{
			{
				name: "List float64 metrics",
				db: func() *mocks.MockDB {
					mockDB := mocks.NewMockDB(ctrl)
					mockDB.EXPECT().
						QuerySlice(gomock.Any(), gomock.Any(), listGaugeMetrics, gomock.Any()).
						DoAndReturn(func(ctx context.Context, dest any, query string, args ...any) error {
							d := dest.(*[]model.Metrics[float64])
							*d = []model.Metrics[float64]{
								{ID: "metric1", Type: model.Gauge, Value: 10.5},
								{ID: "metric2", Type: model.Gauge, Value: 20.3},
							}
							return nil
						})
					return mockDB
				},
				want: []model.Metrics[float64]{
					{ID: "metric1", Type: model.Gauge, Value: 10.5},
					{ID: "metric2", Type: model.Gauge, Value: 20.3},
				},
				wantErr: false,
			},
			{
				name: "List from empty float64 repo",
				db: func() *mocks.MockDB {
					mockDB := mocks.NewMockDB(ctrl)
					mockDB.EXPECT().
						QuerySlice(gomock.Any(), gomock.Any(), listGaugeMetrics, gomock.Any()).
						DoAndReturn(func(ctx context.Context, dest any, query string, args ...any) error {
							d := dest.(*[]model.Metrics[float64])
							*d = []model.Metrics[float64]{}
							return nil
						})
					return mockDB
				},
				want:    []model.Metrics[float64]{},
				wantErr: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				repo := NewMetricPostgresRepo[float64](tt.db())
				got, err := repo.List(context.Background())

				if tt.wantErr {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
					assert.Equal(t, tt.want, got)
				}
			})
		}
	})
}

func TestMetricPostgresRepo_Set(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name      string
		db        func() *mocks.MockDB
		metric    *model.Metrics[float64]
		wantError bool
	}{
		{
			name: "Set gauge metric successfully",
			db: func() *mocks.MockDB {
				mockDB := mocks.NewMockDB(ctrl)
				mockDB.EXPECT().
					Exec(gomock.Any(), setGaugeMetric, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(int64(1), nil)
				return mockDB
			},
			metric:    &model.Metrics[float64]{ID: "metric1", Type: model.Gauge, Value: 123.45},
			wantError: false,
		},
		{
			name: "Set gauge metric with DB error",
			db: func() *mocks.MockDB {
				mockDB := mocks.NewMockDB(ctrl)
				mockDB.EXPECT().
					Exec(gomock.Any(), setGaugeMetric, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(int64(0), errors.New("db error"))
				return mockDB
			},
			metric:    &model.Metrics[float64]{ID: "metric2", Type: model.Gauge, Value: 67.89},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := NewMetricPostgresRepo[float64](tt.db())
			err := repo.Set(context.Background(), tt.metric)

			if tt.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestMetricPostgresRepo_Update(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name      string
		db        func() *mocks.MockDB
		metric    *model.Metrics[int64]
		wantError bool
	}{
		{
			name: "Update counter metric successfully",
			db: func() *mocks.MockDB {
				mockDB := mocks.NewMockDB(ctrl)
				mockDB.EXPECT().
					Exec(gomock.Any(), updateCounterMetric, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(int64(1), nil)
				return mockDB
			},
			metric:    &model.Metrics[int64]{ID: "metric1", Type: model.Counter, Value: 10},
			wantError: false,
		},
		{
			name: "Update counter metric with DB error",
			db: func() *mocks.MockDB {
				mockDB := mocks.NewMockDB(ctrl)
				mockDB.EXPECT().
					Exec(gomock.Any(), updateCounterMetric, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(int64(0), errors.New("db error"))
				return mockDB
			},
			metric:    &model.Metrics[int64]{ID: "metric2", Type: model.Counter, Value: 20},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := NewMetricPostgresRepo[int64](tt.db())
			err := repo.Update(context.Background(), tt.metric)

			if tt.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
