package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/ya-metrics/internal/config"
	"github.com/yogenyslav/ya-metrics/internal/server"
	buildinfo "github.com/yogenyslav/ya-metrics/pkg/build_info"
	"github.com/yogenyslav/ya-metrics/pkg/errs"
)

// To set build info, use the following ldflags:
// -ldflags "-X main.buildVersion=$(VERSION) -X main.buildDate=$(DATE) -X main.buildCommit=$(COMMIT)".
var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	buildinfo.GetInfo(buildVersion, buildDate, buildCommit)
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = srv.Start(ctx)
	if err != nil {
		return errs.Wrap(err, "start server")
	}
	defer srv.Shutdown()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop
	signal.Stop(stop)

	cancel()

	return nil
}
