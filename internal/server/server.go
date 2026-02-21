package server

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/ya-metrics/internal/config"
	"github.com/yogenyslav/ya-metrics/internal/server/audit"
	"github.com/yogenyslav/ya-metrics/internal/server/handler"
	"github.com/yogenyslav/ya-metrics/internal/server/middleware"
	"github.com/yogenyslav/ya-metrics/internal/server/repository"
	"github.com/yogenyslav/ya-metrics/internal/server/service"
	"github.com/yogenyslav/ya-metrics/pkg/database"
	"github.com/yogenyslav/ya-metrics/pkg/errs"
)

// Server serves HTTP requests.
type Server struct {
	router         chi.Router
	cfg            *config.Config
	pg             database.TxDB
	dumper         middleware.Dumper
	dumpOnShutdown func()
}

// NewServer creates new HTTP server.
func NewServer(cfg *config.Config, l *zerolog.Logger) (*Server, error) {
	router := chi.NewRouter()
	router.Use(
		middleware.WithLogging(l),
		middleware.WithCompression(middleware.GzipCompression),
		middleware.WithSignature(cfg.Server.SecureKey),
	)

	srv := &Server{
		router: router,
		cfg:    cfg,
	}

	switch {
	case cfg.DB.Dsn != "":
		pg, err := database.NewPostgres(context.Background(), cfg.DB.Dsn, cfg.Retry)
		if err != nil {
			return nil, errs.Wrap(err, "connect to database")
		}
		srv.pg = pg

		err = database.RunMigration(pg, "postgres")
		if err != nil {
			return nil, errs.Wrap(err, "run migrations")
		}
	case cfg.Dump.FileStoragePath != "":
		srv.dumper = repository.NewDumper(cfg.Dump.FileStoragePath)
	}

	return srv, nil
}

// Start starts the HTTP server.
func (s *Server) Start(ctx context.Context) error {
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
		gaugeRepo = repository.NewMetricPostgresRepo[float64](s.pg)
		counterRepo = repository.NewMetricPostgresRepo[int64](s.pg)
	}
	s.router.Mount("/debug", chimw.Profiler())

	metricService := service.NewService(gaugeRepo, counterRepo, database.NewUnitOfWork(s.pg))
	audit := audit.New(s.cfg.Audit)

	h := handler.NewHandler(metricService, s.pg, audit)
	h.RegisterRoutes(s.router)

	go s.listen()

	return nil
}

func (s *Server) listen() {
	if err := http.ListenAndServe(s.cfg.Server.Addr, s.router); err != nil {
		log.Err(err).Msg("failed serving HTTP")
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

		s.Dumping(ctx, s.dumper, defaultTickerFactory, dumpingGaugeRepo, dumpingCounterRepo)
		s.dumpOnShutdown = func() {
			err := s.dumper.Dump(context.Background(), gaugeRepo, counterRepo)
			if err != nil {
				log.Warn().Err(err).Msg("failed to dump data on shutdown")
			} else {
				log.Info().Msg("successfuly dumped data on shutdown")
			}
		}
		s.router.Use(
			middleware.WithFileDumper(s.dumper, s.cfg.Dump.StoreInterval, dumpingGaugeRepo, dumpingCounterRepo),
		)
	}

	return gaugeRepo, counterRepo, nil
}

// Shutdown performs server shutdown.
func (s *Server) Shutdown() {
	if s.pg != nil {
		s.pg.Close()
	}
	s.dumpOnShutdown()
}
