package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/yogenyslav/ya-metrics/internal/agent"
	"github.com/yogenyslav/ya-metrics/internal/agent/config"
	buildinfo "github.com/yogenyslav/ya-metrics/pkg/build_info"
	"github.com/yogenyslav/ya-metrics/pkg/errs"
	"github.com/yogenyslav/ya-metrics/pkg/secure"
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
		log.Fatal().Err(err).Msg("agent failed")
	}
}

func run() error {
	cfg, err := config.NewConfig()
	if err != nil {
		return errs.Wrap(err, "create config")
	}

	var sg *secure.SignatureGenerator
	if cfg.SecureKey != "" {
		sg = secure.NewSignatureGenerator(cfg.SecureKey)
	}

	l := zerolog.New(os.Stdout).With().Timestamp().Logger()

	a := agent.New(http.DefaultClient, cfg, sg, &l)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = a.Start(ctx)
	if err != nil {
		return errs.Wrap(err, "start agent")
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	signal.Stop(stop)

	cancel()
	a.Shutdown()

	return nil
}
