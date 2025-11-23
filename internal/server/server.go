package server

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/yogenyslav/ya-metrics/internal/config"
	"github.com/yogenyslav/ya-metrics/internal/server/handler"
	"github.com/yogenyslav/ya-metrics/internal/server/middleware"
	"github.com/yogenyslav/ya-metrics/internal/server/repository"
	"github.com/yogenyslav/ya-metrics/internal/server/service"
	"github.com/yogenyslav/ya-metrics/pkg/database"
	"github.com/yogenyslav/ya-metrics/pkg/errs"
)

type Dumper interface {
	middleware.Dumper
	Start(ctx context.Context, gaugeRepo repository.Repo, counterRepo repository.Repo)
}

// Server serves HTTP requests.
type Server struct {
	router chi.Router
	cfg    *config.Config
	pg     database.PostgresTxDB
	dumper Dumper
}

// NewServer creates new HTTP server.
func NewServer(
	cfg *config.Config,
	l *zerolog.Logger,
) (*Server, error) {
	router := chi.NewRouter()
	router.Use(middleware.WithLogging(l))
	router.Use(middleware.WithCompression(middleware.GzipCompression))

	srv := &Server{
		router: router,
		cfg:    cfg,
	}

	switch {
	case cfg.DB.Dsn != "":
		pg, err := database.NewPostgres(context.Background(), cfg.DB.Dsn)
		if err != nil {
			return nil, errs.Wrap(err, "connect to database")
		}
		srv.pg = pg

		err = database.RunMigration(pg, "postgres")
		if err != nil {
			return nil, errs.Wrap(err, "run migrations")
		}
	case cfg.Dump.FileStoragePath != "":
		dumper := repository.NewDumper(cfg.Dump.FileStoragePath, cfg.Dump.StoreInterval)
		srv.dumper = dumper
	}

	return srv, nil
}

// Start starts the HTTP server.
func (s *Server) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var (
		gaugeRepo   service.GaugeRepo
		counterRepo service.CounterRepo
		err         error
	)

	if s.pg == nil {
		gaugeRepo, counterRepo, err = s.initRepos(ctx)
		if err != nil {
			return errs.Wrap(err, "init repositories")
		}
	} else {
		defer s.pg.Close()

		gaugeRepo = repository.NewMetricPostgresRepo[float64](s.pg)
		counterRepo = repository.NewMetricPostgresRepo[int64](s.pg)
	}

	metricService := service.NewService(gaugeRepo, counterRepo)

	h := handler.NewHandler(metricService, s.pg)
	h.RegisterRoutes(s.router)

	go s.listen()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop

	return nil
}

func (s *Server) listen() {
	if err := http.ListenAndServe(s.cfg.Server.Addr, s.router); err != nil {
		panic(err)
	}
}

func (s *Server) initRepos(ctx context.Context) (service.GaugeRepo, service.CounterRepo, error) {
	var (
		gaugeMetrics   repository.StorageState[float64]
		counterMetrics repository.StorageState[int64]
		err            error
	)

	if s.cfg.Dump.Restore {
		gaugeMetrics, counterMetrics, err = repository.RestoreMetrics(s.cfg.Dump.FileStoragePath)
		if err != nil {
			return nil, nil, errs.Wrap(err, "restore metrics")
		}
	}

	gaugeRepo := repository.NewMetricInMemRepo(gaugeMetrics)
	counterRepo := repository.NewMetricInMemRepo(counterMetrics)

	if s.dumper != nil {
		dumpingGaugeRepo, ok := any(gaugeRepo).(repository.Repo)
		if !ok {
			return nil, nil, errors.New("gauge repo does not implement repository.Repo")
		}

		dumpingCounterRepo, ok := any(counterRepo).(repository.Repo)
		if !ok {
			return nil, nil, errors.New("counter repo does not implement repository.Repo")
		}

		s.dumper.Start(ctx, dumpingGaugeRepo, dumpingCounterRepo)
		s.router.Use(
			middleware.WithFileDumper(s.dumper, s.cfg.Dump.StoreInterval, dumpingGaugeRepo, dumpingCounterRepo),
		)
	}

	return gaugeRepo, counterRepo, nil
}
