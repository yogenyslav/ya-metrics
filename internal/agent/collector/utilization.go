package collector

import (
	"runtime"
	"strconv"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/yogenyslav/ya-metrics/internal/model"
	"github.com/yogenyslav/ya-metrics/pkg/errs"
)

// UtilizationMetrics holds system utilization metrics.
type UtilizationMetrics struct {
	TotalMemory    *model.Metrics[float64]
	FreeMemory     *model.Metrics[float64]
	CPUUtilization []*model.Metrics[float64]
}

// NewUtilizationMetrics initializes and returns a new UtilizationMetrics instance.
func NewUtilizationMetrics() *UtilizationMetrics {
	return &UtilizationMetrics{
		TotalMemory: model.NewGaugeMetric("TotalMemory"),
		FreeMemory:  model.NewGaugeMetric("FreeMemory"),
		CPUUtilization: func() []*model.Metrics[float64] {
			metrics := make([]*model.Metrics[float64], 0, runtime.NumCPU())
			for i := 0; i < runtime.NumCPU(); i++ {
				metrics = append(metrics, model.NewGaugeMetric("CPUutilization"+strconv.Itoa(i)))
			}
			return metrics
		}(),
	}
}

func (c *Collector) updateUtilizationMetrics() error {
	memoryUsage, err := mem.VirtualMemory()
	if err != nil {
		return errs.Wrap(err, "update memory utilization")
	}

	c.utilizationMetrics.TotalMemory.Value = float64(memoryUsage.Total)
	c.utilizationMetrics.FreeMemory.Value = float64(memoryUsage.Free)

	cpuUsage, err := cpu.Percent(time.Duration(0), true)
	if err != nil {
		return errs.Wrap(err, "update cpu utilization")
	}
	for i, cpuPercent := range cpuUsage {
		if i >= len(c.utilizationMetrics.CPUUtilization) {
			break
		}
		c.utilizationMetrics.CPUUtilization[i].Value = cpuPercent
	}

	return nil
}
