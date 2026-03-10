package service

import "gpu-scheduler-platform/internal/domain/job"

type FairnessService struct{}

func NewFairnessService() *FairnessService {
	return &FairnessService{}
}

func (s *FairnessService) OrderJobs(in []job.Job) []job.Job {
	// 第一版直接保持 repo 返回顺序。
	// 后续接 priority/fair queue 时在这里替换。
	out := make([]job.Job, len(in))
	copy(out, in)
	return out
}
