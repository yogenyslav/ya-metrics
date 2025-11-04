package config

import (
	"flag"
	"os"

	"github.com/yogenyslav/ya-metrics/pkg"
	"github.com/yogenyslav/ya-metrics/pkg/errs"
)

const (
	defaultServerAddr       string = "localhost:8080"
	defaultStoreIntervalSec int    = 300
)

// Config holds the configuration settings for the server.
type Config struct {
	Addr            string
	LogLevel        string
	FileStoragePath string
	StoreInterval   int
	Restore         bool
}

// NewConfig creates a new Config with cli args or default values.
func NewConfig() (*Config, error) {
	flags := flag.NewFlagSet("server", flag.ExitOnError)
	addrFlag := flags.String("a", defaultServerAddr, "адрес сервера в формате ip:port")
	logLevelFlag := flags.String("l", "debug", "уровень логирования (debug, info, error)")
	fileStoragePathFlag := flags.String("f", "metrics.json", "путь к файлу для хранения метрик")
	storeIntervalFlag := flags.Int(
		"i",
		defaultStoreIntervalSec,
		"интервал сохранения метрик в файл в секундах (значение 0 делает запись синхронной)",
	)
	restoreFlag := flags.Bool("r", false, "восстановление метрик из файла при старте сервера")

	err := flags.Parse(os.Args[1:])
	if err != nil {
		return nil, errs.Wrap(err, "parse flags")
	}

	return &Config{
		Addr:            pkg.GetEnv("ADDRESS", *addrFlag),
		LogLevel:        pkg.GetEnv("LOG_LEVEL", *logLevelFlag),
		FileStoragePath: pkg.GetEnv("FILE_STORAGE_PATH", *fileStoragePathFlag),
		StoreInterval:   pkg.GetEnv("STORE_INTERVAL", *storeIntervalFlag),
		Restore:         pkg.GetEnv("RESTORE", *restoreFlag),
	}, nil
}
