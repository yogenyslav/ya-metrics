package main

import (
	"github.com/yogenyslav/ya-metrics/internal/server"
	"github.com/yogenyslav/ya-metrics/internal/server/config"
)

func main() {
	srv := server.NewServer(config.MustNew())
	if err := srv.Start(); err != nil {
		panic(err)
	}
}
