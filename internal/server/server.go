package server

import (
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/yogenyslav/ya-metrics/internal/server/handler"
	"github.com/yogenyslav/ya-metrics/internal/server/repository"
	"github.com/yogenyslav/ya-metrics/internal/server/service"
	"github.com/yogenyslav/ya-metrics/pkg/errs"
)

const defaultServerAddr string = "localhost:8080"

// Server serves HTTP requests.
type Server struct {
	router chi.Router
	addr   string
}

// NewServer creates new HTTP server.
func NewServer() (*Server, error) {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	flags := flag.NewFlagSet("server", flag.ExitOnError)
	addrFlag := flags.String("a", defaultServerAddr, "адрес сервера в формате ip:port")

	err := flags.Parse(os.Args[1:])
	if err != nil {
		return nil, errs.Wrap(err, "parse flags")
	}

	return &Server{
		router: router,
		addr:   *addrFlag,
	}, nil
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
	if err := http.ListenAndServe(s.addr, s.router); err != nil {
		panic(err)
	}
}
