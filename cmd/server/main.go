package main

import (
	"github.com/yogenyslav/ya-metrics/internal/server"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	srv, err := server.NewServer()
	if err != nil {
		return err
	}

	err = srv.Start()
	if err != nil {
		return err
	}

	return nil
}
