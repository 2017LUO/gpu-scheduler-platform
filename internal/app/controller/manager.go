package controller

import (
	"context"

	appcfg "gpu-scheduler-platform/internal/config"
	ctrlsvc "gpu-scheduler-platform/internal/controller/service"
	obslog "gpu-scheduler-platform/internal/observability/logging"
	repoimpl "gpu-scheduler-platform/internal/repo/mysql"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Manager struct {
	cfg     *appcfg.ControllerAppConfig
	logger  *zap.Logger
	service *ctrlsvc.ControllerService
}

func NewManager(
	cfg *appcfg.ControllerAppConfig,
	lg *zap.Logger,
	db *gorm.DB,
	_ *redis.Client,
	_ any,
) *Manager {
	repos := repoimpl.NewRepositories(db)

	service := ctrlsvc.NewControllerService(
		repos.Jobs,
		repos.JobEvents,
		repos.Reservations,
		repos.Allocations,
		repos.TxManager,
		lg,
		cfg.Controller.ResyncPeriod,
	)

	return &Manager{
		cfg:     cfg,
		logger:  obslog.LoggerOrNop(lg),
		service: service,
	}
}

func (m *Manager) Run(ctx context.Context) error {
	return m.service.Run(ctx)
}
