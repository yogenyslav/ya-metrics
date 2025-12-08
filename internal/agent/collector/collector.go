package collector

import (
	"context"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/yogenyslav/ya-metrics/internal/model"
)

// Collector struct to collect metrics.
type Collector struct {
	PollInterval       int
	memoryMetrics      *MemoryMetrics
	generalMetrics     *GeneralMetrics
	utilizationMetrics *UtilizationMetrics
	l                  *zerolog.Logger
	updaters           []func() error
	mu                 *sync.Mutex
	wg                 *sync.WaitGroup
}

// NewCollector creates a new Collector instance.
func NewCollector(pollInterval int, l *zerolog.Logger) *Collector {
	c := &Collector{
		PollInterval:       pollInterval,
		memoryMetrics:      NewMemoryMetrics(),
		generalMetrics:     NewGeneralMetrics(),
		utilizationMetrics: NewUtilizationMetrics(),
		l:                  l,
		mu:                 &sync.Mutex{},
		wg:                 &sync.WaitGroup{},
	}

	c.updaters = []func() error{
		c.updateMemoryMetrics,
		c.updateGeneralMetrics,
		c.updateUtilizationMetrics,
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
	c.mu.Lock()
	defer c.mu.Unlock()

	errCh := make(chan error, len(c.updaters))

	c.wg.Add(len(c.updaters))
	for _, updater := range c.updaters {
		go func(u func() error) {
			defer c.wg.Done()
			err := u()
			errCh <- err
			if err == nil {
				c.l.Info().Msg("updated metrics")
			}
		}(updater)
	}
	c.wg.Wait()

	success := true
	for range c.updaters {
		if err := <-errCh; err != nil {
			success = false
			c.l.Error().Err(err).Msg("failed to update metrics")
		}
	}

	if success {
		c.l.Info().Msg("updated all metrics")
	}
}

// MemoryMetrics returns the current memory metrics.
func (c *Collector) MemoryMetrics() *MemoryMetrics {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.memoryMetrics
}

// GeneralMetrics returns the current general metrics.
func (c *Collector) GeneralMetrics() *GeneralMetrics {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.generalMetrics
}

// GetAllGaugeMetrics returns all gauge metrics collected by the Collector.
func (c *Collector) GetAllGaugeMetrics() []*model.Metrics[float64] {
	m := []*model.Metrics[float64]{
		c.generalMetrics.RandomValue,
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
		c.utilizationMetrics.TotalMemory,
		c.utilizationMetrics.FreeMemory,
	}
	m = append(m, c.utilizationMetrics.CPUUtilization...)
	return m
}

// GetAllCounterMetrics returns all counter metrics collected by the Collector.
func (c *Collector) GetAllCounterMetrics() []*model.Metrics[int64] {
	return []*model.Metrics[int64]{
		c.generalMetrics.PollCount,
	}
}

// GetAllMetrics returns all metrics collected by the Collector.
func (c *Collector) GetAllMetrics() []*model.MetricsDto {
	c.mu.Lock()
	defer c.mu.Unlock()

	gaugeMetrics := c.GetAllGaugeMetrics()
	counterMetrics := c.GetAllCounterMetrics()

	metrics := make([]*model.MetricsDto, 0, len(gaugeMetrics)+len(counterMetrics))

	for _, metric := range gaugeMetrics {
		metrics = append(metrics, metric.ToDto())
	}

	for _, metric := range counterMetrics {
		metrics = append(metrics, metric.ToDto())
	}

	return metrics
}
