package repository

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/yogenyslav/ya-metrics/internal/model"
	"github.com/yogenyslav/ya-metrics/pkg"
	"github.com/yogenyslav/ya-metrics/tests/mocks"
)

func TestNewDumper(t *testing.T) {
	t.Parallel()

	filePath := "dump.json"
	intervalSec := 10
	want := &fileDumper{
		filePath:    filePath,
		intervalSec: intervalSec,
	}

	dumper := NewDumper("dump.json", 10)
	assert.Equal(t, want, dumper)
}

func TestRestoreMetrics(t *testing.T) {
	t.Parallel()

	metricsData := []byte(`[
		{"id":"gauge1","type":"gauge","value":12.34},
		{"id":"counter1","type":"counter","delta":56}
	]`)

	filePath := t.TempDir() + "/metrics.json"
	err := os.WriteFile(filePath, metricsData, os.ModePerm)
	require.NoError(t, err)

	gaugeMetrics, counterMetrics, err := RestoreMetrics(filePath)
	require.NoError(t, err)

	gaugeMetric, ok := gaugeMetrics["gauge1"]
	require.True(t, ok)
	assert.Equal(t, "gauge1", gaugeMetric.ID)
	assert.Equal(t, "gauge", gaugeMetric.Type)
	assert.Equal(t, 12.34, gaugeMetric.Value)

	counterMetric, ok := counterMetrics["counter1"]
	require.True(t, ok)
	assert.Equal(t, "counter1", counterMetric.ID)
	assert.Equal(t, "counter", counterMetric.Type)
	assert.Equal(t, int64(56), counterMetric.Value)
}

func Test_fileDumper_Dump(t *testing.T) {
	t.Parallel()

	filePath := t.TempDir() + "/metrics.json"
	dumper := NewDumper(filePath, 0)

	gaugeRepo := new(mocks.MockGaugeRepo)
	counterRepo := new(mocks.MockCounterRepo)

	gaugeMetrics := []*model.MetricsDto{
		{ID: "gauge1", Type: model.Gauge, Value: pkg.Ptr(12.34)},
	}
	counterMetrics := []*model.MetricsDto{
		{ID: "counter1", Type: model.Counter, Delta: pkg.Ptr[int64](56)},
	}

	gaugeRepo.On("GetMetrics", mock.Anything).Return(gaugeMetrics, nil)
	counterRepo.On("GetMetrics", mock.Anything).Return(counterMetrics, nil)

	err := dumper.Dump(context.Background(), gaugeRepo, counterRepo)
	require.NoError(t, err)

	storage, err := os.ReadFile(filePath)
	require.NoError(t, err)

	want := `[
		{"id":"gauge1","type":"gauge","value":12.34},
		{"id":"counter1","type":"counter","delta":56}
	]`

	assert.JSONEq(t, want, string(storage))
}

type testTicker struct {
	ch chan time.Time
}

// C implements Ticker.
func (t *testTicker) C() <-chan time.Time {
	return t.ch
}

// Stop implements Ticker.
func (t *testTicker) Stop() {}

func Test_fileDumper_Start(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ticker := &testTicker{ch: make(chan time.Time)}
	getTicker := func(time.Duration) Ticker { return ticker }
	called := make(chan struct{}, 1)
	onTick := func(ctx context.Context, g Repo, c Repo) error {
		called <- struct{}{}
		return nil
	}

	d := &fileDumper{intervalSec: 1}
	d.Start(ctx, nil, nil, getTicker, onTick)

	ticker.ch <- time.Now()

	select {
	case <-called:
		// ok
	case <-time.After(time.Second):
		t.Fail()
	}
}
