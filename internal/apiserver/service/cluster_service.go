package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	model "gpu-scheduler-platform/internal/repo/models"
	repoimpl "gpu-scheduler-platform/internal/repo/mysql"

	"go.uber.org/zap"
)

type ClusterService struct {
	repos  *repoimpl.Repos
	logger *zap.Logger
}

type ListNodesInput struct {
	ClusterName string
	State       string
	Schedulable *bool
	Limit       int
	Offset      int
}

type NodeDetail struct {
	Node            *model.Node
	Heartbeat       *model.NodeHeartbeat
	LatestSnapshot  *model.NodeSnapshot
	GPUDevices      []model.GPUDevice
	MIGDevices      []model.GPUMIGDevice
	RuntimeBindings []model.PodGPUBindingRuntime
}

func NewClusterService(repos *repoimpl.Repos, lg *zap.Logger) *ClusterService {
	if lg == nil {
		lg = zap.NewNop()
	}
	return &ClusterService{
		repos:  repos,
		logger: lg,
	}
}

func (s *ClusterService) ListNodes(ctx context.Context, in ListNodesInput) ([]model.Node, repoimpl.PageQuery, error) {
	page := repoimpl.PageQuery{
		Limit:  in.Limit,
		Offset: in.Offset,
	}.Normalize(100, 1000)

	items, err := s.repos.Nodes.List(ctx, in.ClusterName, in.State, in.Schedulable, page)
	if err != nil {
		return nil, page, err
	}
	return items, page, nil
}

func (s *ClusterService) GetNode(ctx context.Context, nodeName string) (*NodeDetail, error) {
	node, err := s.repos.Nodes.Get(ctx, nodeName)
	if err != nil {
		return nil, err
	}

	detail := &NodeDetail{Node: node}

	if hb, err := s.repos.NodeHeartbeats.Get(ctx, nodeName); err == nil {
		detail.Heartbeat = hb
	} else if !errors.Is(err, repoimpl.ErrNotFound) {
		return nil, err
	}

	snapshot, err := s.repos.NodeSnapshots.GetLatestByNode(ctx, nodeName)
	if err != nil {
		if errors.Is(err, repoimpl.ErrNotFound) {
			return detail, nil
		}
		return nil, err
	}
	detail.LatestSnapshot = snapshot

	gpus, err := s.repos.NodeSnapshots.ListGPUDevicesBySnapshot(ctx, snapshot.ID)
	if err != nil {
		return nil, err
	}
	detail.GPUDevices = gpus

	migs, err := s.repos.NodeSnapshots.ListMIGDevicesBySnapshot(ctx, snapshot.ID)
	if err != nil {
		return nil, err
	}
	detail.MIGDevices = migs

	runtimeBindings, err := s.repos.NodeSnapshots.ListRuntimeBindingsBySnapshot(ctx, snapshot.ID)
	if err != nil {
		return nil, err
	}
	detail.RuntimeBindings = runtimeBindings

	return detail, nil
}

func jsonToMapString(b []byte) map[string]string {
	if len(b) == 0 {
		return nil
	}
	var out map[string]string
	if err := json.Unmarshal(b, &out); err != nil {
		return nil
	}
	return out
}

func jsonToMapAny(b []byte) map[string]any {
	if len(b) == 0 {
		return nil
	}
	var out map[string]any
	if err := json.Unmarshal(b, &out); err != nil {
		return nil
	}
	return out
}

func jsonToStringSlice(b []byte) []string {
	if len(b) == 0 {
		return nil
	}
	var out []string
	if err := json.Unmarshal(b, &out); err != nil {
		return nil
	}
	return out
}

func toRFC3339Ptr(t *time.Time) *string {
	if t == nil || t.IsZero() {
		return nil
	}
	v := t.UTC().Format(time.RFC3339Nano)
	return &v
}
