package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yogenyslav/ya-metrics/internal/model"
	"github.com/yogenyslav/ya-metrics/pkg"
)

func TestMetricInMemRepo_Get(t *testing.T) {
	t.Parallel()

	type args struct {
		name string
	}
	type testCase[T interface{ int64 | float64 }] struct {
		name       string
		r          MetricInMemRepo[T]
		args       args
		wantMetric *model.Metrics[T]
		wantExist  bool
	}

	t.Run(
		"int64", func(t *testing.T) {
			t.Parallel()

			tests := []testCase[int64]{
				{
					name: "Get existing int64 metric",
					r: MetricInMemRepo[int64]{
						storage: StorageState[int64]{
							"metric1": {ID: "metric1", Value: 10},
						},
					},
					args:       args{name: "metric1"},
					wantMetric: &model.Metrics[int64]{ID: "metric1", Value: 10},
					wantExist:  true,
				},
				{
					name:       "Get non-existing int64 metric",
					r:          MetricInMemRepo[int64]{storage: make(StorageState[int64])},
					args:       args{name: "metric2"},
					wantMetric: nil,
					wantExist:  false,
				},
			}

			for _, tt := range tests {
				t.Run(
					tt.name, func(t *testing.T) {
						t.Parallel()
						metric, exist := tt.r.Get(tt.args.name)
						assert.Equal(t, tt.wantMetric, metric)
						assert.Equal(t, tt.wantExist, exist)
					},
				)
			}
		},
	)

	t.Run(
		"float64", func(t *testing.T) {
			t.Parallel()

			tests := []testCase[float64]{
				{
					name: "Get existing float64 metric",
					r: MetricInMemRepo[float64]{
						storage: StorageState[float64]{
							"metric1": {ID: "metric1", Value: 10.5},
						},
					},
					args:       args{name: "metric1"},
					wantMetric: &model.Metrics[float64]{ID: "metric1", Value: 10.5},
					wantExist:  true,
				},
				{
					name:       "Get non-existing float64 metric",
					r:          MetricInMemRepo[float64]{storage: make(StorageState[float64])},
					args:       args{name: "metric2"},
					wantMetric: nil,
					wantExist:  false,
				},
			}

			for _, tt := range tests {
				t.Run(
					tt.name, func(t *testing.T) {
						t.Parallel()
						got, got1 := tt.r.Get(tt.args.name)
						assert.Equal(t, tt.wantMetric, got)
						assert.Equal(t, tt.wantExist, got1)
					},
				)
			}
		},
	)
}

func TestMetricInMemRepo_Set(t *testing.T) {
	t.Parallel()

	type args[T interface{ int64 | float64 }] struct {
		name  string
		value T
	}
	type testCase[T interface{ int64 | float64 }] struct {
		name string
		r    MetricInMemRepo[T]
		args args[T]
	}

	t.Run(
		"int64", func(t *testing.T) {
			t.Parallel()

			tests := []testCase[int64]{
				{
					name: "Set int64 metric",
					r:    MetricInMemRepo[int64]{storage: make(StorageState[int64])},
					args: args[int64]{name: "metric1", value: 10},
				},
			}

			for _, tt := range tests {
				t.Run(
					tt.name, func(t *testing.T) {
						t.Parallel()
						tt.r.Set(tt.args.name, tt.args.value, model.Counter)
						got, exists := tt.r.Get(tt.args.name)
						assert.True(t, exists)
						assert.Equal(t, tt.args.value, got.Value)
					},
				)
			}
		},
	)

	t.Run(
		"float64", func(t *testing.T) {
			t.Parallel()

			tests := []testCase[float64]{
				{
					name: "Set float64 metric",
					r:    MetricInMemRepo[float64]{storage: make(StorageState[float64])},
					args: args[float64]{name: "metric1", value: 10.5},
				},
			}

			for _, tt := range tests {
				t.Run(
					tt.name, func(t *testing.T) {
						t.Parallel()
						tt.r.Set(tt.args.name, tt.args.value, model.Gauge)
						got, exists := tt.r.Get(tt.args.name)
						assert.True(t, exists)
						assert.Equal(t, tt.args.value, got.Value)
					},
				)
			}
		},
	)
}

