package agent

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yogenyslav/ya-metrics/internal/agent/collector"
	"github.com/yogenyslav/ya-metrics/internal/model"
	"github.com/yogenyslav/ya-metrics/pkg/errs"
	"github.com/yogenyslav/ya-metrics/pkg/retry"
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

func encodeMetrics(metrics []*model.MetricsDto, compressionType string) (*bytes.Buffer, error) {
	body, err := json.Marshal(metrics)
	if err != nil {
		return nil, errs.Wrap(err, "marshal metric")
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

	return buf, nil
}

func (a *Agent) sendAllMetrics(ctx context.Context, coll *collector.Collector) error {
	gaugeMetrics := coll.GetAllGaugeMetrics()
	counterMetrics := coll.GetAllCounterMetrics()

	metrics := make([]*model.MetricsDto, 0, len(gaugeMetrics)+len(counterMetrics))

	for _, metric := range gaugeMetrics {
		metrics = append(metrics, metric.ToDto())
	}

	for _, metric := range counterMetrics {
		metrics = append(metrics, metric.ToDto())
	}

	buf, err := encodeMetrics(metrics, a.cfg.CompressionType)
	if err != nil {
		return errs.Wrap(err, "encode metrics")
	}

	if err = a.sendMetricsBatch(ctx, buf); err != nil {
		return errs.Wrap(err, "send metrics batch")
	}

	return nil
}

func (a *Agent) sendMetricsBatch(ctx context.Context, reqBuff *bytes.Buffer) error {
	req, err := http.NewRequestWithContext(
		ctx, http.MethodPost, a.cfg.ServerAddr+"/updates/", reqBuff,
	)
	if err != nil {
		return errs.Wrap(err, "create request")
	}

	req.Header.Set("Content-Type", "application/json")
	if a.cfg.CompressionType != "" {
		req.Header.Set("Accept-Encoding", a.cfg.CompressionType)
		req.Header.Set("Content-Encoding", a.cfg.CompressionType)
	}

	if a.cfg.SecureKey != "" {
		req.Header.Set("Hash", a.sg.SignatureSHA256(reqBuff.Bytes()))
	}

	var buff bytes.Buffer
	err = retry.WithLinearBackoffRetry(ctx, a.cfg.Retry, func(context.Context) error {
		resp, err := a.client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode >= http.StatusInternalServerError {
			return fmt.Errorf("got status code: %d", resp.StatusCode)
		} else if resp.StatusCode >= http.StatusBadRequest {
			return errs.Wrap(retry.ErrUnretriable, fmt.Sprintf("got status code: %d", resp.StatusCode))
		}

		_, err = io.Copy(&buff, resp.Body)
		if err != nil {
			return errs.Wrap(retry.ErrUnretriable, err.Error())
		}

		return nil
	})
	if err != nil {
		if errors.Is(err, retry.ErrUnretriable) {
			return errs.Wrap(ErrUpdateMetric)
		}
		return errs.Wrap(err, "send request")
	}

	return nil
}
