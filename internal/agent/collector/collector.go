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
		PollCount:     model.NewCounterMetric("PollCount"),
		RandomValue:   model.NewGaugeMetric("RandomValue"),
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

// GetAllGaugeMetrics returns all gauge metrics collected by the Collector.
func (c *Collector) GetAllGaugeMetrics() []*model.Metrics[float64] {
	return []*model.Metrics[float64]{
		c.RandomValue,
		c.MemoryMetrics.Alloc,
		c.MemoryMetrics.BuckHashSys,
		c.MemoryMetrics.Frees,
		c.MemoryMetrics.GCCPUFraction,
		c.MemoryMetrics.GCSys,
		c.MemoryMetrics.HeapAlloc,
		c.MemoryMetrics.HeapIdle,
		c.MemoryMetrics.HeapInuse,
		c.MemoryMetrics.HeapObjects,
		c.MemoryMetrics.HeapReleased,
		c.MemoryMetrics.HeapSys,
		c.MemoryMetrics.LastGC,
		c.MemoryMetrics.Lookups,
		c.MemoryMetrics.MCacheInuse,
		c.MemoryMetrics.MCacheSys,
		c.MemoryMetrics.MSpanInuse,
		c.MemoryMetrics.MSpanSys,
		c.MemoryMetrics.Mallocs,
		c.MemoryMetrics.NextGC,
		c.MemoryMetrics.NumForcedGC,
		c.MemoryMetrics.NumGC,
		c.MemoryMetrics.OtherSys,
		c.MemoryMetrics.PauseTotalNs,
		c.MemoryMetrics.StackInuse,
		c.MemoryMetrics.StackSys,
		c.MemoryMetrics.Sys,
		c.MemoryMetrics.TotalAlloc,
	}
}

// GetAllCounterMetrics returns all counter metrics collected by the Collector.
func (c *Collector) GetAllCounterMetrics() []*model.Metrics[int64] {
	return []*model.Metrics[int64]{
		c.PollCount,
	}
}
