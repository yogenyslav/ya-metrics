package middleware

import (
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/ya-metrics/internal/server/repository"
	"github.com/yogenyslav/ya-metrics/pkg/errs"
)

type dumper interface {
	Dump(gaugeRepo repository.Repo, counterRepo repository.Repo) error
}

// WithFileDumper returns a dumping middleware if intervalSec is <= 0.
func WithFileDumper(d dumper, intervalSec int, gaugeRepo, counterRepo repository.Repo) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
			if intervalSec > 0 {
				return
			}

			err := d.Dump(gaugeRepo, counterRepo)
			if err != nil {
				log.Ctx(r.Context()).Err(errs.Wrap(err)).Msg("dump metrics to file")
			}
		})
	}
}
