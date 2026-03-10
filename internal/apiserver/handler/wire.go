package handler

import (
	"gpu-scheduler-platform/internal/apiserver/service"

	"go.uber.org/zap"
)

type Handlers struct {
	Job           *JobHandler
	InternalAgent *InternalAgentHandler
}

func NewHandlers(services *service.Services, lg *zap.Logger) *Handlers {
	if services == nil {
		return &Handlers{}
	}
	return &Handlers{
		Job:           NewJobHandler(services.Job, lg),
		InternalAgent: NewInternalAgentHandler(services.AgentIngest, lg),
	}
}
