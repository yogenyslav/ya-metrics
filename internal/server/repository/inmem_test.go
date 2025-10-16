package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yogenyslav/ya-metrics/internal/model"
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
						storage: map[string]*model.Metrics[int64]{
							"metric1": {Name: "metric1", Value: 10},
						},
					},
					args:       args{name: "metric1"},
					wantMetric: &model.Metrics[int64]{Name: "metric1", Value: 10},
					wantExist:  true,
				},
				{
					name:       "Get non-existing int64 metric",
					r:          MetricInMemRepo[int64]{storage: make(map[string]*model.Metrics[int64])},
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
						storage: map[string]*model.Metrics[float64]{
							"metric1": {Name: "metric1", Value: 10.5},
						},
					},
					args:       args{name: "metric1"},
					wantMetric: &model.Metrics[float64]{Name: "metric1", Value: 10.5},
					wantExist:  true,
				},
				{
					name:       "Get non-existing float64 metric",
					r:          MetricInMemRepo[float64]{storage: make(map[string]*model.Metrics[float64])},
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
					r:    MetricInMemRepo[int64]{storage: make(map[string]*model.Metrics[int64])},
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
					r:    MetricInMemRepo[float64]{storage: make(map[string]*model.Metrics[float64])},
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
						storage: map[string]*model.Metrics[int64]{
							"metric1": {Name: "metric1", Value: 10},
						},
					},
					args: args[int64]{name: "metric1", delta: 5},
				},
				{
					name: "Update non-existing int64 metric",
					r:    MetricInMemRepo[int64]{storage: make(map[string]*model.Metrics[int64])},
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
						storage: map[string]*model.Metrics[float64]{
							"metric1": {Name: "metric1", Value: 10.5},
						},
					},
					args: args[float64]{name: "metric1", delta: 5.5},
				},
				{
					name: "Update non-existing float64 metric",
					r:    MetricInMemRepo[float64]{storage: make(map[string]*model.Metrics[float64])},
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
		name string
		want *MetricInMemRepo[T]
	}

	t.Run(
		"int64", func(t *testing.T) {
			t.Parallel()

			tests := []testCase[int64]{
				{
					name: "Create new int64 MetricInMemRepo",
					want: &MetricInMemRepo[int64]{storage: make(map[string]*model.Metrics[int64])},
				},
			}

			for _, tt := range tests {
				t.Run(
					tt.name, func(t *testing.T) {
						t.Parallel()
						repo := NewMetricInMemRepo[int64]()
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
					name: "Create new float64 MetricInMemRepo",
					want: &MetricInMemRepo[float64]{storage: make(map[string]*model.Metrics[float64])},
				},
			}

			for _, tt := range tests {
				t.Run(
					tt.name, func(t *testing.T) {
						t.Parallel()
						repo := NewMetricInMemRepo[float64]()
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
						storage: map[string]*model.Metrics[int64]{
							"metric1": {Name: "metric1", Value: 10},
							"metric2": {Name: "metric2", Value: 20},
						},
					},
					want: []model.Metrics[int64]{
						{Name: "metric1", Value: 10},
						{Name: "metric2", Value: 20},
					},
				},
				{
					name: "List from empty int64 repo",
					r:    MetricInMemRepo[int64]{storage: make(map[string]*model.Metrics[int64])},
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
						storage: map[string]*model.Metrics[float64]{
							"metric1": {Name: "metric1", Value: 10.5},
							"metric2": {Name: "metric2", Value: 20.3},
						},
					},
					want: []model.Metrics[float64]{
						{Name: "metric1", Value: 10.5},
						{Name: "metric2", Value: 20.3},
					},
				},
				{
					name: "List from empty float64 repo",
					r:    MetricInMemRepo[float64]{storage: make(map[string]*model.Metrics[float64])},
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
