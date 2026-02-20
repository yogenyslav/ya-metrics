package config

import (
	"flag"
	"os"
	"strings"

	"github.com/yogenyslav/ya-metrics/pkg"
	"github.com/yogenyslav/ya-metrics/pkg/errs"
	"github.com/yogenyslav/ya-metrics/pkg/retry"
)

const (
	defaultServerAddr     = "localhost:8080"
	defaultPollInterval   = 2
	defaultReportInterval = 10
	defaultBatchSize      = 3
)

// Config holds the configuration settings for the agent.
type Config struct {
	ServerAddr        string
	PollIntervalSec   int
	ReportIntervalSec int
	CompressionType   string
	Retry             *retry.Config
	SecureKey         string
	RateLimit         int
	BatchSize         int
}

// NewConfig creates a new Config with cli args or default values.
func NewConfig() (*Config, error) {
	flags := flag.NewFlagSet("agent", flag.ExitOnError)
	serverAddrFlag := flags.String("a", defaultServerAddr, "адрес сервера в формате ip:port")
	pollIntervalFlag := flags.Int("p", defaultPollInterval, "интервал опроса метрик, сек.")
	reportIntervalFlag := flags.Int("r", defaultReportInterval, "интервал отправки метрик на сервер, сек. ")
	compressionTypeFlag := flags.String("c", "", "тип сжатия при отправке метрик на сервер")
	secureKeyFlag := flags.String("k", "", "ключ для подписи сигнатуры сообщений")
	rateLimitFlag := flags.Int("l", 1, "максимальное число одновременных запросов к серверу")

	if err := flags.Parse(os.Args[1:]); err != nil {
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
		Retry: &retry.Config{
			MaxRetries:         retry.DefaultRetries,
			LinearBackoffMilli: retry.DefaultLinearBackoffMilli,
		},
		SecureKey: pkg.GetEnv("KEY", *secureKeyFlag),
		RateLimit: pkg.GetEnv("RATE_LIMIT", *rateLimitFlag),
		BatchSize: defaultBatchSize,
	}, nil
}
