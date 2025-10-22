package server

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/yogenyslav/ya-metrics/internal/config"
	"github.com/yogenyslav/ya-metrics/internal/server/handler"
	"github.com/yogenyslav/ya-metrics/internal/server/repository"
	"github.com/yogenyslav/ya-metrics/internal/server/service"
)

// Server serves HTTP requests.
type Server struct {
	router chi.Router
	cfg    *config.Config
}

// NewServer creates new HTTP server.
func NewServer(cfg *config.Config) *Server {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	return &Server{
		router: router,
		cfg:    cfg,
	}
}

// Start starts the HTTP server.
func (s *Server) Start() error {
	gaugeRepo := repository.NewMetricInMemRepo[float64]()
	counterRepo := repository.NewMetricInMemRepo[int64]()

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
	if err := http.ListenAndServe(s.cfg.Addr, s.router); err != nil {
		panic(err)
	}
}
