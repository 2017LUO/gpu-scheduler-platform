package algorithm

import (
	"context"

	model "gpu-scheduler-platform/internal/repo/models"
	schedframework "gpu-scheduler-platform/internal/scheduler/framework"
)

func Reserve(
	ctx context.Context,
	deps Dependencies,
	cs *schedframework.CycleState,
	job *model.GPUJob,
	node *model.Node,
) *schedframework.Status {
	if deps.Framework == nil {
		return schedframework.NewStatus(schedframework.CodeError, "framework is nil")
	}
	return deps.Framework.RunReserve(cs, job, node)
}

func CleanupReservation(ctx context.Context, deps Dependencies, job *model.GPUJob) {
	if job == nil {
		return
	}
	if deps.ReservationCache != nil {
		deps.ReservationCache.Delete(job.ID)
	}
	if deps.Repos != nil && deps.Repos.Reservations != nil {
		_ = deps.Repos.Reservations.DeleteByJobID(ctx, job.ID)
	}
}
