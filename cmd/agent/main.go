package main

import (
	"net/http"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/yogenyslav/ya-metrics/internal/agent"
	"github.com/yogenyslav/ya-metrics/internal/agent/config"
	"github.com/yogenyslav/ya-metrics/pkg/secure"
)

func main() {
	if err := run(); err != nil {
		log.Fatal().Err(err).Msg("agent failed")
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.NewConfig()
	if err != nil {
		return err
	}

	var sg *secure.SignatureGenerator
	if cfg.SecureKey != "" {
		sg = secure.NewSignatureGenerator(cfg.SecureKey)
	}

	l := zerolog.New(os.Stdout).With().Timestamp().Logger()

	a := agent.New(http.DefaultClient, cfg, sg, &l)

	err = a.Start()
	if err != nil {
		return err
	}

	return nil
}
