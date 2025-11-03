package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yogenyslav/ya-metrics/internal/agent/collector"
	"github.com/yogenyslav/ya-metrics/internal/model"
	"github.com/yogenyslav/ya-metrics/pkg/errs"
)

// Start begins the metric collection and reporting process.
func (a *Agent) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	coll := collector.NewCollector(a.cfg.PollIntervalSec)

	go func() {
		ticker := time.NewTicker(time.Second * time.Duration(a.cfg.ReportIntervalSec))
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

	gaugeMetrics := coll.GetAllGaugeMetrics()
	for _, metric := range gaugeMetrics {
		sendErr := sendMetric(ctx, metric, a.cfg.ServerAddr, a.client)
		if sendErr != nil {
			err = append(err, sendErr)
		}
	}

	counterMetric := coll.PollCount
	sendErr := sendMetric(ctx, counterMetric, a.cfg.ServerAddr, a.client)
	if sendErr != nil {
		err = append(err, sendErr)
	}

	return errors.Join(err...)
}

func sendMetric[T int64 | float64](ctx context.Context, metric *model.Metrics[T], host string, client Client) error {
	data := metric.ToDto()
	body, err := json.Marshal(data)
	if err != nil {
		return errs.Wrap(err, "marshal metric")
	}

	req, err := http.NewRequestWithContext(
		ctx, http.MethodPost, host+"/update/", bytes.NewReader(body),
	)
	if err != nil {
		return errs.Wrap(err, "create request")
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return errs.Wrap(err, "send request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errs.Wrap(
			ErrUpdateMetric, fmt.Sprintf(
				"metric '%s' of type '%s' with value '%v' not updated, status code: %d", metric.ID, metric.Type,
				metric.Value, resp.StatusCode,
			),
		)
	}

	return nil
}
