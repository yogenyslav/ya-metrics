package agent

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/ya-metrics/internal/agent/collector"
	"github.com/yogenyslav/ya-metrics/internal/agent/config"
)

const gracefulShutdownTimeout = 60 * time.Second

// ErrUpdateMetric indicates a failure to update a metric.
var ErrUpdateMetric = errors.New("failed to update metric")

// Client is an interface that defines the Do method for making HTTP requests.
type Client interface {
	Do(r *http.Request) (*http.Response, error)
}

// SignatureGenerator is an interface for generating hash signatures.
type SignatureGenerator interface {
	SignatureSHA256(data []byte) string
}

// Agent struct to collect and send metrics to server.
type Agent struct {
	client   Client
	cfg      *config.Config
	sg       SignatureGenerator
	l        *zerolog.Logger
	shutdown chan struct{}
}

// New creates a new Agent instance.
func New(client Client, cfg *config.Config, sg SignatureGenerator, l *zerolog.Logger) *Agent {
	return &Agent{
		client:   client,
		cfg:      cfg,
		sg:       sg,
		l:        l,
		shutdown: make(chan struct{}, 1),
	}
}

// Start begins the metric collection and reporting process.
func (a *Agent) Start(ctx context.Context) error {
	coll := collector.NewCollector(a.cfg.PollIntervalSec, a.l)
	coll.Collect(ctx)

	go func() {
		defer func() {
			a.shutdown <- struct{}{}
		}()

		ticker := time.NewTicker(time.Second * time.Duration(a.cfg.ReportIntervalSec))
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := a.sendAllMetrics(context.Background(), coll); err != nil {
					slog.Error("failed to send metrics", "error", err)
				}
			}
		}
	}()

	return nil
}

// Shutdown performs graceful shutdown for agent.
func (a *Agent) Shutdown() {
	select {
	case <-a.shutdown:
		log.Info().Msg("agent shutdown gracefully")
	case <-time.After(gracefulShutdownTimeout):
		log.Warn().Msg("graceful shutdown timeout reached, exiting")
	}
}
