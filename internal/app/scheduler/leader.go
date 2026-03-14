package scheduler

import (
	"context"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"time"

	appcfg "gpu-scheduler-platform/internal/config"

	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
)

type LeaderRunner struct {
	cfg      *appcfg.SchedulerConfig
	logger   *zap.Logger
	runner   *Runner
	client   kubernetes.Interface
	identity string
	active   atomic.Bool
}

func NewLeaderRunner(
	cfg *appcfg.SchedulerConfig,
	lg *zap.Logger,
	runner *Runner,
	client kubernetes.Interface,
	identity string,
) *LeaderRunner {
	if lg == nil {
		lg = zap.NewNop()
	}
	if identity == "" {
		identity = defaultLeaderIdentity()
	}

	return &LeaderRunner{
		cfg:      cfg,
		logger:   lg,
		runner:   runner,
		client:   client,
		identity: identity,
	}
}

func (l *LeaderRunner) Run(ctx context.Context) error {
	if l == nil {
		return fmt.Errorf("leader runner is nil")
	}
	if l.cfg == nil {
		return fmt.Errorf("scheduler config is nil")
	}
	if l.runner == nil {
		return fmt.Errorf("runner is nil")
	}

	leCfg := l.cfg.Scheduler.LeaderElection
	if !leCfg.Enabled {
		l.logger.Info("leader election disabled, running scheduler directly")
		return l.runner.Run(ctx)
	}

	if l.client == nil {
		return fmt.Errorf("kubernetes client is nil")
	}
	if leCfg.LeaseName == "" {
		return fmt.Errorf("leader election lease_name is empty")
	}
	if leCfg.LeaseNamespace == "" {
		return fmt.Errorf("leader election lease_namespace is empty")
	}

	if !l.active.CompareAndSwap(false, true) {
		l.logger.Warn("leader runner already active")
		return nil
	}
	defer l.active.Store(false)

	lock, err := resourcelock.New(
		resourcelock.LeasesResourceLock,
		leCfg.LeaseNamespace,
		leCfg.LeaseName,
		l.client.CoreV1(),
		l.client.CoordinationV1(),
		resourcelock.ResourceLockConfig{
			Identity: l.identity,
		},
	)
	if err != nil {
		return fmt.Errorf("create leader election lock: %w", err)
	}

	leaseDuration := durationOrDefault(leCfg.LeaseDuration, 15*time.Second)
	renewDeadline := durationOrDefault(leCfg.RenewDeadline, 10*time.Second)
	retryPeriod := durationOrDefault(leCfg.RetryPeriod, 2*time.Second)

	if renewDeadline >= leaseDuration {
		return fmt.Errorf("invalid leader election config: renew_deadline must be less than lease_duration")
	}
	if retryPeriod >= renewDeadline {
		return fmt.Errorf("invalid leader election config: retry_period must be less than renew_deadline")
	}

	l.logger.Info("starting leader election",
		zap.String("identity", l.identity),
		zap.String("lease_name", leCfg.LeaseName),
		zap.String("lease_namespace", leCfg.LeaseNamespace),
		zap.Duration("lease_duration", leaseDuration),
		zap.Duration("renew_deadline", renewDeadline),
		zap.Duration("retry_period", retryPeriod),
	)

	leadCtx, cancelLead := context.WithCancel(ctx)
	defer cancelLead()

	var (
		runErr   error
		runErrMu sync.Mutex
	)

	elector, err := leaderelection.NewLeaderElector(leaderelection.LeaderElectionConfig{
		Lock:            lock,
		ReleaseOnCancel: true,
		LeaseDuration:   leaseDuration,
		RenewDeadline:   renewDeadline,
		RetryPeriod:     retryPeriod,
		Name:            leCfg.LeaseName,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: func(c context.Context) {
				l.logger.Info("scheduler leadership acquired",
					zap.String("identity", l.identity),
					zap.String("lease_name", leCfg.LeaseName),
					zap.String("lease_namespace", leCfg.LeaseNamespace),
				)

				if err := l.runner.Run(c); err != nil {
					runErrMu.Lock()
					runErr = err
					runErrMu.Unlock()

					l.logger.Error("scheduler runner exited with error",
						zap.String("identity", l.identity),
						zap.Error(err),
					)
					cancelLead()
					return
				}

				l.logger.Info("scheduler runner exited",
					zap.String("identity", l.identity),
				)
				cancelLead()
			},
			OnStoppedLeading: func() {
				l.logger.Warn("scheduler leadership lost",
					zap.String("identity", l.identity),
					zap.String("lease_name", leCfg.LeaseName),
					zap.String("lease_namespace", leCfg.LeaseNamespace),
				)
				cancelLead()
			},
			OnNewLeader: func(identity string) {
				if identity == "" {
					return
				}
				if identity == l.identity {
					l.logger.Debug("current instance is leader",
						zap.String("identity", identity),
					)
					return
				}
				l.logger.Info("new scheduler leader observed",
					zap.String("current_identity", l.identity),
					zap.String("leader_identity", identity),
				)
			},
		},
	})
	if err != nil {
		return fmt.Errorf("create leader elector: %w", err)
	}

	elector.Run(leadCtx)

	runErrMu.Lock()
	defer runErrMu.Unlock()
	return runErr
}

func durationOrDefault(v, def time.Duration) time.Duration {
	if v > 0 {
		return v
	}
	return def
}

func defaultLeaderIdentity() string {
	host, err := os.Hostname()
	if err != nil || host == "" {
		return fmt.Sprintf("scheduler-%d", time.Now().UnixNano())
	}
	return fmt.Sprintf("%s-%d", host, time.Now().UnixNano())
}
