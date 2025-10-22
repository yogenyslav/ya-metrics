package config

import (
	"flag"
	"os"
	"strings"

	"github.com/yogenyslav/ya-metrics/pkg/errs"
)

const (
	defaultServerAddr     = "localhost:8080"
	defaultPollInterval   = 2
	defaultReportInterval = 10
)

// Config holds the configuration settings for the agent.
type Config struct {
	ServerAddr        string
	PollIntervalSec   int
	ReportIntervalSec int
}

// NewConfig creates a new Config with cli args or default values.
func NewConfig() (*Config, error) {
	flags := flag.NewFlagSet("agent", flag.ExitOnError)
	serverAddrFlag := flags.String("a", defaultServerAddr, "адрес сервера в формате ip:port")
	pollIntervalFlag := flags.Int("p", defaultPollInterval, "интервал опроса метрик, сек.")
	reportIntervalFlag := flags.Int("r", defaultReportInterval, "интервал отправки метрик на сервер, сек. ")

	err := flags.Parse(os.Args[1:])
	if err != nil {
		return nil, errs.Wrap(err, "parse flags")
	}

	serverAddr := *serverAddrFlag
	if !strings.HasPrefix(serverAddr, "http://") && !strings.HasPrefix(serverAddr, "https://") {
		serverAddr = "http://" + serverAddr
	}

	return &Config{
		ServerAddr:        serverAddr,
		PollIntervalSec:   *pollIntervalFlag,
		ReportIntervalSec: *reportIntervalFlag,
	}, nil
}
