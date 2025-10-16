package agent

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yogenyslav/ya-metrics/internal/agent/collector"
	"github.com/yogenyslav/ya-metrics/internal/agent/config"
	models "github.com/yogenyslav/ya-metrics/internal/model"
	"github.com/yogenyslav/ya-metrics/pkg/errs"
)

// ErrUpdateMetric indicates a failure to update a metric.
var ErrUpdateMetric = errors.New("failed to update metric")

// Client is an interface that defines the Do method for making HTTP requests.
type Client interface {
	Do(r *http.Request) (*http.Response, error)
}

// Agent struct to collect and send metrics to server.
type Agent struct {
	cfg    *config.Config
	client Client
}

// New creates a new Agent instance.
func New(cfg *config.Config, client Client) *Agent {
	return &Agent{
		cfg:    cfg,
		client: client,
	}
}

// Start begins the metric collection and reporting process.
func (a *Agent) Start(ctx context.Context) error {
	coll := collector.NewCollector(a.cfg.Agent.PollIntervalSec)

	go func() {
		ticker := time.NewTicker(time.Second * time.Duration(a.cfg.Agent.ReportIntervalSec))
		defer ticker.Stop()

		coll.Collect(ctx)

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				err := a.sendAllMetrics(ctx, coll)
				if err != nil {
					log.Printf("failed to send metrics: %v\n", err)
				}
			}
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	return nil
}

func (a *Agent) sendAllMetrics(ctx context.Context, coll *collector.Collector) error {
	err := make([]error, 27)
	err = append(err, sendMetric(ctx, coll.MemoryMetrics.Alloc, a.cfg.ServerURL, a.client))
	err = append(err, sendMetric(ctx, coll.MemoryMetrics.BuckHashSys, a.cfg.ServerURL, a.client))
	err = append(err, sendMetric(ctx, coll.MemoryMetrics.Frees, a.cfg.ServerURL, a.client))
	err = append(err, sendMetric(ctx, coll.MemoryMetrics.GCCPUFraction, a.cfg.ServerURL, a.client))
	err = append(err, sendMetric(ctx, coll.MemoryMetrics.GCSys, a.cfg.ServerURL, a.client))
	err = append(err, sendMetric(ctx, coll.MemoryMetrics.HeapAlloc, a.cfg.ServerURL, a.client))
	err = append(err, sendMetric(ctx, coll.MemoryMetrics.HeapIdle, a.cfg.ServerURL, a.client))
	err = append(err, sendMetric(ctx, coll.MemoryMetrics.HeapInuse, a.cfg.ServerURL, a.client))
	err = append(err, sendMetric(ctx, coll.MemoryMetrics.HeapObjects, a.cfg.ServerURL, a.client))
	err = append(err, sendMetric(ctx, coll.MemoryMetrics.HeapReleased, a.cfg.ServerURL, a.client))
	err = append(err, sendMetric(ctx, coll.MemoryMetrics.HeapSys, a.cfg.ServerURL, a.client))
	err = append(err, sendMetric(ctx, coll.MemoryMetrics.LastGC, a.cfg.ServerURL, a.client))
	err = append(err, sendMetric(ctx, coll.MemoryMetrics.Lookups, a.cfg.ServerURL, a.client))
	err = append(err, sendMetric(ctx, coll.MemoryMetrics.MCacheInuse, a.cfg.ServerURL, a.client))
	err = append(err, sendMetric(ctx, coll.MemoryMetrics.MCacheSys, a.cfg.ServerURL, a.client))
	err = append(err, sendMetric(ctx, coll.MemoryMetrics.MSpanInuse, a.cfg.ServerURL, a.client))
	err = append(err, sendMetric(ctx, coll.MemoryMetrics.MSpanSys, a.cfg.ServerURL, a.client))
	err = append(err, sendMetric(ctx, coll.MemoryMetrics.Mallocs, a.cfg.ServerURL, a.client))
	err = append(err, sendMetric(ctx, coll.MemoryMetrics.NextGC, a.cfg.ServerURL, a.client))
	err = append(err, sendMetric(ctx, coll.MemoryMetrics.NumForcedGC, a.cfg.ServerURL, a.client))
	err = append(err, sendMetric(ctx, coll.MemoryMetrics.NumGC, a.cfg.ServerURL, a.client))
	err = append(err, sendMetric(ctx, coll.MemoryMetrics.OtherSys, a.cfg.ServerURL, a.client))
	err = append(err, sendMetric(ctx, coll.MemoryMetrics.PauseTotalNs, a.cfg.ServerURL, a.client))
	err = append(err, sendMetric(ctx, coll.MemoryMetrics.StackInuse, a.cfg.ServerURL, a.client))
	err = append(err, sendMetric(ctx, coll.MemoryMetrics.StackSys, a.cfg.ServerURL, a.client))
	err = append(err, sendMetric(ctx, coll.MemoryMetrics.Sys, a.cfg.ServerURL, a.client))
	err = append(err, sendMetric(ctx, coll.MemoryMetrics.TotalAlloc, a.cfg.ServerURL, a.client))

	return errors.Join(err...)
}

func sendMetric[T int64 | float64](ctx context.Context, metric *models.Metrics[T], host string, client Client) error {
	req, err := http.NewRequestWithContext(
		ctx, http.MethodPost, fmt.Sprintf(host+"/update/%s/%s/%v", metric.Type, metric.Name, metric.Value), nil,
	)
	if err != nil {
		return errs.Wrap(err, "create request")
	}

	req.Header.Set("Content-Type", "text/plain")

	resp, err := client.Do(req)
	if err != nil {
		return errs.Wrap(err, "send request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errs.Wrap(
			ErrUpdateMetric, fmt.Sprintf(
				"metric '%s' of type '%s' with value '%v' not updated, status code: %d", metric.Name, metric.Type,
				metric.Value, resp.StatusCode,
			),
		)
	}

	return nil
}
