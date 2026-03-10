package mysql

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"gpu-scheduler-platform/internal/domain/cluster"
	"gpu-scheduler-platform/internal/repo"
	"gpu-scheduler-platform/internal/repo/models"
	"gpu-scheduler-platform/internal/util"

	"gorm.io/gorm"
)

type NodeSnapshotRepo struct {
	db *gorm.DB
}

func NewNodeSnapshotRepo(db *gorm.DB) *NodeSnapshotRepo {
	return &NodeSnapshotRepo{db: db}
}

func (r *NodeSnapshotRepo) UpsertSnapshot(ctx context.Context, s *cluster.Snapshot) error {
	if r == nil || r.db == nil || s == nil {
		return util.ErrInvalidArgument
	}

	txDB := dbFromContext(ctx, r.db).WithContext(ctx)

	for _, n := range s.Nodes {
		payload := nodeSnapshotPayload{
			Labels:      n.Labels,
			Annotations: n.Annotations,
			Topology:    n.Topology,
			GPUs:        n.GPUs,
			MIGs:        n.MIGs,
		}

		snapshotModel := &models.NodeSnapshot{
			Version:         s.Version,
			NodeName:        n.Name,
			NodeState:       string(n.State),
			Schedulable:     n.Schedulable,
			LabelsJSON:      mustJSON(n.Labels),
			AnnotationsJSON: mustJSON(n.Annotations),
			TopologyJSON:    mustJSON(payload),
			CreatedAt:       s.CreatedAt,
		}

		if err := txDB.Create(snapshotModel).Error; err != nil {
			return fmt.Errorf("create node snapshot: %w", err)
		}

		for _, g := range n.GPUs {
			device := &models.GPUDevice{
				ID:              g.ID,
				SnapshotID:      snapshotModel.ID,
				NodeName:        n.Name,
				UUID:            g.UUID,
				GPUIndex:        g.Index,
				Model:           g.Model,
				Vendor:          g.Vendor,
				Type:            string(g.Type),
				MemoryMiB:       g.MemoryMiB,
				FreeMemoryMiB:   g.FreeMemoryMiB,
				Healthy:         g.Healthy,
				Health:          string(g.Health),
				MIGEnabled:      g.MIGEnabled,
				MIGProfile:      g.MIGProfile,
				LabelsJSON:      mustJSON(g.Labels),
				AnnotationsJSON: mustJSON(g.Annotations),
				Allocated:       g.Allocated,
				Reserved:        g.Reserved,
			}
			if err := txDB.Create(device).Error; err != nil {
				return fmt.Errorf("create gpu device: %w", err)
			}
		}
	}

	return nil
}

func (r *NodeSnapshotRepo) GetLatest(ctx context.Context) (*cluster.Snapshot, error) {
	if r == nil || r.db == nil {
		return nil, util.ErrInvalidArgument
	}

	txDB := dbFromContext(ctx, r.db).WithContext(ctx)

	var latest models.NodeSnapshot
	if err := txDB.Order("created_at DESC, id DESC").Take(&latest).Error; err != nil {
		return nil, wrapNotFound(err, "get latest node snapshot")
	}

	var snapshotRows []models.NodeSnapshot
	if err := txDB.Where("version = ?", latest.Version).Order("node_name ASC").Find(&snapshotRows).Error; err != nil {
		return nil, fmt.Errorf("list latest node snapshots: %w", err)
	}

	nodes := make([]cluster.Node, 0, len(snapshotRows))
	for _, row := range snapshotRows {
		nodes = append(nodes, nodeFromSnapshotRow(row))
	}

	return &cluster.Snapshot{
		Version:   latest.Version,
		Nodes:     nodes,
		CreatedAt: latest.CreatedAt,
	}, nil
}

func (r *NodeSnapshotRepo) List(ctx context.Context, filter repo.NodeSnapshotListFilter) ([]cluster.Node, int64, error) {
	if r == nil || r.db == nil {
		return nil, 0, util.ErrInvalidArgument
	}

	dbq := dbFromContext(ctx, r.db).WithContext(ctx).Model(&models.NodeSnapshot{})
	if filter.NodeName != "" {
		dbq = dbq.Where("node_name = ?", filter.NodeName)
	}

	var total int64
	if err := dbq.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count node snapshots: %w", err)
	}

	limit := filter.Limit
	if limit <= 0 {
		limit = 50
	}
	if limit > 500 {
		limit = 500
	}
	offset := filter.Offset
	if offset < 0 {
		offset = 0
	}

	var rows []models.NodeSnapshot
	if err := dbq.Order("created_at DESC, id DESC").Limit(limit).Offset(offset).Find(&rows).Error; err != nil {
		return nil, 0, fmt.Errorf("list node snapshots: %w", err)
	}

	out := make([]cluster.Node, 0, len(rows))
	for _, row := range rows {
		out = append(out, nodeFromSnapshotRow(row))
	}
	return out, total, nil
}

func nodeFromSnapshotRow(row models.NodeSnapshot) cluster.Node {
	n := cluster.Node{
		Name:        row.NodeName,
		State:       cluster.NodeState(row.NodeState),
		Schedulable: row.Schedulable,
		Labels:      parseJSONMap(row.LabelsJSON),
		Annotations: parseJSONMap(row.AnnotationsJSON),
	}

	if strings.TrimSpace(row.TopologyJSON) != "" {
		var payload nodeSnapshotPayload
		if err := json.Unmarshal([]byte(row.TopologyJSON), &payload); err == nil {
			n.Topology = payload.Topology
			if len(payload.GPUs) > 0 {
				n.GPUs = payload.GPUs
			}
			if len(payload.MIGs) > 0 {
				n.MIGs = payload.MIGs
			}
			if len(payload.Labels) > 0 && len(n.Labels) == 0 {
				n.Labels = payload.Labels
			}
			if len(payload.Annotations) > 0 && len(n.Annotations) == 0 {
				n.Annotations = payload.Annotations
			}
		}
	}

	return n
}
