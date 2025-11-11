package config

import (
	"flag"
	"os"
	"strings"

	"github.com/yogenyslav/ya-metrics/pkg"
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
	CompressionType   string
}

// NewConfig creates a new Config with cli args or default values.
func NewConfig() (*Config, error) {
	flags := flag.NewFlagSet("agent", flag.ExitOnError)
	serverAddrFlag := flags.String("a", defaultServerAddr, "адрес сервера в формате ip:port")
	pollIntervalFlag := flags.Int("p", defaultPollInterval, "интервал опроса метрик, сек.")
	reportIntervalFlag := flags.Int("r", defaultReportInterval, "интервал отправки метрик на сервер, сек. ")
	compressionTypeFlag := flags.String("c", "", "тип сжатия при отправке метрик на сервер")

	err := flags.Parse(os.Args[1:])
	if err != nil {
		return nil, errs.Wrap(err, "parse flags")
	}

	serverAddr := pkg.GetEnv("ADDRESS", *serverAddrFlag)
	if !strings.HasPrefix(serverAddr, "http://") && !strings.HasPrefix(serverAddr, "https://") {
		serverAddr = "http://" + serverAddr
	}

	return &Config{
		ServerAddr:        serverAddr,
		PollIntervalSec:   pkg.GetEnv("POLL_INTERVAL", *pollIntervalFlag),
		ReportIntervalSec: pkg.GetEnv("REPORT_INTERVAL", *reportIntervalFlag),
		CompressionType:   pkg.GetEnv("COMPRESSION_TYPE", *compressionTypeFlag),
	}, nil
}