func TestMetricInMemRepo_Update(t *testing.T) {
	t.Parallel()

	type args[T interface{ int64 | float64 }] struct {
		name  string
		delta T
	}
	type testCase[T interface{ int64 | float64 }] struct {
		name string
		r    MetricInMemRepo[T]
		args args[T]
	}

	t.Run(
		"int64", func(t *testing.T) {
			t.Parallel()

			tests := []testCase[int64]{
				{
					name: "Update existing int64 metric",
					r: MetricInMemRepo[int64]{
						storage: StorageState[int64]{
							"metric1": {ID: "metric1", Value: 10},
						},
					},
					args: args[int64]{name: "metric1", delta: 5},
				},
				{
					name: "Update non-existing int64 metric",
					r:    MetricInMemRepo[int64]{storage: make(StorageState[int64])},
					args: args[int64]{name: "metric2", delta: 7},
				},
			}

			for _, tt := range tests {
				t.Run(
					tt.name, func(t *testing.T) {
						t.Parallel()
						tt.r.Update(tt.args.name, tt.args.delta, model.Counter)
						got, exists := tt.r.Get(tt.args.name)
						assert.True(t, exists)
						expectedValue := tt.args.delta
						if original, ok := tt.r.storage[tt.args.name]; ok {
							expectedValue += original.Value - tt.args.delta
						}
						assert.Equal(t, expectedValue, got.Value)
					},
				)
			}
		},
	)

	t.Run(
		"float64", func(t *testing.T) {
			t.Parallel()

			tests := []testCase[float64]{
				{
					name: "Update existing float64 metric",
					r: MetricInMemRepo[float64]{
						storage: StorageState[float64]{
							"metric1": {ID: "metric1", Value: 10.5},
						},
					},
					args: args[float64]{name: "metric1", delta: 5.5},
				},
				{
					name: "Update non-existing float64 metric",
					r:    MetricInMemRepo[float64]{storage: make(StorageState[float64])},
					args: args[float64]{name: "metric2", delta: 7.3},
				},
			}

			for _, tt := range tests {
				t.Run(
					tt.name, func(t *testing.T) {
						t.Parallel()
						tt.r.Update(tt.args.name, tt.args.delta, model.Gauge)
						got, exists := tt.r.Get(tt.args.name)
						assert.True(t, exists)
						expectedValue := tt.args.delta
						if original, ok := tt.r.storage[tt.args.name]; ok {
							expectedValue += original.Value - tt.args.delta
						}
						assert.Equal(t, expectedValue, got.Value)
					},
				)
			}
		},
	)
}

func TestNewMetricInMemRepo(t *testing.T) {
	t.Parallel()

	type testCase[T interface{ int64 | float64 }] struct {
		name  string
		state StorageState[T]
		want  *MetricInMemRepo[T]
	}

	t.Run(
		"int64", func(t *testing.T) {
			t.Parallel()

			tests := []testCase[int64]{
				{
					name:  "Create new int64 MetricInMemRepo",
					state: nil,
					want:  &MetricInMemRepo[int64]{storage: make(StorageState[int64])},
				},
				{
					name: "Create new int64 MetricInMemRepo with initial state",
					state: StorageState[int64]{
						"metric1": {ID: "metric1", Value: 10},
					},
					want: &MetricInMemRepo[int64]{storage: StorageState[int64]{
						"metric1": {ID: "metric1", Value: 10},
					}},
				},
			}

			for _, tt := range tests {
				t.Run(
					tt.name, func(t *testing.T) {
						t.Parallel()
						repo := NewMetricInMemRepo[int64](tt.state)
						assert.Equal(t, *tt.want, *repo)
					},
				)
			}
		},
	)

	t.Run(
		"float64", func(t *testing.T) {
			t.Parallel()

			tests := []testCase[float64]{
				{
					name:  "Create new float64 MetricInMemRepo",
					state: nil,
					want:  &MetricInMemRepo[float64]{storage: make(StorageState[float64])},
				},
				{
					name: "Create new float64 MetricInMemRepo with initial state",
					state: StorageState[float64]{
						"metric1": {ID: "metric1", Value: 10.5},
					},
					want: &MetricInMemRepo[float64]{storage: StorageState[float64]{
						"metric1": {ID: "metric1", Value: 10.5},
					}},
				},
			}

			for _, tt := range tests {
				t.Run(
					tt.name, func(t *testing.T) {
						t.Parallel()
						repo := NewMetricInMemRepo[float64](tt.state)
						assert.Equal(t, *tt.want, *repo)
					},
				)
			}
		},
	)
}

