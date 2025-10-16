package collector

import (
	"context"
	"time"

	"github.com/yogenyslav/ya-metrics/internal/model"
)

// Collector struct to collect metrics.
type Collector struct {
	MemoryMetrics *MemoryMetrics
	PollCount     *model.Metrics[int64]
	RandomValue   *model.Metrics[float64]
	PollInterval  int
}

// NewCollector creates a new Collector instance.
func NewCollector(pollInterval int) *Collector {
	c := &Collector{
		MemoryMetrics: NewMemoryMetrics(),
		PollCount:     model.NewCounterMetric("poll_count"),
		RandomValue:   model.NewGaugeMetric("random_value"),
		PollInterval:  pollInterval,
	}

	return c
}

// Collect starts collecting metrics at specified intervals.
func (c *Collector) Collect(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(time.Second * time.Duration(c.PollInterval))
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				c.updateMetrics()
			}
		}
	}()
}

func (c *Collector) updateMetrics() {
	c.PollCount.Value++
	c.RandomValue.Value = float64(time.Now().UnixNano()%100) + 1
	c.updateMemoryMetrics()
}
