package service

import (
	repoimpl "gpu-scheduler-platform/internal/repo/mysql"
)

type Services struct {
	Job         *JobService
	AgentIngest *AgentIngestService
}

func NewServices(repos *repoimpl.Repositories) *Services {
	if repos == nil {
		return &Services{}
	}
	return &Services{
		Job: NewJobService(
			repos.Jobs,
			repos.JobEvents,
			repos.Tenants,
			repos.TxManager,
		),
		AgentIngest: NewAgentIngestService(
			repos.NodeSnapshots,
		),
	}
}
