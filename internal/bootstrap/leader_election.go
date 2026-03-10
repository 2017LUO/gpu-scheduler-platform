package bootstrap

import (
	"context"

	appcfg "gpu-scheduler-platform/internal/config"

	"go.uber.org/zap"
)

type LeaderElector interface {
	Run(ctx context.Context, onStartedLeading func(context.Context)) error
}

type noopLeaderElector struct {
	logger *zap.Logger
}

func (n *noopLeaderElector) Run(ctx context.Context, onStartedLeading func(context.Context)) error {
	if onStartedLeading != nil {
		onStartedLeading(ctx)
	}
	<-ctx.Done()
	return nil
}

func NewLeaderElector(_ appcfg.LeaderElectionConfig, lg *zap.Logger) LeaderElector {
	return &noopLeaderElector{logger: lg}
}
