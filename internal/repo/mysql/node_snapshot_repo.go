package mysql

import (
	"context"
	"fmt"
	model "gpu-scheduler-platform/internal/repo/models"

	"gorm.io/gorm"
)

type NodeSnapshotRepo struct {
	db *gorm.DB
}

func NewNodeSnapshotRepo(db *gorm.DB) (*NodeSnapshotRepo, error) {
	if db == nil {
		return nil, ErrNilDB
	}
	return &NodeSnapshotRepo{db: db}, nil
}

func (r *NodeSnapshotRepo) Create(ctx context.Context, m *model.NodeSnapshot) error {
	if r == nil || r.db == nil || m == nil || m.NodeName == "" || m.ReportTime.IsZero() {
		return ErrInvalidArgument
	}
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return fmt.Errorf("create node snapshot: %w", err)
	}
	return nil
}

func (r *NodeSnapshotRepo) BatchCreateGPUDevices(ctx context.Context, items []model.GPUDevice) error {
	if r == nil || r.db == nil {
		return ErrNilDB
	}
	if len(items) == 0 {
		return nil
	}
	if err := r.db.WithContext(ctx).CreateInBatches(items, 200).Error; err != nil {
		return fmt.Errorf("batch create gpu devices: %w", err)
	}
	return nil
}

func (r *NodeSnapshotRepo) BatchCreateMIGDevices(ctx context.Context, items []model.GPUMIGDevice) error {
	if r == nil || r.db == nil {
		return ErrNilDB
	}
	if len(items) == 0 {
		return nil
	}
	if err := r.db.WithContext(ctx).CreateInBatches(items, 200).Error; err != nil {
		return fmt.Errorf("batch create mig devices: %w", err)
	}
	return nil
}

func (r *NodeSnapshotRepo) BatchCreateRuntimeBindings(ctx context.Context, items []model.PodGPUBindingRuntime) error {
	if r == nil || r.db == nil {
		return ErrNilDB
	}
	if len(items) == 0 {
		return nil
	}
	if err := r.db.WithContext(ctx).CreateInBatches(items, 200).Error; err != nil {
		return fmt.Errorf("batch create runtime bindings: %w", err)
	}
	return nil
}

func (r *NodeSnapshotRepo) GetLatestByNode(ctx context.Context, nodeName string) (*model.NodeSnapshot, error) {
	if r == nil || r.db == nil || nodeName == "" {
		return nil, ErrInvalidArgument
	}
	var m model.NodeSnapshot
	if err := r.db.WithContext(ctx).
		Where("node_name = ?", nodeName).
		Order("report_time DESC, id DESC").
		First(&m).Error; err != nil {
		return nil, mapDBError(err)
	}
	return &m, nil
}

func (r *NodeSnapshotRepo) ListGPUDevicesBySnapshot(ctx context.Context, snapshotID uint64) ([]model.GPUDevice, error) {
	if r == nil || r.db == nil || snapshotID == 0 {
		return nil, ErrInvalidArgument
	}
	var out []model.GPUDevice
	if err := r.db.WithContext(ctx).
		Where("snapshot_id = ?", snapshotID).
		Order("gpu_index ASC, id ASC").
		Find(&out).Error; err != nil {
		return nil, fmt.Errorf("list gpu devices by snapshot: %w", err)
	}
	return out, nil
}

func (r *NodeSnapshotRepo) ListMIGDevicesBySnapshot(ctx context.Context, snapshotID uint64) ([]model.GPUMIGDevice, error) {
	if r == nil || r.db == nil || snapshotID == 0 {
		return nil, ErrInvalidArgument
	}
	var out []model.GPUMIGDevice
	if err := r.db.WithContext(ctx).
		Where("snapshot_id = ?", snapshotID).
		Order("parent_gpu_uuid ASC, mig_uuid ASC, id ASC").
		Find(&out).Error; err != nil {
		return nil, fmt.Errorf("list mig devices by snapshot: %w", err)
	}
	return out, nil
}

func (r *NodeSnapshotRepo) ListRuntimeBindingsBySnapshot(ctx context.Context, snapshotID uint64) ([]model.PodGPUBindingRuntime, error) {
	if r == nil || r.db == nil || snapshotID == 0 {
		return nil, ErrInvalidArgument
	}
	var out []model.PodGPUBindingRuntime
	if err := r.db.WithContext(ctx).
		Where("snapshot_id = ?", snapshotID).
		Order("namespace ASC, pod_name ASC, id ASC").
		Find(&out).Error; err != nil {
		return nil, fmt.Errorf("list runtime bindings by snapshot: %w", err)
	}
	return out, nil
}
