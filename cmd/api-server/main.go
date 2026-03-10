package main

import (
	"context"
	"flag"
	"os"

	"gpu-scheduler-platform/internal/app/apiserver"
	"gpu-scheduler-platform/internal/bootstrap"

	"go.uber.org/zap"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "configs/api-server.yaml", "path to api-server config file")
	flag.Parse()

	cfg, err := bootstrap.LoadAPIServerConfig(configPath)
	if err != nil {
		panic(err)
	}

	app, err := apiserver.New(cfg)
	if err != nil {
		panic(err)
	}

	root := context.Background()
	ctx, cancel := bootstrap.WithSignalCancel(root)
	defer cancel()

	if err := app.Start(ctx); err != nil {
		app.Logger.Error("api-server exited with error", zap.Error(err))
		_ = app.Stop(context.Background())
		os.Exit(1)
	}

	if err := app.Stop(context.Background()); err != nil {
		app.Logger.Error("api-server shutdown failed", zap.Error(err))
		os.Exit(1)
	}
}
