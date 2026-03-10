package main

import (
	"context"
	"flag"
	"os"

	"gpu-scheduler-platform/internal/app/scheduler"
	"gpu-scheduler-platform/internal/bootstrap"

	"go.uber.org/zap"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "configs/scheduler.yaml", "path to scheduler config file")
	flag.Parse()

	cfg, err := bootstrap.LoadSchedulerConfig(configPath)
	if err != nil {
		panic(err)
	}

	app, err := scheduler.New(cfg)
	if err != nil {
		panic(err)
	}

	root := context.Background()
	ctx, cancel := bootstrap.WithSignalCancel(root)
	defer cancel()

	if err := app.Start(ctx); err != nil {
		app.Logger.Error("scheduler exited with error", zap.Error(err))
		_ = app.Stop(context.Background())
		os.Exit(1)
	}

	if err := app.Stop(context.Background()); err != nil {
		app.Logger.Error("scheduler shutdown failed", zap.Error(err))
		os.Exit(1)
	}
}
