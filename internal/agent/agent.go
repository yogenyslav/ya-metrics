package agent

import (
	"bytes"
	"errors"
	"net/http"
	"sync"

	"github.com/rs/zerolog"
	"github.com/yogenyslav/ya-metrics/internal/agent/config"
)

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
	buffPool *sync.Pool
	mu       *sync.Mutex
}

// New creates a new Agent instance.
func New(client Client, cfg *config.Config, sg SignatureGenerator, l *zerolog.Logger) *Agent {
	return &Agent{
		client: client,
		cfg:    cfg,
		sg:     sg,
		l:      l,
		buffPool: &sync.Pool{
			New: func() any {
				return &bytes.Buffer{}
			},
		},
		mu: &sync.Mutex{},
	}
}
