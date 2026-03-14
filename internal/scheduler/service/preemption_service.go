package service

import (
	"context"
	"fmt"
	"time"

	model "gpu-scheduler-platform/internal/repo/models"
	repoimpl "gpu-scheduler-platform/internal/repo/mysql"

	"go.uber.org/zap"
)

type PreemptionService struct {
	repos   *repoimpl.Repos
	logger  *zap.Logger
	nowFunc func() time.Time
}

func NewPreemptionService(repos *repoimpl.Repos, lg *zap.Logger) *PreemptionService {
	if lg == nil {
		lg = zap.NewNop()
	}
	return &PreemptionService{
		repos:   repos,
		logger:  lg,
		nowFunc: func() time.Time { return time.Now().UTC() },
	}
}

func (s *PreemptionService) TryPreempt(ctx context.Context, job *model.GPUJob, nodes []*model.Node) (bool, string, error) {
	if job == nil {
		return false, "", repoimpl.ErrInvalidArgument
	}
	if s.repos == nil || s.repos.Allocations == nil || s.repos.Bindings == nil || s.repos.GPUJobs == nil {
		return false, "", fmt.Errorf("preemption dependencies are incomplete")
	}

	for _, node := range nodes {
		if node == nil {
			continue
		}
		victims, ok, err := s.findVictimsOnNode(ctx, job, node)
		if err != nil {
			return false, "", err
		}
		if !ok {
			continue
		}

		msg := fmt.Sprintf("preempted %d job(s) on node %s for job %s", len(victims), node.NodeName, job.ID)
		for _, victim := range victims {
			if err := s.preemptVictim(ctx, job, victim); err != nil {
				return false, "", err
			}
		}
		return true, msg, nil
	}
	return false, "", nil
}

func (s *PreemptionService) findVictimsOnNode(ctx context.Context, job *model.GPUJob, node *model.Node) ([]*model.GPUJob, bool, error) {
	items, err := s.repos.Allocations.ListCommittedByNode(ctx, node.NodeName, repoimpl.PageQuery{Limit: 1000, Offset: 0})
	if err != nil {
		return nil, false, err
	}
	if len(items) == 0 {
		return nil, false, nil
	}

	requiredGPU := job.GPUCount
	requiredMem := int64(job.GPUCount) * job.GPUMemoryMiB

	freeGPU := node.HealthyGPUCount
	freeMem := node.FreeMemoryMiB

	if freeGPU >= requiredGPU && freeMem >= requiredMem {
		return nil, false, nil
	}

	victims := make([]*model.GPUJob, 0, len(items))
	for _, alloc := range items {
		victim, err := s.repos.GPUJobs.Get(ctx, alloc.JobID)
		if err != nil {
			continue
		}
		if !victim.Preemptible {
			continue
		}
		if normalizePriority(victim.Priority) >= normalizePriority(job.Priority) {
			continue
		}
		victims = append(victims, victim)
		freeGPU += victim.GPUCount
		freeMem += int64(victim.GPUCount) * victim.GPUMemoryMiB

		if freeGPU >= requiredGPU && freeMem >= requiredMem {
			return victims, true, nil
		}
	}
	return nil, false, nil
}

func (s *PreemptionService) preemptVictim(ctx context.Context, by *model.GPUJob, victim *model.GPUJob) error {
	now := s.nowFunc()
	message := fmt.Sprintf("preempted by higher priority job %s", by.ID)

	alloc, err := s.repos.Allocations.GetByJobID(ctx, victim.ID)
	if err == nil && alloc != nil {
		if err := s.repos.Allocations.MarkReleased(ctx, alloc.ID, now, &message); err != nil {
			return err
		}
	}
	_ = s.repos.Bindings.DeleteByJobID(ctx, victim.ID)
	_ = s.repos.Reservations.DeleteByJobID(ctx, victim.ID)

	if err := s.repos.GPUJobs.UpdateStatus(ctx, victim.ID, "QUEUED", &message); err != nil {
		return err
	}

	_ = s.repos.GPUJobEvents.Create(ctx, &model.GPUJobEvent{
		ID:         fmt.Sprintf("evt-%d", now.UnixNano()),
		JobID:      victim.ID,
		TenantID:   victim.TenantID,
		Reason:     "JOB_PREEMPTED",
		Message:    &message,
		Source:     "scheduler",
		OccurredAt: now,
		CreatedAt:  now,
	})

	_ = s.repos.Outbox.Create(ctx, &model.Outbox{
		Topic:       "job.preempted",
		EventKey:    victim.ID,
		PayloadJSON: mustJSON(map[string]any{"job_id": victim.ID, "by_job_id": by.ID}),
		Status:      "PENDING",
		AvailableAt: now,
		CreatedAt:   now,
		UpdatedAt:   now,
	})

	s.logger.Info("victim job preempted",
		zap.String("victim_job_id", victim.ID),
		zap.String("by_job_id", by.ID),
	)
	return nil
}

func normalizePriority(p string) int {
	switch p {
	case "CRITICAL":
		return 100
	case "HIGH":
		return 80
	case "MEDIUM":
		return 50
	case "NORMAL":
		return 40
	case "LOW":
		return 20
	default:
		return 10
	}
}
