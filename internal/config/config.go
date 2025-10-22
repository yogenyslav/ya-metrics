package config

import (
	"flag"
	"os"

	"github.com/yogenyslav/ya-metrics/pkg/errs"
)

const defaultServerAddr string = "localhost:8080"

// Config holds the configuration settings for the agent.
type Config struct {
	Addr string
}

// NewConfig creates a new Config with cli args or default values.
func NewConfig() (*Config, error) {
	flags := flag.NewFlagSet("server", flag.ExitOnError)
	addrFlag := flags.String("a", defaultServerAddr, "адрес сервера в формате ip:port")

	err := flags.Parse(os.Args[1:])
	if err != nil {
		return nil, errs.Wrap(err, "parse flags")
	}

	return &Config{
		Addr: *addrFlag,
	}, nil
}
