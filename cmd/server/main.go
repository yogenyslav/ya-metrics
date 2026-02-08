package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/ya-metrics/internal/config"
	"github.com/yogenyslav/ya-metrics/internal/server"
	"github.com/yogenyslav/ya-metrics/pkg/errs"
)

func main() {
	if err := run(); err != nil {
		log.Fatal().Err(err).Msg("server failed")
	}
}

func run() error {
	cfg, err := config.NewConfig()
	if err != nil {
		return errs.Wrap(err, "create config")
	}

	logLevel, err := zerolog.ParseLevel(cfg.Server.LogLevel)
	if err != nil {
		return errs.Wrap(err, "parse log level")
	}
	l := zerolog.New(os.Stdout).With().Timestamp().Logger().Level(logLevel)

	srv, err := server.NewServer(cfg, &l)
	if err != nil {
		return errs.Wrap(err, "create server")
	}

	err = srv.Start()
	if err != nil {
		return errs.Wrap(err, "start server")
	}

	return nil
}
