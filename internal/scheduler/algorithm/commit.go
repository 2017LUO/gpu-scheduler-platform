package algorithm

import (
	"context"
	"fmt"
	"time"

	model "gpu-scheduler-platform/internal/repo/models"
	schedframework "gpu-scheduler-platform/internal/scheduler/framework"

	"gorm.io/gorm"
)

func Commit(
	ctx context.Context,
	deps Dependencies,
	cs *schedframework.CycleState,
	job *model.GPUJob,
	node *model.Node,
	gpuUUIDs []string,
) error {
	if deps.DB == nil {
		return fmt.Errorf("db is nil")
	}
	if job == nil || node == nil || len(gpuUUIDs) == 0 {
		return fmt.Errorf("invalid commit arguments")
	}

	now := deps.Now()
	if now.IsZero() {
		now = time.Now().UTC()
	}

	reservationID, _ := schedframework.ReadReservationID(cs)
	message := fmt.Sprintf("allocated on node %s", node.NodeName)

	allocationID := newID("alloc")
	bindingID := newID("bind")

	err := deps.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		allocation := &model.Allocation{
			ID:            allocationID,
			ReservationID: reservationID,
			JobID:         job.ID,
			TenantID:      job.TenantID,
			NodeName:      node.NodeName,
			GPUIDsJSON:    mustJSONStringSlice(gpuUUIDs),
			Status:        "COMMITTED",
			Message:       &message,
			CommittedAt:   &now,
			CreatedAt:     now,
			UpdatedAt:     now,
		}
		if err := tx.Create(allocation).Error; err != nil {
			return err
		}

		binding := &model.Binding{
			ID:           bindingID,
			AllocationID: allocationID,
			JobID:        job.ID,
			NodeName:     node.NodeName,
			GPUIDsJSON:   mustJSONStringSlice(gpuUUIDs),
			PodName:      job.Name,
			Namespace:    job.Namespace,
			CreatedAt:    now,
		}
		if err := tx.Create(binding).Error; err != nil {
			return err
		}

		if err := tx.Model(&model.GPUJob{}).
			Where("id = ?", job.ID).
			Updates(map[string]any{
				"status":       "ALLOCATED",
				"message":      &message,
				"scheduled_at": now,
				"updated_at":   now,
			}).Error; err != nil {
			return err
		}

		event := &model.GPUJobEvent{
			ID:         newID("evt"),
			JobID:      job.ID,
			TenantID:   job.TenantID,
			Reason:     "JOB_ALLOCATED",
			Message:    &message,
			Source:     "scheduler",
			OccurredAt: now,
			CreatedAt:  now,
		}
		if err := tx.Create(event).Error; err != nil {
			return err
		}

		outbox := &model.Outbox{
			Topic:       "job.allocated",
			EventKey:    job.ID,
			PayloadJSON: mustJSON(map[string]any{"job_id": job.ID, "node_name": node.NodeName, "gpu_ids": gpuUUIDs}),
			Status:      "PENDING",
			AvailableAt: now,
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		if err := tx.Create(outbox).Error; err != nil {
			return err
		}

		if reservationID != "" {
			if err := tx.Delete(&model.Reservation{}, "id = ?", reservationID).Error; err != nil {
				return err
			}
		} else {
			if err := tx.Delete(&model.Reservation{}, "job_id = ?", job.ID).Error; err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	if deps.ReservationCache != nil {
		deps.ReservationCache.Delete(job.ID)
	}
	return nil
}
