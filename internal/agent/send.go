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
	"slices"
	"syscall"
	"time"

	"github.com/yogenyslav/ya-metrics/internal/agent/collector"
	"github.com/yogenyslav/ya-metrics/internal/model"
	"github.com/yogenyslav/ya-metrics/pkg/errs"
	"github.com/yogenyslav/ya-metrics/pkg/retry"
	"golang.org/x/sync/errgroup"
)

// Start begins the metric collection and reporting process.
func (a *Agent) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	coll := collector.NewCollector(a.cfg.PollIntervalSec, a.l)
	coll.Collect(ctx)

	go func() {
		ticker := time.NewTicker(time.Second * time.Duration(a.cfg.ReportIntervalSec))
		defer ticker.Stop()

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

func (a *Agent) encodeMetrics(metrics []*model.MetricsDto, compressionType string) ([]byte, error) {
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

	return buf.Bytes(), nil
}

func (a *Agent) sendAllMetrics(ctx context.Context, coll *collector.Collector) error {
	metrics := coll.GetAllMetrics()

	// не уверен, что я понял идею применения worker pool именно тут корректно, так как изначально у нас отправлялся один большой запрос с метриками.
	// может быть, если бы они обрабатывались ощутимое время на стороне сервера, тогда я бы лучше ощутил эту идею.
	// в любом случае, пояснение к реализации:
	// если считать batchSize через RateLimit, то будем получать batchCount <= RateLimit, что приведет к простою воркеров.
	// для примера выбран batchSize = 3, чтобы воркеры в количестве RateLimit были зафиксированы и могли ждать получения новых батчей при batchCount > RateLimit.
	// вопрос с оптимальным подбором размера батча, как мне кажется, на реальной задаче можно было бы определить только эмпирически.

	batchCount := (len(metrics) + a.cfg.BatchSize - 1) / a.cfg.BatchSize

	a.l.Info().Int("batchSize", a.cfg.BatchSize).Int("batchCount", batchCount).Msg("sending metrics")

	g, ctx := errgroup.WithContext(ctx)
	batchCh := make(chan []*model.MetricsDto, batchCount)
	for range a.cfg.RateLimit {
		g.Go(func() error {
			return a.sendMetricsBatch(ctx, batchCh)
		})
	}

	for batch := range slices.Chunk(metrics, a.cfg.BatchSize) {
		batchCh <- batch
	}
	close(batchCh)

	a.l.Debug().Msg("waiting for batch tasks to complete")
	if err := g.Wait(); err != nil {
		return errs.Wrap(errors.New("send all metrics"))
	}

	return nil
}

func (a *Agent) sendMetricsBatch(ctx context.Context, batchCh <-chan []*model.MetricsDto) error {
	for batch := range batchCh {
		req, err := a.createRequest(ctx, batch)
		if err != nil {
			return errs.Wrap(err, "create request")
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

		a.l.Info().Msg("sent metrics batch successfully")
	}

	return nil
}

func (a *Agent) createRequest(ctx context.Context, batch []*model.MetricsDto) (*http.Request, error) {
	data, err := a.encodeMetrics(batch, a.cfg.CompressionType)
	if err != nil {
		return nil, errs.Wrap(err, "encode metrics")
	}

	req, err := http.NewRequestWithContext(
		ctx, http.MethodPost, a.cfg.ServerAddr+"/updates/", bytes.NewReader(data),
	)
	if err != nil {
		return nil, errs.Wrap(err, "create request")
	}

	req.Header.Set("Content-Type", "application/json")
	if a.cfg.CompressionType != "" {
		req.Header.Set("Accept-Encoding", a.cfg.CompressionType)
		req.Header.Set("Content-Encoding", a.cfg.CompressionType)
	}

	if a.cfg.SecureKey != "" {
		req.Header.Set("HashSHA256", a.sg.SignatureSHA256(data))
	}

	return req, nil
}
