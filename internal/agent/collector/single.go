package collector

import (
	"time"

	"github.com/yogenyslav/ya-metrics/internal/model"
)

// GeneralMetrics holds general application metrics.
//
// generate:reset
type GeneralMetrics struct {
	PollCount   *model.Metrics[int64]
	RandomValue *model.Metrics[float64]
}

// NewGeneralMetrics initializes and returns a new GeneralMetrics instance.
func NewGeneralMetrics() *GeneralMetrics {
	return &GeneralMetrics{
		PollCount:   model.NewCounterMetric("PollCount"),
		RandomValue: model.NewGaugeMetric("RandomValue"),
	}
}

func (c *Collector) updateGeneralMetrics() error {
	c.generalMetrics.PollCount.Value++
	c.generalMetrics.RandomValue.Value = float64(time.Now().UnixNano()%100) + 1

	return nil
}
