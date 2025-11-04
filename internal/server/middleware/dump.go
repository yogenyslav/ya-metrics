package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/ya-metrics/internal/model"
	"github.com/yogenyslav/ya-metrics/internal/server/repository"
	"github.com/yogenyslav/ya-metrics/pkg/errs"
)

type inMemRepo interface {
	GetMetrics() []*model.MetricsDto
}

// fileDumper is a struct to dump data to file.
type fileDumper struct {
	filePath    string
	intervalSec int
	restore     bool
}

// NewDumper creates a new instance of fileDumper.
func NewDumper(filePath string, intervalSec int, restore bool) *fileDumper {
	return &fileDumper{
		filePath:    filePath,
		intervalSec: intervalSec,
		restore:     restore,
	}
}

// Middleware returns a dumping middleware if intervalSec is <= 0.
func (d *fileDumper) Middleware(gaugeRepo inMemRepo, counterRepo inMemRepo) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
			if d.intervalSec > 0 {
				return
			}

			err := d.Dump(gaugeRepo, counterRepo)
			if err != nil {
				log.Ctx(r.Context()).Err(errs.Wrap(err)).Msg("dump metrics to file")
			}
		})
	}
}

// Start starts the dumping process.
func (d *fileDumper) Start(ctx context.Context, gaugeRepo inMemRepo, counterRepo inMemRepo) {
	if d.intervalSec <= 0 {
		return
	}

	go func() {
		var err error

		ticker := time.NewTicker(time.Duration(d.intervalSec) * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				err = d.Dump(gaugeRepo, counterRepo)
				if err != nil {
					log.Ctx(ctx).Err(errs.Wrap(err)).Msg("dump metrics to file")
				}
			}
		}
	}()
}

// Restore the data from file.
func (d *fileDumper) Restore() (gaugeMetrics repository.StorageState[float64], counterMetrics repository.StorageState[int64], err error) {
	if !d.restore {
		return nil, nil, nil
	}

	gaugeMetrics = make(repository.StorageState[float64])
	counterMetrics = make(repository.StorageState[int64])

	f, err := os.OpenFile(d.filePath, os.O_RDONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, nil, errs.Wrap(err, "open file for restore")
	}
	defer f.Close()

	data, err := os.ReadFile(d.filePath)
	if err != nil {
		return nil, nil, errs.Wrap(err, "read data from file")
	}

	if len(data) == 0 {
		return nil, nil, nil
	}

	var v []*model.MetricsDto
	err = json.Unmarshal(data, &v)
	if err != nil {
		return nil, nil, errs.Wrap(err, "unmarshal data")
	}

	for _, m := range v {
		switch m.Type {
		case model.Gauge:
			gaugeMetrics[m.ID] = m.ToGaugeMetric()
		case model.Counter:
			counterMetrics[m.ID] = m.ToCounterMetric()
		}
	}

	return gaugeMetrics, counterMetrics, nil
}

// Dump data to file.
func (d *fileDumper) Dump(gaugeRepo inMemRepo, counterRepo inMemRepo) error {
	gaugeMetrics := gaugeRepo.GetMetrics()
	counterMetrics := counterRepo.GetMetrics()

	v := make([]*model.MetricsDto, 0, len(gaugeMetrics)+len(counterMetrics))
	v = append(v, gaugeMetrics...)
	v = append(v, counterMetrics...)

	f, err := os.OpenFile(d.filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return errs.Wrap(err, "open file for dump")
	}
	defer f.Close()

	data, err := json.Marshal(v)
	if err != nil {
		return errs.Wrap(err, "marshal data")
	}

	_, err = f.Write(data)
	if err != nil {
		return errs.Wrap(err, "write data to file")
	}

	return nil
}
