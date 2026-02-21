package server

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/ya-metrics/internal/server/middleware"
	"github.com/yogenyslav/ya-metrics/internal/server/repository"
	"github.com/yogenyslav/ya-metrics/pkg/errs"
)

type ticker struct {
	*time.Ticker
}

// C implements Ticker for dumper.
func (t *ticker) C() <-chan time.Time {
	return t.Ticker.C
}

// Ticker is a wrapper interface for time.Ticker.
type Ticker interface {
	C() <-chan time.Time
	Stop()
}

// TickerFactory is a function type for creating new Ticker instances.
type TickerFactory func(d time.Duration) Ticker

var defaultTickerFactory = func(d time.Duration) Ticker {
	return &ticker{time.NewTicker(d)}
}

// Start the dumping process.
func (s *Server) Dumping(
	ctx context.Context,
	dumper middleware.Dumper,
	newTicker TickerFactory,
	gaugeRepo, counterRepo repository.Repo,
) {
	if s.cfg.Dump.StoreInterval <= 0 {
		return
	}

	ticker := newTicker(time.Second * time.Duration(s.cfg.Dump.StoreInterval))

	go func() {
		defer ticker.Stop()

		var err error
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C():
				err = dumper.Dump(context.Background(), gaugeRepo, counterRepo)
				if err != nil {
					log.Err(errs.Wrap(err)).Msg("dump metrics to file")
				}
			}
		}
	}()
}
