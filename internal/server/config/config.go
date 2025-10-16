package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/yogenyslav/ya-metrics/pkg/errs"
)

// Server holds server configuration settings.
type Server struct {
	Port string `yaml:"port" env:"SERVER_PORT" env-default:"8080"`
}

// Config holds the overall application configuration.
type Config struct {
	Server Server `yaml:"server"`
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
