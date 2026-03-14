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

// App：agent 进程的应用层总装配与生命周期管理入口。
// 负责把配置、日志、K8s 客户端、tracing、metrics、HTTP server、ServiceRunner
// 组装为一个完整可运行的 agent 应用，并提供 Start / Stop / Readiness 能力。
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

	// 启动期尽早失败，避免错误配置拖到运行时
	if err := ValidateAgentConfig(cfg); err != nil {
		return nil, fmt.Errorf("validate agent config: %w", err)
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

	tracingCloser, err := bootstrap.InitTracing(cfg.Service, cfg.Observability.Tracing)
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

	serviceRunner, err := NewServiceRunner(cfg, lg, metricsReg, k8sClients)
	if err != nil {
		return nil, fmt.Errorf("init service runner: %w", err)
	}
	app.Service = serviceRunner

	app.readyChecks = []ReadyCheck{
		{
			Name: "config",
			Fn: func(ctx context.Context) error {
				_ = ctx
				return ValidateAgentConfig(app.Config)
			},
		},
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
		{
			Name: "service_runner",
			Fn: func(ctx context.Context) error {
				_ = ctx
				if app.Service == nil {
					return fmt.Errorf("service runner not initialized")
				}
				if app.Service.Service() == nil {
					return fmt.Errorf("agent service not initialized")
				}
				if app.Service.Reporter() == nil {
					return fmt.Errorf("agent reporter not initialized")
				}
				return nil
			},
		},
		{
			Name: "http_server",
			Fn: func(ctx context.Context) error {
				_ = ctx
				if app.HTTPServer == nil {
					return fmt.Errorf("http server not initialized")
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
	if a == nil {
		return fmt.Errorf("agent app is nil")
	}
	if a.Lifecycle == nil {
		return fmt.Errorf("agent lifecycle is nil")
	}
	if ctx == nil {
		ctx = context.Background()
	}

	if err := a.Lifecycle.Start(ctx); err != nil {
		return err
	}
	return a.run(ctx)
}

func (a *App) Stop(ctx context.Context) error {
	if a == nil {
		return nil
	}
	if ctx == nil {
		ctx = context.Background()
	}

	var err error
	if a.Lifecycle != nil {
		err = a.Lifecycle.Stop(ctx)
	}

	// 兜底关闭 ServiceRunner（例如 gRPC reporter 长连接）
	if a.Service != nil {
		if closeErr := a.Service.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}

	return err
}

func (a *App) Readiness(ctx context.Context) (bool, map[string]string) {
	if a == nil {
		return false, map[string]string{
			"app": "agent app is nil",
		}
	}

	if ctx == nil {
		ctx = context.Background()
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	details := make(map[string]string, len(a.readyChecks))
	ok := true

	for _, chk := range a.readyChecks {
		if chk.Fn == nil {
			details[chk.Name] = "skip"
			continue
		}
		if err := chk.Fn(ctx); err != nil {
			ok = false
			details[chk.Name] = err.Error()
		} else {
			details[chk.Name] = "ok"
		}
	}

	return ok, details
}
