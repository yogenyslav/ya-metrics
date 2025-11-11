package agent

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
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

type combinedError struct {
	err []error
}

// NewCombinedError creates an instance of combinedError.
func NewCombinedError() *combinedError {
	return &combinedError{
		err: make([]error, 0),
	}
}

// Error implements the error interface.
func (e *combinedError) Error() string {
	b := strings.Builder{}
	for _, err := range e.err {
		b.WriteString(err.Error() + "; ")
	}
	return b.String()
}

// Add an error to combinedError.
func (e *combinedError) Add(err error) {
	e.err = append(e.err, err)
}

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
	var err error
	combinedErr := NewCombinedError()

	gaugeMetrics := coll.GetAllGaugeMetrics()
	for _, metric := range gaugeMetrics {
		err = sendMetric(ctx, metric, a.cfg.ServerAddr, a.client, a.cfg.CompressionType)
		if err != nil {
			combinedErr.Add(err)
		}
	}

	counterMetric := coll.PollCount
	err = sendMetric(ctx, counterMetric, a.cfg.ServerAddr, a.client, a.cfg.CompressionType)
	if err != nil {
		combinedErr.Add(err)
	}

	return errs.Wrap(combinedErr, "send all metrics")
}

func sendMetric[T int64 | float64](
	ctx context.Context,
	metric *model.Metrics[T],
	host string,
	client Client,
	compressionType string,
) error {
	data := metric.ToDto()
	body, err := json.Marshal(data)
	if err != nil {
		return errs.Wrap(err, "marshal metric")
	}

	buf := &bytes.Buffer{}
	switch compressionType {
	case "gzip":
		gz := gzip.NewWriter(buf)
		gz.Write(body)
		gz.Close()
	default:
		buf.Write(body)
	}

	req, err := http.NewRequestWithContext(
		ctx, http.MethodPost, host+"/update/", buf,
	)
	if err != nil {
		return errs.Wrap(err, "create request")
	}

	req.Header.Set("Content-Type", "application/json")
	if compressionType != "" {
		req.Header.Set("Accept-Encoding", compressionType)
		req.Header.Set("Content-Encoding", compressionType)
	}

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
