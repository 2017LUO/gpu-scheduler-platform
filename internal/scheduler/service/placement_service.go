package service

import (
	"context"
	"sort"
	"strings"

	model "gpu-scheduler-platform/internal/repo/models"
	repoimpl "gpu-scheduler-platform/internal/repo/mysql"
	schedframework "gpu-scheduler-platform/internal/scheduler/framework"

	"go.uber.org/zap"
)

type PlacementService struct {
	repos  *repoimpl.Repos
	logger *zap.Logger
}

func NewPlacementService(repos *repoimpl.Repos, lg *zap.Logger) *PlacementService {
	if lg == nil {
		lg = zap.NewNop()
	}
	return &PlacementService{
		repos:  repos,
		logger: lg,
	}
}

func (s *PlacementService) LoadNodeInventory(ctx context.Context, nodes []*model.Node) (schedframework.NodeGPUInventory, error) {
	out := make(schedframework.NodeGPUInventory, len(nodes))
	for _, node := range nodes {
		if node == nil {
			continue
		}

		snapshot, err := s.repos.NodeSnapshots.GetLatestByNode(ctx, node.NodeName)
		if err != nil {
			continue
		}
		gpus, err := s.repos.NodeSnapshots.ListGPUDevicesBySnapshot(ctx, snapshot.ID)
		if err != nil {
			return nil, err
		}
		out[node.NodeName] = gpus
	}
	return out, nil
}

func (s *PlacementService) SelectNodeGPUs(ctx context.Context, job *model.GPUJob, node *model.Node) ([]model.GPUDevice, error) {
	if job == nil || node == nil {
		return nil, repoimpl.ErrInvalidArgument
	}

	snapshot, err := s.repos.NodeSnapshots.GetLatestByNode(ctx, node.NodeName)
	if err != nil {
		return nil, err
	}
	items, err := s.repos.NodeSnapshots.ListGPUDevicesBySnapshot(ctx, snapshot.ID)
	if err != nil {
		return nil, err
	}

	out := make([]model.GPUDevice, 0, len(items))
	for _, item := range items {
		if item.Allocated || item.Reserved {
			continue
		}
		if job.RequireHealthy && !item.Healthy {
			continue
		}
		if strings.TrimSpace(job.GPUModel) != "" && !strings.EqualFold(strings.TrimSpace(job.GPUModel), strings.TrimSpace(item.Model)) {
			continue
		}
		if job.RequireMIG && !item.MIGEnabled {
			continue
		}
		if strings.TrimSpace(job.MIGProfile) != "" && !strings.EqualFold(strings.TrimSpace(job.MIGProfile), strings.TrimSpace(item.MIGProfile)) {
			continue
		}
		if job.GPUMemoryMiB > 0 && item.FreeMemoryMiB < job.GPUMemoryMiB {
			continue
		}
		out = append(out, item)
	}

	sort.SliceStable(out, func(i, j int) bool {
		if out[i].FreeMemoryMiB != out[j].FreeMemoryMiB {
			return out[i].FreeMemoryMiB > out[j].FreeMemoryMiB
		}
		if out[i].UtilizationGPU != out[j].UtilizationGPU {
			return out[i].UtilizationGPU < out[j].UtilizationGPU
		}
		return out[i].GPUIndex < out[j].GPUIndex
	})

	return out, nil
}
