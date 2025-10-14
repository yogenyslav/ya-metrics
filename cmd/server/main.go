package main

import (
	"github.com/yogenyslav/ya-metrics/internal/config"
	"github.com/yogenyslav/ya-metrics/internal/server"
)

func main() {
	srv := server.NewServer(config.MustNew())
	if err := srv.Start(); err != nil {
		panic(err)
	}
}
