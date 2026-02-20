package repository

import (
	"context"
	"encoding/json"
	"os"

	"github.com/yogenyslav/ya-metrics/internal/model"
	"github.com/yogenyslav/ya-metrics/pkg/errs"
)

// Repo is an interface for repositories that can get metrics.
type Repo interface {
	GetMetrics(ctx context.Context) ([]*model.MetricsDto, error)
}

// fileDumper is a struct to dump data to file.
type fileDumper struct {
	filePath string
}

// NewDumper creates a new instance of fileDumper.
func NewDumper(filePath string) *fileDumper {
	return &fileDumper{
		filePath: filePath,
	}
}

// Dump data to file.
func (d *fileDumper) Dump(ctx context.Context, gaugeRepo Repo, counterRepo Repo) error {
	gaugeMetrics, err := gaugeRepo.GetMetrics(ctx)
	if err != nil {
		return errs.Wrap(err, "get gauge metrics")
	}
	counterMetrics, err := counterRepo.GetMetrics(ctx)
	if err != nil {
		return errs.Wrap(err, "get counter metrics")
	}

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
