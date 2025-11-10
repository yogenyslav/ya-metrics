package server

import (
	"context"
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
	"github.com/yogenyslav/ya-metrics/pkg/errs"
)

// Server serves HTTP requests.
type Server struct {
	router chi.Router
	cfg    *config.Config
}

// NewServer creates new HTTP server.
func NewServer(cfg *config.Config, l *zerolog.Logger) *Server {
	router := chi.NewRouter()
	router.Use(middleware.WithLogging(l))
	router.Use(middleware.WithCompression(middleware.GzipCompression))

	return &Server{
		router: router,
		cfg:    cfg,
	}
}

// Start starts the HTTP server.
func (s *Server) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dumper := repository.NewDumper(s.cfg.Dump.FileStoragePath, s.cfg.Dump.StoreInterval)

	var (
		gaugeMetrics   repository.StorageState[float64]
		counterMetrics repository.StorageState[int64]
		err            error
	)

	if s.cfg.Dump.Restore {
		gaugeMetrics, counterMetrics, err = repository.RestoreMetrics(s.cfg.Dump.FileStoragePath)
		if err != nil {
			return errs.Wrap(err, "restore metrics")
		}
	}

	gaugeRepo := repository.NewMetricInMemRepo(gaugeMetrics)
	counterRepo := repository.NewMetricInMemRepo(counterMetrics)

	dumper.Start(ctx, gaugeRepo, counterRepo)
	s.router.Use(middleware.WithFileDumper(dumper, s.cfg.Dump.StoreInterval, gaugeRepo, counterRepo))

	metricService := service.NewService(gaugeRepo, counterRepo)

	h := handler.NewHandler(metricService)
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
