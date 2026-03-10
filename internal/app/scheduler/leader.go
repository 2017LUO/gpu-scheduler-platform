package scheduler

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

func (a *App) runWithLeaderElection(ctx context.Context) error {
	if a.Config == nil {
		return fmt.Errorf("scheduler config is nil")
	}
	if a.LeaderElector == nil {
		return fmt.Errorf("leader elector is nil")
	}
	if a.Runner == nil {
		return fmt.Errorf("runner is nil")
	}

	a.Logger.Info("scheduler starting",
		zap.Bool("leader_election_enabled", a.Config.LeaderElection.Enabled),
		zap.String("lease_name", a.Config.LeaderElection.LeaseName),
		zap.String("lease_namespace", a.Config.LeaderElection.LeaseNamespace),
	)

	return a.LeaderElector.Run(ctx, func(runCtx context.Context) {
		a.Logger.Info("scheduler became leader")
		if err := a.Runner.Run(runCtx); err != nil {
			a.Logger.Error("scheduler runner exited with error", zap.Error(err))
		}
	})
}
