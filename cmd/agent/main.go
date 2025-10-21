package main

import (
	"log"
	"net/http"

	"github.com/yogenyslav/ya-metrics/internal/agent"
	"github.com/yogenyslav/ya-metrics/internal/agent/config"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg, err := config.NewConfig()
	if err != nil {
		return err
	}

	a := agent.New(http.DefaultClient, cfg)

	err = a.Start()
	if err != nil {
		return err
	}

	return nil
}
