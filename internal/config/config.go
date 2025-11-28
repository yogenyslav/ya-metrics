package config

import (
	"flag"
	"os"

	"github.com/yogenyslav/ya-metrics/pkg"
	"github.com/yogenyslav/ya-metrics/pkg/errs"
	"github.com/yogenyslav/ya-metrics/pkg/retry"
)

const (
	defaultServerAddr       string = "localhost:8080"
	defaultStoreIntervalSec int    = 300
)

// DatabaseConfig holds the configuration settings for the database.
type DatabaseConfig struct {
	Dsn string
}

// ServerConfig holds the configuration settings for the server.
type ServerConfig struct {
	Addr     string
	LogLevel string
}

// DumpConfig holds settings for repository dumping into file.
type DumpConfig struct {
	FileStoragePath string
	StoreInterval   int
	Restore         bool
}

// Config holds the entire application settings.
type Config struct {
	Server *ServerConfig
	Dump   *DumpConfig
	DB     *DatabaseConfig
	Retry  *retry.Config
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
	dbDsnFlag := flags.String("d", "", "строка с адресом подключения к БД")

	err := flags.Parse(os.Args[1:])
	if err != nil {
		return nil, errs.Wrap(err, "parse flags")
	}

	return &Config{
		Server: &ServerConfig{
			Addr:     pkg.GetEnv("ADDRESS", *addrFlag),
			LogLevel: pkg.GetEnv("LOG_LEVEL", *logLevelFlag),
		},
		Dump: &DumpConfig{
			FileStoragePath: pkg.GetEnv("FILE_STORAGE_PATH", *fileStoragePathFlag),
			StoreInterval:   pkg.GetEnv("STORE_INTERVAL", *storeIntervalFlag),
			Restore:         pkg.GetEnv("RESTORE", *restoreFlag),
		},
		DB: &DatabaseConfig{
			Dsn: pkg.GetEnv("DATABASE_DSN", *dbDsnFlag),
		},
		Retry: &retry.Config{
			MaxRetries:         retry.DefaultRetries,
			LinearBackoffMilli: retry.DefaultLinearBackoffMilli,
		},
	}, nil
}
