package algorithm

import (
	"context"
	"fmt"
	"time"

	model "gpu-scheduler-platform/internal/repo/models"
	schedframework "gpu-scheduler-platform/internal/scheduler/framework"
)

type Outcome struct {
	Scheduled       bool
	Unschedulable   bool
	NeedsPreemption bool
	Message         string
}

func ScheduleOne(
	ctx context.Context,
	deps Dependencies,
	job *model.GPUJob,
	nodes []*model.Node,
) (Outcome, error) {
	if job == nil {
		return Outcome{}, fmt.Errorf("job is nil")
	}
	if deps.Repos == nil || deps.Framework == nil {
		return Outcome{}, fmt.Errorf("scheduler dependencies are incomplete")
	}
	if deps.Now == nil {
		deps.Now = func() time.Time { return time.Now().UTC() }
	}

	now := deps.Now()
	attemptNo := nextAttemptNo(ctx, deps.Repos, job.ID)

	cycleState := schedframework.NewCycleState()

	if st := deps.Framework.RunPreFilter(cycleState, job); !st.IsSuccess() {
		msg := firstReason(st, "prefilter rejected job")
		recordAttempt(ctx, deps.Repos, now, job, attemptNo, "", nil, nil, nil, "REJECTED", &msg)
		recordJobEvent(ctx, deps.Repos, now, job, "SCHEDULING_REJECTED", &msg)
		return Outcome{Unschedulable: true, Message: msg}, nil
	}

	selectedNode, candidateNodes, scores, filterReasons, st := SelectNode(ctx, deps, cycleState, job, nodes)
	if !st.IsSuccess() {
		msg := firstReason(st, "no feasible node")
		recordAttempt(ctx, deps.Repos, now, job, attemptNo, "", candidateNodes, scores, filterReasons, "UNSCHEDULABLE", &msg)
		recordJobEvent(ctx, deps.Repos, now, job, "SCHEDULING_FAILED", &msg)
		_ = deps.Repos.GPUJobs.UpdateStatus(ctx, job.ID, "QUEUED", &msg)
		return Outcome{
			Unschedulable:   true,
			NeedsPreemption: true,
			Message:         msg,
		}, nil
	}

	gpuUUIDs, st := SelectGPU(ctx, deps, cycleState, job, selectedNode)
	if !st.IsSuccess() {
		msg := firstReason(st, "no feasible gpu on selected node")
		recordAttempt(ctx, deps.Repos, now, job, attemptNo, selectedNode.NodeName, candidateNodes, scores, filterReasons, "UNSCHEDULABLE", &msg)
		recordJobEvent(ctx, deps.Repos, now, job, "SCHEDULING_FAILED", &msg)
		_ = deps.Repos.GPUJobs.UpdateStatus(ctx, job.ID, "QUEUED", &msg)
		return Outcome{
			Unschedulable:   true,
			NeedsPreemption: true,
			Message:         msg,
		}, nil
	}

	if st := Reserve(ctx, deps, cycleState, job, selectedNode); !st.IsSuccess() {
		msg := firstReason(st, "reservation failed")
		recordAttempt(ctx, deps.Repos, now, job, attemptNo, selectedNode.NodeName, candidateNodes, scores, filterReasons, "ERROR", &msg)
		recordJobEvent(ctx, deps.Repos, now, job, "RESERVATION_FAILED", &msg)
		return Outcome{Message: msg}, st.Err()
	}

	if st := deps.Framework.RunPermit(cycleState, job, selectedNode); !st.IsSuccess() {
		msg := firstReason(st, "permit rejected job")
		CleanupReservation(ctx, deps, job)
		recordAttempt(ctx, deps.Repos, now, job, attemptNo, selectedNode.NodeName, candidateNodes, scores, filterReasons, "WAIT", &msg)
		recordJobEvent(ctx, deps.Repos, now, job, "PERMIT_WAIT", &msg)
		return Outcome{Message: msg}, nil
	}

	if err := Commit(ctx, deps, cycleState, job, selectedNode, gpuUUIDs); err != nil {
		msg := err.Error()
		CleanupReservation(ctx, deps, job)
		recordAttempt(ctx, deps.Repos, now, job, attemptNo, selectedNode.NodeName, candidateNodes, scores, filterReasons, "ERROR", &msg)
		recordJobEvent(ctx, deps.Repos, now, job, "COMMIT_FAILED", &msg)
		return Outcome{Message: msg}, err
	}

	successMsg := fmt.Sprintf("scheduled on node %s with %d gpu(s)", selectedNode.NodeName, len(gpuUUIDs))
	recordAttempt(ctx, deps.Repos, now, job, attemptNo, selectedNode.NodeName, candidateNodes, scores, filterReasons, "SUCCESS", &successMsg)
	return Outcome{
		Scheduled: true,
		Message:   successMsg,
	}, nil
}

func firstReason(st *schedframework.Status, fallback string) string {
	if st == nil {
		return fallback
	}
	reasons := st.Reasons()
	if len(reasons) > 0 && reasons[0] != "" {
		return reasons[0]
	}
	if err := st.Err(); err != nil {
		return err.Error()
	}
	return fallback
}
