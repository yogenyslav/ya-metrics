package collector

import (
	"runtime"

	models "github.com/yogenyslav/ya-metrics/internal/model"
)

// MemoryMetrics holds various memory statistics.
type MemoryMetrics struct {
	Alloc         *models.Metrics[float64]
	BuckHashSys   *models.Metrics[float64]
	Frees         *models.Metrics[float64]
	GCCPUFraction *models.Metrics[float64]
	GCSys         *models.Metrics[float64]
	HeapAlloc     *models.Metrics[float64]
	HeapIdle      *models.Metrics[float64]
	HeapInuse     *models.Metrics[float64]
	HeapObjects   *models.Metrics[float64]
	HeapReleased  *models.Metrics[float64]
	HeapSys       *models.Metrics[float64]
	LastGC        *models.Metrics[float64]
	Lookups       *models.Metrics[float64]
	MCacheInuse   *models.Metrics[float64]
	MCacheSys     *models.Metrics[float64]
	MSpanInuse    *models.Metrics[float64]
	MSpanSys      *models.Metrics[float64]
	Mallocs       *models.Metrics[float64]
	NextGC        *models.Metrics[float64]
	NumForcedGC   *models.Metrics[float64]
	NumGC         *models.Metrics[float64]
	OtherSys      *models.Metrics[float64]
	PauseTotalNs  *models.Metrics[float64]
	StackInuse    *models.Metrics[float64]
	StackSys      *models.Metrics[float64]
	Sys           *models.Metrics[float64]
	TotalAlloc    *models.Metrics[float64]
}

// NewMemoryMetrics initializes and returns a new MemoryMetrics instance.
func NewMemoryMetrics() *MemoryMetrics {
	return &MemoryMetrics{
		Alloc:         models.NewGaugeMetric("Alloc"),
		BuckHashSys:   models.NewGaugeMetric("BuckHashSys"),
		Frees:         models.NewGaugeMetric("Frees"),
		GCCPUFraction: models.NewGaugeMetric("GCCPUFraction"),
		GCSys:         models.NewGaugeMetric("GCSys"),
		HeapAlloc:     models.NewGaugeMetric("HeapAlloc"),
		HeapIdle:      models.NewGaugeMetric("HeapIdle"),
		HeapInuse:     models.NewGaugeMetric("HeapInuse"),
		HeapObjects:   models.NewGaugeMetric("HeapObjects"),
		HeapReleased:  models.NewGaugeMetric("HeapReleased"),
		HeapSys:       models.NewGaugeMetric("HeapSys"),
		LastGC:        models.NewGaugeMetric("LastGC"),
		Lookups:       models.NewGaugeMetric("Lookups"),
		MCacheInuse:   models.NewGaugeMetric("MCacheInuse"),
		MCacheSys:     models.NewGaugeMetric("MCacheSys"),
		MSpanInuse:    models.NewGaugeMetric("MSpanInuse"),
		MSpanSys:      models.NewGaugeMetric("MSpanSys"),
		Mallocs:       models.NewGaugeMetric("Mallocs"),
		NextGC:        models.NewGaugeMetric("NextGC"),
		NumForcedGC:   models.NewGaugeMetric("NumForcedGC"),
		NumGC:         models.NewGaugeMetric("NumGC"),
		OtherSys:      models.NewGaugeMetric("OtherSys"),
		PauseTotalNs:  models.NewGaugeMetric("PauseTotalNs"),
		StackInuse:    models.NewGaugeMetric("StackInuse"),
		StackSys:      models.NewGaugeMetric("StackSys"),
		Sys:           models.NewGaugeMetric("Sys"),
		TotalAlloc:    models.NewGaugeMetric("TotalAlloc"),
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
