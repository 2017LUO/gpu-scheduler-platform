package algorithm

import (
	"sort"

	"gpu-scheduler-platform/internal/domain/cluster"
	"gpu-scheduler-platform/internal/domain/job"
)

func SelectGPUs(n cluster.Node, req job.Requirement) ([]cluster.GPU, bool) {
	candidates := make([]cluster.GPU, 0, len(n.GPUs))
	for _, g := range n.GPUs {
		if !gpuFits(g, req) {
			continue
		}
		candidates = append(candidates, g)
	}

	if len(candidates) < req.GPUCount {
		return nil, false
	}

	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].FreeMemoryMiB == candidates[j].FreeMemoryMiB {
			return candidates[i].Index < candidates[j].Index
		}
		return candidates[i].FreeMemoryMiB < candidates[j].FreeMemoryMiB
	})

	out := make([]cluster.GPU, 0, req.GPUCount)
	for _, g := range candidates {
		out = append(out, g)
		if len(out) == req.GPUCount {
			return out, true
		}
	}
	return nil, false
}

func gpuFits(g cluster.GPU, req job.Requirement) bool {
	if g.Allocated || g.Reserved {
		return false
	}
	if req.RequireHealthy && !g.Healthy {
		return false
	}
	if req.GPUMemoryMiB > 0 && g.FreeMemoryMiB < req.GPUMemoryMiB {
		return false
	}
	if req.GPUModel != "" && g.Model != req.GPUModel {
		return false
	}
	if req.RequireMIG && !g.MIGEnabled {
		return false
	}
	if req.RequireMIG && req.MIGProfile != "" && g.MIGProfile != req.MIGProfile {
		return false
	}
	return true
}
