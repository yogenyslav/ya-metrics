package audit

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/ya-metrics/internal/config"
	"github.com/yogenyslav/ya-metrics/pkg/errs"
	"golang.org/x/sync/errgroup"
)

type source interface {
	Log(ctx context.Context, data []byte) error
}

type fileSource struct {
	filePath string
	mu       *sync.Mutex
}

func (fs *fileSource) Log(ctx context.Context, data []byte) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	file, err := os.OpenFile(fs.filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return errs.Wrap(err, "open audit log file")
	}
	defer file.Close()

	_, err = file.Write(append(data, '\n'))
	if err != nil {
		return errs.Wrap(err, "write audit log entry to file")
	}
	return nil
}

type serviceSource struct {
	url string
}

func (ss *serviceSource) Log(ctx context.Context, data []byte) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, ss.url, bytes.NewReader(data))
	if err != nil {
		return errs.Wrap(err, "create audit log request")
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errs.Wrap(err, "send audit log request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Ctx(ctx).Error().Int("status_code", resp.StatusCode).Msg("audit log service status")
	}

	return nil
}

// Entry is the structure for audit log entries.
type Entry struct {
	Ts      int64    `json:"ts"`
	Metrics []string `json:"metrics"`
	IPAddr  string   `json:"ip_address"`
}

// Audit handles audit logging.
type Audit struct {
	cfg     *config.AuditConfig
	sources []source
}

// New creates a new Audit instance.
func New(cfg *config.AuditConfig) *Audit {
	var sources []source
	if cfg.File != "" {
		sources = append(sources, &fileSource{
			filePath: cfg.File,
			mu:       &sync.Mutex{},
		})
	}
	if cfg.URL != "" {
		sources = append(sources, &serviceSource{
			url: cfg.URL,
		})
	}
	return &Audit{
		cfg:     cfg,
		sources: sources,
	}
}

// LogMetrics logs the given metrics as audit entry.
func (a *Audit) LogMetrics(ctx context.Context, metrics []string, ipAddr string) error {
	auditEntry := &Entry{
		Ts:      time.Now().Unix(),
		Metrics: metrics,
		IPAddr:  ipAddr,
	}
	data, err := json.Marshal(auditEntry)
	if err != nil {
		return errs.Wrap(err, "marshal audit entry")
	}

	g, ctx := errgroup.WithContext(ctx)
	for _, src := range a.sources {
		g.Go(func() error {
			return src.Log(ctx, data)
		})
	}

	if err := g.Wait(); err != nil {
		return errs.Wrap(err, "log audit entry to sources")
	}

	return nil
}
