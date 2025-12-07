package collector

import (
	"context"
	"sync"
	"time"

	"github.com/yogenyslav/ya-metrics/internal/model"
)

// Collector struct to collect metrics.
type Collector struct {
	PollInterval  int
	memoryMetrics *MemoryMetrics
	pollCount     *model.Metrics[int64]
	randomValue   *model.Metrics[float64]
	mu            *sync.Mutex
}

// NewCollector creates a new Collector instance.
func NewCollector(pollInterval int) *Collector {
	return &Collector{
		memoryMetrics: NewMemoryMetrics(),
		pollCount:     model.NewCounterMetric("PollCount"),
		randomValue:   model.NewGaugeMetric("RandomValue"),
		PollInterval:  pollInterval,
		mu:            &sync.Mutex{},
	}
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
	c.mu.Lock()
	defer c.mu.Unlock()

	c.pollCount.Value++
	c.randomValue.Value = float64(time.Now().UnixNano()%100) + 1
	c.updateMemoryMetrics()
}

// MemoryMetrics returns the current memory metrics.
func (c *Collector) MemoryMetrics() *MemoryMetrics {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.memoryMetrics
}

// PollCount returns the current poll count metric.
func (c *Collector) PollCount() *model.Metrics[int64] {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.pollCount
}

// RandomValue returns the current random value metric.
func (c *Collector) RandomValue() *model.Metrics[float64] {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.randomValue
}

// GetAllGaugeMetrics returns all gauge metrics collected by the Collector.
func (c *Collector) GetAllGaugeMetrics() []*model.Metrics[float64] {
	c.mu.Lock()
	defer c.mu.Unlock()

	return []*model.Metrics[float64]{
		c.randomValue,
		c.memoryMetrics.Alloc,
		c.memoryMetrics.BuckHashSys,
		c.memoryMetrics.Frees,
		c.memoryMetrics.GCCPUFraction,
		c.memoryMetrics.GCSys,
		c.memoryMetrics.HeapAlloc,
		c.memoryMetrics.HeapIdle,
		c.memoryMetrics.HeapInuse,
		c.memoryMetrics.HeapObjects,
		c.memoryMetrics.HeapReleased,
		c.memoryMetrics.HeapSys,
		c.memoryMetrics.LastGC,
		c.memoryMetrics.Lookups,
		c.memoryMetrics.MCacheInuse,
		c.memoryMetrics.MCacheSys,
		c.memoryMetrics.MSpanInuse,
		c.memoryMetrics.MSpanSys,
		c.memoryMetrics.Mallocs,
		c.memoryMetrics.NextGC,
		c.memoryMetrics.NumForcedGC,
		c.memoryMetrics.NumGC,
		c.memoryMetrics.OtherSys,
		c.memoryMetrics.PauseTotalNs,
		c.memoryMetrics.StackInuse,
		c.memoryMetrics.StackSys,
		c.memoryMetrics.Sys,
		c.memoryMetrics.TotalAlloc,
	}
}

// GetAllCounterMetrics returns all counter metrics collected by the Collector.
func (c *Collector) GetAllCounterMetrics() []*model.Metrics[int64] {
	c.mu.Lock()
	defer c.mu.Unlock()

	return []*model.Metrics[int64]{
		c.pollCount,
	}
}
