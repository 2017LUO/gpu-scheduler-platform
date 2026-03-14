package service

import (
	"context"
	"time"

	"gpu-scheduler-platform/internal/apiserver/dto"
	model "gpu-scheduler-platform/internal/repo/models"
	repoimpl "gpu-scheduler-platform/internal/repo/mysql"

	"go.uber.org/zap"
	"gorm.io/datatypes"
)

type InternalAgentService struct {
	repos   *repoimpl.Repos
	logger  *zap.Logger
	nowFunc func() time.Time
}

func NewInternalAgentService(repos *repoimpl.Repos, lg *zap.Logger) *InternalAgentService {
	if lg == nil {
		lg = zap.NewNop()
	}
	return &InternalAgentService{
		repos:   repos,
		logger:  lg,
		nowFunc: func() time.Time { return time.Now().UTC() },
	}
}

func (s *InternalAgentService) HandleHeartbeat(ctx context.Context, req dto.AgentHeartbeatRequest) (time.Time, error) {
	seenAt, err := parseTimeOrNow(req.SeenAt, s.nowFunc)
	if err != nil {
		return time.Time{}, err
	}
	if req.NodeName == "" || req.Status == "" {
		return time.Time{}, repoimpl.ErrInvalidArgument
	}

	hb := &model.NodeHeartbeat{
		NodeName:   req.NodeName,
		Status:     req.Status,
		Message:    req.Message,
		LastSeenAt: seenAt,
		UpdatedAt:  seenAt,
	}
	if err := s.repos.NodeHeartbeats.Upsert(ctx, hb); err != nil {
		return time.Time{}, err
	}

	_ = s.repos.Nodes.UpdateHeartbeatTime(ctx, req.NodeName, seenAt)

	return seenAt, nil
}

func (s *InternalAgentService) HandleReport(ctx context.Context, req dto.AgentReportRequest) (uint64, time.Time, error) {
	reportTime, err := parseTimeOrNow(req.ReportTime, s.nowFunc)
	if err != nil {
		return 0, time.Time{}, err
	}
	if req.NodeName == "" {
		return 0, time.Time{}, repoimpl.ErrInvalidArgument
	}

	node := &model.Node{
		NodeName:          req.NodeName,
		ClusterName:       req.ClusterName,
		Source:            defaultString(req.Source, "agent"),
		State:             defaultString(req.NodeState, "READY"),
		Schedulable:       req.Schedulable,
		GPUCount:          req.GPUCount,
		HealthyGPUCount:   req.HealthyGPUCount,
		TotalMemoryMiB:    req.TotalMemoryMiB,
		FreeMemoryMiB:     req.FreeMemoryMiB,
		LabelsJSON:        mustJSON(req.Labels),
		AnnotationsJSON:   mustJSON(req.Annotations),
		TopologyJSON:      mustJSON(req.Topology),
		LastReportTime:    &reportTime,
		LastHeartbeatTime: nil,
		UpdatedAt:         reportTime,
	}
	if err := s.repos.Nodes.Upsert(ctx, node); err != nil {
		return 0, time.Time{}, err
	}

	snapshot := &model.NodeSnapshot{
		Version:         defaultString(req.Version, "v1"),
		AgentVersion:    req.AgentVersion,
		ClusterName:     req.ClusterName,
		NodeName:        req.NodeName,
		Source:          defaultString(req.Source, "agent"),
		NodeState:       defaultString(req.NodeState, "READY"),
		Schedulable:     req.Schedulable,
		GPUCount:        req.GPUCount,
		HealthyGPUCount: req.HealthyGPUCount,
		TotalMemoryMiB:  req.TotalMemoryMiB,
		FreeMemoryMiB:   req.FreeMemoryMiB,
		LabelsJSON:      mustJSON(req.Labels),
		AnnotationsJSON: mustJSON(req.Annotations),
		TopologyJSON:    mustJSON(req.Topology),
		ReportTime:      reportTime,
		CreatedAt:       reportTime,
	}
	if err := s.repos.NodeSnapshots.Create(ctx, snapshot); err != nil {
		return 0, time.Time{}, err
	}

	gpuItems := make([]model.GPUDevice, 0, len(req.GPUs))
	for _, item := range req.GPUs {
		gpuItems = append(gpuItems, model.GPUDevice{
			SnapshotID:        snapshot.ID,
			NodeName:          req.NodeName,
			UUID:              item.UUID,
			GPUIndex:          item.GPUIndex,
			Model:             item.Model,
			Vendor:            defaultString(item.Vendor, "nvidia"),
			Type:              defaultString(item.Type, "GPU"),
			MemoryMiB:         item.MemoryMiB,
			FreeMemoryMiB:     item.FreeMemoryMiB,
			Healthy:           item.Healthy,
			Health:            defaultString(item.Health, "OK"),
			MIGEnabled:        item.MIGEnabled,
			MIGProfile:        item.MIGProfile,
			UtilizationGPU:    item.UtilizationGPU,
			UtilizationMemory: item.UtilizationMemory,
			Temperature:       item.Temperature,
			PowerWatts:        item.PowerWatts,
			LabelsJSON:        mustJSON(item.Labels),
			AnnotationsJSON:   mustJSON(item.Annotations),
			Allocated:         item.Allocated,
			Reserved:          item.Reserved,
		})
	}
	if err := s.repos.NodeSnapshots.BatchCreateGPUDevices(ctx, gpuItems); err != nil {
		return 0, time.Time{}, err
	}

	migItems := make([]model.GPUMIGDevice, 0, len(req.MIGs))
	for _, item := range req.MIGs {
		migItems = append(migItems, model.GPUMIGDevice{
			SnapshotID:    snapshot.ID,
			NodeName:      req.NodeName,
			ParentGPUUUID: item.ParentGPUUUID,
			MIGUUID:       item.MIGUUID,
			Profile:       item.Profile,
			MemoryMiB:     item.MemoryMiB,
			Healthy:       item.Healthy,
			Allocated:     item.Allocated,
			Reserved:      item.Reserved,
		})
	}
	if err := s.repos.NodeSnapshots.BatchCreateMIGDevices(ctx, migItems); err != nil {
		return 0, time.Time{}, err
	}

	runtimeItems := make([]model.PodGPUBindingRuntime, 0, len(req.RuntimeBindings))
	for _, item := range req.RuntimeBindings {
		runtimeItems = append(runtimeItems, model.PodGPUBindingRuntime{
			SnapshotID: snapshot.ID,
			NodeName:   req.NodeName,
			Namespace:  item.Namespace,
			PodName:    item.PodName,
			GPUIDsJSON: datatypes.JSON(mustJSON(item.GPUIDs)),
		})
	}
	if err := s.repos.NodeSnapshots.BatchCreateRuntimeBindings(ctx, runtimeItems); err != nil {
		return 0, time.Time{}, err
	}

	return snapshot.ID, reportTime, nil
}

func parseTimeOrNow(v string, nowFunc func() time.Time) (time.Time, error) {
	if v == "" {
		return nowFunc(), nil
	}
	t, err := time.Parse(time.RFC3339Nano, v)
	if err != nil {
		return time.Time{}, err
	}
	return t.UTC(), nil
}
