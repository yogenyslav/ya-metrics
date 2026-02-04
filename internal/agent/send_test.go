package agent

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"math/rand/v2"
	"net/http"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/yogenyslav/ya-metrics/internal/agent/collector"
	"github.com/yogenyslav/ya-metrics/internal/agent/config"
	"github.com/yogenyslav/ya-metrics/internal/model"
	"github.com/yogenyslav/ya-metrics/tests/mocks"
)

func TestAgent_encodeMetrics(t *testing.T) {
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

			a := New(http.DefaultClient, tt.cfg, nil, zerolog.Ctx(context.Background()))

			metrics, err := a.encodeMetrics(tt.metrics, tt.cfg.CompressionType)
			if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, metrics)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantMetrics, metrics)
			}
		})
	}
}

func TestAgent_sendAllMetrics(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		rateLimit int
		batchSize int
		wantErr   bool
	}{
		{
			name:      "rateLimit = 1, batchSize = 3",
			rateLimit: 1,
			batchSize: 3,
			wantErr:   false,
		},
		{
			name:      "rateLimit = 1, batchSize = 1",
			rateLimit: 1,
			batchSize: 1,
			wantErr:   false,
		},
		{
			name:      "rateLimit = 5, batchSize = 3",
			rateLimit: 5,
			batchSize: 3,
			wantErr:   false,
		},
		{
			name:      "rateLimit = 5, batchSize = 1",
			rateLimit: 5,
			batchSize: 1,
			wantErr:   false,
		},
		{
			name:      "rateLimit = 10, batchSize = 3, has error",
			rateLimit: 1,
			batchSize: 3,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			pollInterval := 1
			reportInterval := 1
			c := collector.NewCollector(pollInterval, zerolog.Ctx(context.Background()))

			metricsNum := len(c.GetAllMetrics())
			batchCount := (metricsNum + tt.batchSize - 1) / tt.batchSize

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			c.Collect(ctx)

			client := new(mocks.HTTPClient)
			a := New(client, &config.Config{
				PollIntervalSec:   pollInterval,
				ReportIntervalSec: reportInterval,
				RateLimit:         tt.rateLimit,
				BatchSize:         tt.batchSize,
			}, nil, zerolog.Ctx(ctx))

			successCalls := max(rand.IntN(batchCount), 1)
			if !tt.wantErr {
				client.On("Do", mock.Anything).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       http.NoBody,
				}, nil).Times(batchCount)
			} else {
				client.On("Do", mock.Anything).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       http.NoBody,
				}, nil).Times(successCalls)

				client.On("Do", mock.Anything).Return(&http.Response{
					StatusCode: http.StatusBadRequest,
					Body:       http.NoBody,
				}, nil).Times(batchCount - successCalls)
			}

			err := a.sendAllMetrics(ctx, c)
			if tt.wantErr {
				require.Error(t, err)
				client.AssertNumberOfCalls(t, "Do", successCalls+1)
			} else {
				require.NoError(t, err)
				client.AssertNumberOfCalls(t, "Do", batchCount)
			}
		})
	}
}
