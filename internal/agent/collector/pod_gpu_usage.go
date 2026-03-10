package collector

import "context"

type PodGPUInfo struct {
	PodName   string   `json:"pod_name"`
	Namespace string   `json:"namespace"`
	GPUIDs    []string `json:"gpu_ids"`
}

type PodGPUUsageCollector struct{}

func NewPodGPUUsageCollector() *PodGPUUsageCollector {
	return &PodGPUUsageCollector{}
}

func (c *PodGPUUsageCollector) Collect(ctx context.Context) ([]PodGPUInfo, error) {
	_ = ctx
	return nil, nil
}
