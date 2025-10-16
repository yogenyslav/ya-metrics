package collector

import (
	"runtime"

	"github.com/yogenyslav/ya-metrics/internal/model"
)

// MemoryMetrics holds various memory statistics.
type MemoryMetrics struct {
	Alloc         *model.Metrics[float64]
	BuckHashSys   *model.Metrics[float64]
	Frees         *model.Metrics[float64]
	GCCPUFraction *model.Metrics[float64]
	GCSys         *model.Metrics[float64]
	HeapAlloc     *model.Metrics[float64]
	HeapIdle      *model.Metrics[float64]
	HeapInuse     *model.Metrics[float64]
	HeapObjects   *model.Metrics[float64]
	HeapReleased  *model.Metrics[float64]
	HeapSys       *model.Metrics[float64]
	LastGC        *model.Metrics[float64]
	Lookups       *model.Metrics[float64]
	MCacheInuse   *model.Metrics[float64]
	MCacheSys     *model.Metrics[float64]
	MSpanInuse    *model.Metrics[float64]
	MSpanSys      *model.Metrics[float64]
	Mallocs       *model.Metrics[float64]
	NextGC        *model.Metrics[float64]
	NumForcedGC   *model.Metrics[float64]
	NumGC         *model.Metrics[float64]
	OtherSys      *model.Metrics[float64]
	PauseTotalNs  *model.Metrics[float64]
	StackInuse    *model.Metrics[float64]
	StackSys      *model.Metrics[float64]
	Sys           *model.Metrics[float64]
	TotalAlloc    *model.Metrics[float64]
}

// NewMemoryMetrics initializes and returns a new MemoryMetrics instance.
func NewMemoryMetrics() *MemoryMetrics {
	return &MemoryMetrics{
		Alloc:         model.NewGaugeMetric("alloc"),
		BuckHashSys:   model.NewGaugeMetric("buck_hash_sys"),
		Frees:         model.NewGaugeMetric("frees"),
		GCCPUFraction: model.NewGaugeMetric("gc_cpu_fraction"),
		GCSys:         model.NewGaugeMetric("gc_sys"),
		HeapAlloc:     model.NewGaugeMetric("heap_alloc"),
		HeapIdle:      model.NewGaugeMetric("heap_idle"),
		HeapInuse:     model.NewGaugeMetric("heap_inuse"),
		HeapObjects:   model.NewGaugeMetric("heap_objects"),
		HeapReleased:  model.NewGaugeMetric("heap_released"),
		HeapSys:       model.NewGaugeMetric("heap_sys"),
		LastGC:        model.NewGaugeMetric("last_gc"),
		Lookups:       model.NewGaugeMetric("lookups"),
		MCacheInuse:   model.NewGaugeMetric("mcache_inuse"),
		MCacheSys:     model.NewGaugeMetric("mcache_sys"),
		MSpanInuse:    model.NewGaugeMetric("mspan_inuse"),
		MSpanSys:      model.NewGaugeMetric("mspan_sys"),
		Mallocs:       model.NewGaugeMetric("mallocs"),
		NextGC:        model.NewGaugeMetric("next_gc"),
		NumForcedGC:   model.NewGaugeMetric("num_forced_gc"),
		NumGC:         model.NewGaugeMetric("num_gc"),
		OtherSys:      model.NewGaugeMetric("other_sys"),
		PauseTotalNs:  model.NewGaugeMetric("pause_total_ns"),
		StackInuse:    model.NewGaugeMetric("stack_inuse"),
		StackSys:      model.NewGaugeMetric("stack_sys"),
		Sys:           model.NewGaugeMetric("sys"),
		TotalAlloc:    model.NewGaugeMetric("total_alloc"),
	}
}

func (c *Collector) updateMemoryMetrics() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	c.MemoryMetrics.Alloc.Value = float64(memStats.Alloc)
	c.MemoryMetrics.BuckHashSys.Value = float64(memStats.BuckHashSys)
	c.MemoryMetrics.Frees.Value = float64(memStats.Frees)
	c.MemoryMetrics.GCCPUFraction.Value = memStats.GCCPUFraction
	c.MemoryMetrics.GCSys.Value = float64(memStats.GCSys)
	c.MemoryMetrics.HeapAlloc.Value = float64(memStats.HeapAlloc)
	c.MemoryMetrics.HeapIdle.Value = float64(memStats.HeapIdle)
	c.MemoryMetrics.HeapInuse.Value = float64(memStats.HeapInuse)
	c.MemoryMetrics.HeapObjects.Value = float64(memStats.HeapObjects)
	c.MemoryMetrics.HeapReleased.Value = float64(memStats.HeapReleased)
	c.MemoryMetrics.HeapSys.Value = float64(memStats.HeapSys)
	c.MemoryMetrics.LastGC.Value = float64(memStats.LastGC)
	c.MemoryMetrics.Lookups.Value = float64(memStats.Lookups)
	c.MemoryMetrics.MCacheInuse.Value = float64(memStats.MCacheInuse)
	c.MemoryMetrics.MCacheSys.Value = float64(memStats.MCacheSys)
	c.MemoryMetrics.MSpanInuse.Value = float64(memStats.MSpanInuse)
	c.MemoryMetrics.MSpanSys.Value = float64(memStats.MSpanSys)
	c.MemoryMetrics.Mallocs.Value = float64(memStats.Mallocs)
	c.MemoryMetrics.NextGC.Value = float64(memStats.NextGC)
	c.MemoryMetrics.NumForcedGC.Value = float64(memStats.NumForcedGC)
	c.MemoryMetrics.NumGC.Value = float64(memStats.NumGC)
	c.MemoryMetrics.OtherSys.Value = float64(memStats.OtherSys)
	c.MemoryMetrics.PauseTotalNs.Value = float64(memStats.PauseTotalNs)
	c.MemoryMetrics.StackInuse.Value = float64(memStats.StackInuse)
	c.MemoryMetrics.StackSys.Value = float64(memStats.StackSys)
	c.MemoryMetrics.Sys.Value = float64(memStats.Sys)
	c.MemoryMetrics.TotalAlloc.Value = float64(memStats.TotalAlloc)
}
