package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/yogenyslav/ya-metrics/pkg/errs"
)

// Agent holds agent configuration settings.
type Agent struct {
	PollIntervalSec   int `yaml:"poll_interval" env:"POLL_INTERVAL" env-default:"2"`
	ReportIntervalSec int `yaml:"report_interval" env:"REPORT_INTERVAL" env-default:"10"`
}

// Config holds the overall application configuration.
type Config struct {
	Agent     Agent  `yaml:"agent"`
	ServerURL string `yaml:"server_url" env:"SERVER_URL" env-default:"http://localhost:8080"`
}

// MustNew creates a new Config or panics.
func MustNew(path ...string) *Config {
	var (
		err error
		cfg Config
	)

	if len(path) > 0 {
		err = cleanenv.ReadConfig(path[0], &cfg)
	} else {
		err = cleanenv.ReadEnv(&cfg)
	}

	if err != nil {
		panic(errs.Wrap(err, "failed to load config"))
	}

	return &cfg
}
