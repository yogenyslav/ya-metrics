package main

import (
	"context"
	"net/http"

	"github.com/yogenyslav/ya-metrics/internal/agent"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	a, err := agent.New(http.DefaultClient)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = a.Start(ctx)
	if err != nil {
		return err
	}

	return nil
}
