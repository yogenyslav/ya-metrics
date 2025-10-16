package server

import (
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/yogenyslav/ya-metrics/internal/server/config"
	"github.com/yogenyslav/ya-metrics/internal/server/handler"
	"github.com/yogenyslav/ya-metrics/internal/server/repository"
	"github.com/yogenyslav/ya-metrics/internal/server/service"
)

// Server serves HTTP requests.
type Server struct {
	mux *http.ServeMux
	cfg *config.Config
}

// NewServer creates new HTTP server.
func NewServer(cfg *config.Config) *Server {
	return &Server{
		mux: http.NewServeMux(),
		cfg: cfg,
	}
}

// Start starts the HTTP server.
func (s *Server) Start() error {
	gaugeRepo := repository.NewMetricInMemRepo[float64]()
	counterRepo := repository.NewMetricInMemRepo[int64]()

	metricService := service.NewService(gaugeRepo, counterRepo)

	h := handler.NewHandler(metricService)
	h.RegisterRoutes(s.mux)

	go s.listen()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop

	return nil
}

func (s *Server) listen() {
	addr := net.JoinHostPort("", s.cfg.Server.Port)
	log.Println("starting server at", addr)
	if err := http.ListenAndServe(addr, s.mux); err != nil {
		panic(err)
	}
}
