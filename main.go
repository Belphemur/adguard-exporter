package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/belphemur/adguard-exporter/internal/adguard"
	"github.com/belphemur/adguard-exporter/internal/config"
	"github.com/belphemur/adguard-exporter/internal/http"
	"github.com/belphemur/adguard-exporter/internal/metrics"
	"github.com/belphemur/adguard-exporter/internal/worker"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	metrics.Init()
	global, err := config.FromEnv()
	if err != nil {
		panic(err)
	}

	clients := []*adguard.Client{}
	for _, conf := range global.Configs {
		clients = append(clients, adguard.NewClient(conf))
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	http := http.NewHttp(global.Server.Debug, global.Server.Port)
	log.Printf("Running server version %s (%s), built at %s\n", version, commit, date)
	go http.Serve()
	go worker.Work(ctx, global.Server.Interval, clients)

	<-sigs
	if err := http.Stop(ctx); err != nil {
		panic(err)
	}
}
