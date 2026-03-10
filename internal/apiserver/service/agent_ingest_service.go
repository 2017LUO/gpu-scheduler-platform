package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"gpu-scheduler-platform/internal/domain/cluster"
	"gpu-scheduler-platform/internal/repo"
	"gpu-scheduler-platform/internal/util"

	"github.com/google/uuid"
)

type AgentIngestService struct {
	snapshots repo.NodeSnapshotRepository
}

func NewAgentIngestService(snapshots repo.NodeSnapshotRepository) *AgentIngestService {
	return &AgentIngestService{
		snapshots: snapshots,
	}
}

type IngestAgentReportRequest struct {
	NodeName    string
	Timestamp   time.Time
	GPUs        []IngestGPUInfo
	MIGs        []IngestMIGInfo
	Topology    []IngestGPULink
	PodBindings []IngestPodGPUInfo
}

type IngestGPUInfo struct {
	ID            string
	UUID          string
	Index         int
	Model         string
	MemoryMiB     int64
	FreeMemoryMiB int64
	Healthy       bool
	Health        string
}

type IngestMIGInfo struct {
	ID         string
	ParentUUID string
	Profile    string
	MemoryMiB  int64
}

type IngestGPULink struct {
	From string
	To   string
	Type string
}

type IngestPodGPUInfo struct {
	PodName   string
	Namespace string
	GPUIDs    []string
}

type IngestAgentReportResult struct {
	Accepted        bool
	SnapshotVersion string
	NodeName        string
	GPUCount        int
}

func (s *AgentIngestService) IngestReport(ctx context.Context, req IngestAgentReportRequest) (*IngestAgentReportResult, error) {
	if s == nil || s.snapshots == nil {
		return nil, util.ErrUnavailable
	}

	req.NodeName = strings.TrimSpace(req.NodeName)
	if req.NodeName == "" {
		return nil, util.ErrInvalidArgument
	}

	if req.Timestamp.IsZero() {
		req.Timestamp = time.Now().UTC()
	}

	nodes := []cluster.Node{
		buildNodeFromReport(req),
	}

	snapshot := &cluster.Snapshot{
		Version:   uuid.NewString(),
		Nodes:     nodes,
		CreatedAt: req.Timestamp,
	}

	if err := s.snapshots.UpsertSnapshot(ctx, snapshot); err != nil {
		return nil, fmt.Errorf("upsert snapshot from agent report: %w", err)
	}

	return &IngestAgentReportResult{
		Accepted:        true,
		SnapshotVersion: snapshot.Version,
		NodeName:        req.NodeName,
		GPUCount:        len(req.GPUs),
	}, nil
}

func buildNodeFromReport(req IngestAgentReportRequest) cluster.Node {
	gpus := make([]cluster.GPU, 0, len(req.GPUs))
	for _, g := range req.GPUs {
		health := cluster.GPUHealth(g.Health)
		if health == "" {
			if g.Healthy {
				health = cluster.GPUHealthHealthy
			} else {
				health = cluster.GPUHealthUnknown
			}
		}

		gpus = append(gpus, cluster.GPU{
			ID:            g.ID,
			UUID:          g.UUID,
			NodeName:      req.NodeName,
			Index:         g.Index,
			Model:         g.Model,
			Vendor:        "nvidia",
			Type:          cluster.GPUTypeFull,
			MemoryMiB:     g.MemoryMiB,
			FreeMemoryMiB: g.FreeMemoryMiB,
			Healthy:       g.Healthy,
			Health:        health,
			MIGEnabled:    false,
			MIGProfile:    "",
			Labels:        map[string]string{},
			Annotations:   map[string]string{},
			Allocated:     false,
			Reserved:      false,
		})
	}

	migs := make([]cluster.MIGDevice, 0, len(req.MIGs))
	for _, m := range req.MIGs {
		migs = append(migs, cluster.MIGDevice{
			ID:            m.ID,
			ParentGPUUUID: m.ParentUUID,
			NodeName:      req.NodeName,
			Profile:       m.Profile,
			MemoryMiB:     m.MemoryMiB,
			Healthy:       true,
			Allocated:     false,
			Reserved:      false,
		})
	}

	links := make([]cluster.TopologyLink, 0, len(req.Topology))
	for _, l := range req.Topology {
		linkType := cluster.LinkType(strings.ToLower(strings.TrimSpace(l.Type)))
		if linkType == "" {
			linkType = cluster.LinkUnknown
		}
		links = append(links, cluster.TopologyLink{
			FromGPU: l.From,
			ToGPU:   l.To,
			Type:    linkType,
			Weight:  10,
		})
	}

	return cluster.Node{
		Name:        req.NodeName,
		State:       cluster.NodeStateReady,
		Schedulable: true,
		Labels: map[string]string{
			"source": "agent",
		},
		Annotations: map[string]string{},
		GPUs:        gpus,
		MIGs:        migs,
		Topology: cluster.Topology{
			NodeName: req.NodeName,
			Links:    links,
		},
	}
}
