package agent

import (
	"errors"
	"net/http"

	"github.com/yogenyslav/ya-metrics/internal/agent/config"
)

// ErrUpdateMetric indicates a failure to update a metric.
var ErrUpdateMetric = errors.New("failed to update metric")

// Client is an interface that defines the Do method for making HTTP requests.
type Client interface {
	Do(r *http.Request) (*http.Response, error)
}

// Agent struct to collect and send metrics to server.
type Agent struct {
	client Client
	cfg    *config.Config
}

// New creates a new Agent instance.
func New(client Client, cfg *config.Config) *Agent {
	return &Agent{
		client: client,
		cfg:    cfg,
	}
}
