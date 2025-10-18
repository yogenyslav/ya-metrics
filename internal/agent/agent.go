package agent

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/yogenyslav/ya-metrics/internal/agent/collector"
	"github.com/yogenyslav/ya-metrics/internal/model"
	"github.com/yogenyslav/ya-metrics/pkg/errs"
)

const (
	defaultServerAddr     = "localhost:8080"
	defaultPollInterval   = 2
	defaultReportInterval = 10
)

// ErrUpdateMetric indicates a failure to update a metric.
var ErrUpdateMetric = errors.New("failed to update metric")

// Client is an interface that defines the Do method for making HTTP requests.
type Client interface {
	Do(r *http.Request) (*http.Response, error)
}

// Agent struct to collect and send metrics to server.
type Agent struct {
	client            Client
	serverAddr        string
	pollIntervalSec   int
	reportIntervalSec int
}

// New creates a new Agent instance.
func New(client Client) (*Agent, error) {
	flags := flag.NewFlagSet("agent", flag.ExitOnError)
	serverAddrFlag := flags.String("a", defaultServerAddr, "адрес сервера в формате ip:port")
	pollIntervalFlag := flags.Int("p", defaultPollInterval, "интервал опроса метрик, сек.")
	reportIntervalFlag := flags.Int("r", defaultReportInterval, "интервал отправки метрик на сервер, сек. ")

	err := flags.Parse(os.Args[1:])
	if err != nil {
		return nil, errs.Wrap(err, "parse flags")
	}

	serverAddr := *serverAddrFlag
	if !strings.HasPrefix(serverAddr, "http://") && !strings.HasPrefix(serverAddr, "https://") {
		serverAddr = "http://" + serverAddr
	}

	return &Agent{
		client:            client,
		serverAddr:        serverAddr,
		pollIntervalSec:   *pollIntervalFlag,
		reportIntervalSec: *reportIntervalFlag,
	}, nil
}

// Start begins the metric collection and reporting process.
func (a *Agent) Start(ctx context.Context) error {
	coll := collector.NewCollector(a.pollIntervalSec)

	go func() {
		ticker := time.NewTicker(time.Second * time.Duration(a.reportIntervalSec))
		defer ticker.Stop()

		coll.Collect(ctx)

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				err := a.sendAllMetrics(ctx, coll)
				if err != nil {
					slog.Error("failed to send metrics", "error", err)
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
	err := make([]error, 0)

	err = append(err, sendMetric(ctx, coll.MemoryMetrics.Alloc, a.serverAddr, a.client))
	err = append(err, sendMetric(ctx, coll.MemoryMetrics.BuckHashSys, a.serverAddr, a.client))
	err = append(err, sendMetric(ctx, coll.MemoryMetrics.Frees, a.serverAddr, a.client))
	err = append(err, sendMetric(ctx, coll.MemoryMetrics.GCCPUFraction, a.serverAddr, a.client))
	err = append(err, sendMetric(ctx, coll.MemoryMetrics.GCSys, a.serverAddr, a.client))
	err = append(err, sendMetric(ctx, coll.MemoryMetrics.HeapAlloc, a.serverAddr, a.client))
	err = append(err, sendMetric(ctx, coll.MemoryMetrics.HeapIdle, a.serverAddr, a.client))
	err = append(err, sendMetric(ctx, coll.MemoryMetrics.HeapInuse, a.serverAddr, a.client))
	err = append(err, sendMetric(ctx, coll.MemoryMetrics.HeapObjects, a.serverAddr, a.client))
	err = append(err, sendMetric(ctx, coll.MemoryMetrics.HeapReleased, a.serverAddr, a.client))
	err = append(err, sendMetric(ctx, coll.MemoryMetrics.HeapSys, a.serverAddr, a.client))
	err = append(err, sendMetric(ctx, coll.MemoryMetrics.LastGC, a.serverAddr, a.client))
	err = append(err, sendMetric(ctx, coll.MemoryMetrics.Lookups, a.serverAddr, a.client))
	err = append(err, sendMetric(ctx, coll.MemoryMetrics.MCacheInuse, a.serverAddr, a.client))
	err = append(err, sendMetric(ctx, coll.MemoryMetrics.MCacheSys, a.serverAddr, a.client))
	err = append(err, sendMetric(ctx, coll.MemoryMetrics.MSpanInuse, a.serverAddr, a.client))
	err = append(err, sendMetric(ctx, coll.MemoryMetrics.MSpanSys, a.serverAddr, a.client))
	err = append(err, sendMetric(ctx, coll.MemoryMetrics.Mallocs, a.serverAddr, a.client))
	err = append(err, sendMetric(ctx, coll.MemoryMetrics.NextGC, a.serverAddr, a.client))
	err = append(err, sendMetric(ctx, coll.MemoryMetrics.NumForcedGC, a.serverAddr, a.client))
	err = append(err, sendMetric(ctx, coll.MemoryMetrics.NumGC, a.serverAddr, a.client))
	err = append(err, sendMetric(ctx, coll.MemoryMetrics.OtherSys, a.serverAddr, a.client))
	err = append(err, sendMetric(ctx, coll.MemoryMetrics.PauseTotalNs, a.serverAddr, a.client))
	err = append(err, sendMetric(ctx, coll.MemoryMetrics.StackInuse, a.serverAddr, a.client))
	err = append(err, sendMetric(ctx, coll.MemoryMetrics.StackSys, a.serverAddr, a.client))
	err = append(err, sendMetric(ctx, coll.MemoryMetrics.Sys, a.serverAddr, a.client))
	err = append(err, sendMetric(ctx, coll.MemoryMetrics.TotalAlloc, a.serverAddr, a.client))
	err = append(err, sendMetric(ctx, coll.PollCount, a.serverAddr, a.client))
	err = append(err, sendMetric(ctx, coll.RandomValue, a.serverAddr, a.client))

	return errors.Join(err...)
}

func sendMetric[T int64 | float64](ctx context.Context, metric *model.Metrics[T], host string, client Client) error {
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
