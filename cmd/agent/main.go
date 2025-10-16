package main

import (
	"context"
	"net/http"

	"github.com/yogenyslav/ya-metrics/internal/agent"
	"github.com/yogenyslav/ya-metrics/internal/agent/config"
)

func main() {
	cfg := config.MustNew()
	a := agent.New(cfg, http.DefaultClient)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := a.Start(ctx)
	if err != nil {
		panic(err)
	}
}
