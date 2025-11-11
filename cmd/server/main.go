package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/ya-metrics/internal/config"
	"github.com/yogenyslav/ya-metrics/internal/server"
)

func main() {
	if err := run(); err != nil {
		log.Fatal().Err(err).Msg("server failed")
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.NewConfig()
	if err != nil {
		return err
	}

	logLevel, err := zerolog.ParseLevel(cfg.Server.LogLevel)
	if err != nil {
		return err
	}
	l := zerolog.New(os.Stdout).With().Timestamp().Logger().Level(logLevel)

	srv := server.NewServer(cfg, &l)

	err = srv.Start()
	if err != nil {
		return err
	}

	return nil
}
