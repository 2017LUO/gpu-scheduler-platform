package main

import (
	"context"
	"flag"
	"os"

	"gpu-scheduler-platform/internal/app/agent"
	"gpu-scheduler-platform/internal/bootstrap"

	"go.uber.org/zap"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "configs/agent.yaml", "path to agent config file")
	flag.Parse()

	cfg, err := bootstrap.LoadAgentConfig(configPath)
	if err != nil {
		panic(err)
	}

	app, err := agent.New(cfg)
	if err != nil {
		panic(err)
	}

	root := context.Background()
	ctx, cancel := bootstrap.WithSignalCancel(root)
	defer cancel()

	if err := app.Start(ctx); err != nil {
		app.Logger.Error("agent exited with error", zap.Error(err))
		_ = app.Stop(context.Background())
		os.Exit(1)
	}

	if err := app.Stop(context.Background()); err != nil {
		app.Logger.Error("agent shutdown failed", zap.Error(err))
		os.Exit(1)
	}
}
