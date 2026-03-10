package agent

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"gpu-scheduler-platform/internal/bootstrap"
	appcfg "gpu-scheduler-platform/internal/config"
	obslog "gpu-scheduler-platform/internal/observability/logging"
	obsmetrics "gpu-scheduler-platform/internal/observability/metrics"

	"go.uber.org/zap"
)

type App struct {
	Config        *appcfg.AgentConfig
	Logger        *zap.Logger
	K8s           *bootstrap.K8sClients
	Metrics       *obsmetrics.Registry
	HTTPServer    *http.Server
	Lifecycle     *bootstrap.Lifecycle
	TracingCloser bootstrap.TracingCloser
	Service       *ServiceRunner
	readyChecks   []ReadyCheck
}

type ReadyCheck struct {
	Name string
	Fn   func(context.Context) error
}

func New(cfg *appcfg.AgentConfig) (*App, error) {
	if cfg == nil {
		return nil, fmt.Errorf("agent config is nil")
	}

	lg, err := bootstrap.NewLogger(cfg.Logging)
	if err != nil {
		return nil, fmt.Errorf("init logger: %w", err)
	}
	lg = obslog.WithComponent(lg, "agent")

	k8sClients, err := bootstrap.NewKubernetesClients(cfg.Kubernetes)
	if err != nil {
		return nil, fmt.Errorf("init kubernetes clients: %w", err)
	}

	tracingCloser, err := bootstrap.InitTracing(cfg.Observability.Tracing)
	if err != nil {
		return nil, fmt.Errorf("init tracing: %w", err)
	}

	metricsReg := obsmetrics.NewRegistry()
	lifecycle := bootstrap.NewLifecycle(lg)

	app := &App{
		Config:        cfg,
		Logger:        lg,
		K8s:           k8sClients,
		Metrics:       metricsReg,
		Lifecycle:     lifecycle,
		TracingCloser: tracingCloser,
	}

	app.Service = NewServiceRunner(cfg, lg, k8sClients)

	app.readyChecks = []ReadyCheck{
		{
			Name: "kubernetes",
			Fn: func(ctx context.Context) error {
				_ = ctx
				if app.K8s == nil || app.K8s.Clientset == nil {
					return fmt.Errorf("kubernetes client not initialized")
				}
				return nil
			},
		},
	}

	mux := app.buildMux()
	app.HTTPServer = bootstrap.NewHTTPServer(app.metricsHTTPConfig(), mux)
	app.registerLifecycleHooks()

	return app, nil
}

func (a *App) Start(ctx context.Context) error {
	if err := a.Lifecycle.Start(ctx); err != nil {
		return err
	}
	return a.run(ctx)
}

func (a *App) Stop(ctx context.Context) error {
	return a.Lifecycle.Stop(ctx)
}

func (a *App) Readiness(ctx context.Context) (bool, map[string]string) {
	if ctx == nil {
		ctx = context.Background()
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	details := make(map[string]string, len(a.readyChecks))
	ok := true
	for _, chk := range a.readyChecks {
		if err := chk.Fn(ctx); err != nil {
			ok = false
			details[chk.Name] = err.Error()
		} else {
			details[chk.Name] = "ok"
		}
	}
	return ok, details
}
