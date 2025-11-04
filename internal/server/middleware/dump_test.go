package middleware

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yogenyslav/ya-metrics/internal/model"
	"github.com/yogenyslav/ya-metrics/pkg"
	"github.com/yogenyslav/ya-metrics/tests/mocks"
)

func TestNewDumper(t *testing.T) {
	t.Parallel()

	filePath := "dump.json"
	intervalSec := 10
	restore := true
	want := &fileDumper{
		filePath:    filePath,
		intervalSec: intervalSec,
		restore:     restore,
	}

	dumper := NewDumper("dump.json", 10, true)
	assert.Equal(t, want, dumper)
}

func Test_fileDumper_Restore(t *testing.T) {
	t.Parallel()

	metricsData := []byte(`[
		{"id":"gauge1","type":"gauge","value":12.34},
		{"id":"counter1","type":"counter","delta":56}
	]`)

	filePath := t.TempDir() + "/metrics.json"
	err := os.WriteFile(filePath, metricsData, os.ModePerm)
	require.NoError(t, err)

	dumper := NewDumper(filePath, 0, true)
	gaugeMetrics, counterMetrics, err := dumper.Restore()
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
	dumper := NewDumper(filePath, 0, false)

	gaugeRepo := new(mocks.MockGaugeRepo)
	counterRepo := new(mocks.MockCounterRepo)

	gaugeMetrics := []*model.MetricsDto{
		{ID: "gauge1", Type: model.Gauge, Value: pkg.Ptr(12.34)},
	}
	counterMetrics := []*model.MetricsDto{
		{ID: "counter1", Type: model.Counter, Delta: pkg.Ptr[int64](56)},
	}

	gaugeRepo.On("GetMetrics").Return(gaugeMetrics)
	counterRepo.On("GetMetrics").Return(counterMetrics)

	err := dumper.Dump(gaugeRepo, counterRepo)
	require.NoError(t, err)

	storage, err := os.ReadFile(filePath)
	require.NoError(t, err)

	want := `[
		{"id":"gauge1","type":"gauge","value":12.34},
		{"id":"counter1","type":"counter","delta":56}
	]`

	assert.JSONEq(t, want, string(storage))
}