func TestMetricInMemRepo_List(t *testing.T) {
	t.Parallel()

	type testCase[T interface{ int64 | float64 }] struct {
		name string
		r    MetricInMemRepo[T]
		want []model.Metrics[T]
	}

	t.Run(
		"int64", func(t *testing.T) {
			t.Parallel()

			tests := []testCase[int64]{
				{
					name: "List int64 metrics",
					r: MetricInMemRepo[int64]{
						storage: StorageState[int64]{
							"metric1": {ID: "metric1", Value: 10},
							"metric2": {ID: "metric2", Value: 20},
						},
					},
					want: []model.Metrics[int64]{
						{ID: "metric1", Value: 10},
						{ID: "metric2", Value: 20},
					},
				},
				{
					name: "List from empty int64 repo",
					r:    MetricInMemRepo[int64]{storage: make(StorageState[int64])},
					want: []model.Metrics[int64]{},
				},
			}

			for _, tt := range tests {
				t.Run(
					tt.name, func(t *testing.T) {
						t.Parallel()
						metrics := tt.r.List()
						assert.ElementsMatch(t, tt.want, metrics)
					},
				)
			}
		},
	)

	t.Run(
		"float64", func(t *testing.T) {
			t.Parallel()

			tests := []testCase[float64]{
				{
					name: "List float64 metrics",
					r: MetricInMemRepo[float64]{
						storage: StorageState[float64]{
							"metric1": {ID: "metric1", Value: 10.5},
							"metric2": {ID: "metric2", Value: 20.3},
						},
					},
					want: []model.Metrics[float64]{
						{ID: "metric1", Value: 10.5},
						{ID: "metric2", Value: 20.3},
					},
				},
				{
					name: "List from empty float64 repo",
					r:    MetricInMemRepo[float64]{storage: make(StorageState[float64])},
					want: []model.Metrics[float64]{},
				},
			}

			for _, tt := range tests {
				t.Run(
					tt.name, func(t *testing.T) {
						t.Parallel()
						metrics := tt.r.List()
						assert.ElementsMatch(t, tt.want, metrics)
					},
				)
			}
		},
	)
}

func TestMetricInMemRepo_GetMetrics(t *testing.T) {
	t.Parallel()

	t.Run("int64", func(t *testing.T) {
		t.Parallel()

		repo := MetricInMemRepo[int64]{
			storage: StorageState[int64]{
				"metric1": {ID: "metric1", Type: model.Counter, Value: 10},
				"metric2": {ID: "metric2", Type: model.Counter, Value: 20},
			},
		}

		want := []*model.MetricsDto{
			{ID: "metric1", Type: model.Counter, Delta: pkg.Ptr[int64](10)},
			{ID: "metric2", Type: model.Counter, Delta: pkg.Ptr[int64](20)},
		}

		got := repo.GetMetrics()
		assert.ElementsMatch(t, want, got)
	})

	t.Run("float64", func(t *testing.T) {
		t.Parallel()

		repo := MetricInMemRepo[float64]{
			storage: StorageState[float64]{
				"metric1": {ID: "metric1", Type: model.Gauge, Value: 10.5},
				"metric2": {ID: "metric2", Type: model.Gauge, Value: 20.3},
			},
		}

		want := []*model.MetricsDto{
			{ID: "metric1", Type: model.Gauge, Value: pkg.Ptr(10.5)},
			{ID: "metric2", Type: model.Gauge, Value: pkg.Ptr(20.3)},
		}

		got := repo.GetMetrics()
		assert.ElementsMatch(t, want, got)
	})
}
