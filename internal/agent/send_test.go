package agent

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yogenyslav/ya-metrics/internal/agent/config"
	"github.com/yogenyslav/ya-metrics/internal/model"
)

func Test_encodeMetrics(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		cfg         *config.Config
		metrics     []*model.MetricsDto
		wantMetrics []byte
		wantErr     bool
	}{
		{
			name: "Encode without compression",
			cfg: &config.Config{
				CompressionType: "",
			},
			metrics: []*model.MetricsDto{
				{
					ID:    "test_counter",
					Type:  model.Counter,
					Delta: new(int64),
				},
				{
					ID:    "test_gauge",
					Type:  model.Gauge,
					Value: new(float64),
				},
			},
			wantMetrics: func() []byte {
				data, _ := json.Marshal([]*model.MetricsDto{
					{
						ID:    "test_counter",
						Type:  model.Counter,
						Delta: new(int64),
					},
					{
						ID:    "test_gauge",
						Type:  model.Gauge,
						Value: new(float64),
					},
				})
				return data
			}(),
			wantErr: false,
		},
		{
			name: "Encode with gzip compression",
			cfg: &config.Config{
				CompressionType: "gzip",
			},
			metrics: []*model.MetricsDto{
				{
					ID:    "test_counter",
					Type:  model.Counter,
					Delta: new(int64),
				},
				{
					ID:    "test_gauge",
					Type:  model.Gauge,
					Value: new(float64),
				},
			},
			wantMetrics: func() []byte {
				data, _ := json.Marshal([]*model.MetricsDto{
					{
						ID:    "test_counter",
						Type:  model.Counter,
						Delta: new(int64),
					},
					{
						ID:    "test_gauge",
						Type:  model.Gauge,
						Value: new(float64),
					},
				})

				buf := &bytes.Buffer{}
				w := gzip.NewWriter(nil)
				w.Reset(buf)
				w.Write(data)
				w.Close()

				return buf.Bytes()
			}(),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			metrics, err := encodeMetrics(tt.metrics, tt.cfg.CompressionType)
			if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, metrics)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantMetrics, metrics.Bytes())
			}
		})
	}
}
