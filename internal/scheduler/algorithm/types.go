package algorithm

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	model "gpu-scheduler-platform/internal/repo/models"
	repoimpl "gpu-scheduler-platform/internal/repo/mysql"
	schedcache "gpu-scheduler-platform/internal/scheduler/cache"
	schedframework "gpu-scheduler-platform/internal/scheduler/framework"

	"go.uber.org/zap"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type LoadNodeInventoryFunc func(context.Context, []*model.Node) (schedframework.NodeGPUInventory, error)
type SelectNodeGPUsFunc func(context.Context, *model.GPUJob, *model.Node) ([]model.GPUDevice, error)

type Dependencies struct {
	DB               *gorm.DB
	Repos            *repoimpl.Repos
	Framework        *schedframework.Framework
	ReservationCache *schedcache.ReservationCache
	Logger           *zap.Logger
	Now              func() time.Time
	ReservationTTL   time.Duration
	TopK             int

	LoadNodeInventory LoadNodeInventoryFunc
	SelectNodeGPUs    SelectNodeGPUsFunc
}

type NodeScore struct {
	Node  *model.Node
	Score int64
}

type Decision struct {
	Job            *model.GPUJob
	Node           *model.Node
	GPUUUIDs       []string
	AttemptNo      int
	CandidateNodes []string
	Scores         map[string]int64
	FilterReasons  map[string][]string
}

func mustJSON(v any) datatypes.JSON {
	if v == nil {
		return datatypes.JSON([]byte("{}"))
	}
	b, err := json.Marshal(v)
	if err != nil || len(b) == 0 {
		return datatypes.JSON([]byte("{}"))
	}
	return datatypes.JSON(b)
}

func mustJSONStringSlice(v []string) datatypes.JSON {
	if v == nil {
		return datatypes.JSON([]byte("[]"))
	}
	b, err := json.Marshal(v)
	if err != nil || len(b) == 0 {
		return datatypes.JSON([]byte("[]"))
	}
	return datatypes.JSON(b)
}

func newID(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, time.Now().UTC().UnixNano())
}

func nextAttemptNo(ctx context.Context, repos *repoimpl.Repos, jobID string) int {
	if repos == nil || repos.SchedulingAttempts == nil || jobID == "" {
		return 1
	}
	items, err := repos.SchedulingAttempts.ListByJob(ctx, jobID, repoimpl.PageQuery{Limit: 1000, Offset: 0})
	if err != nil {
		return 1
	}
	return len(items) + 1
}

func recordAttempt(
	ctx context.Context,
	repos *repoimpl.Repos,
	now time.Time,
	job *model.GPUJob,
	attemptNo int,
	selectedNode string,
	candidateNodes []string,
	scores map[string]int64,
	filterReasons map[string][]string,
	result string,
	message *string,
) {
	if repos == nil || repos.SchedulingAttempts == nil || job == nil {
		return
	}
	_ = repos.SchedulingAttempts.Create(ctx, &model.SchedulingAttempt{
		JobID:              job.ID,
		TenantID:           job.TenantID,
		AttemptNo:          attemptNo,
		Phase:              "SCHEDULE",
		SelectedNode:       selectedNode,
		CandidateNodesJSON: mustJSON(candidateNodes),
		ScoresJSON:         mustJSON(scores),
		FilterReasonsJSON:  mustJSON(filterReasons),
		Result:             result,
		Message:            message,
		CreatedAt:          now,
	})
}

func recordJobEvent(
	ctx context.Context,
	repos *repoimpl.Repos,
	now time.Time,
	job *model.GPUJob,
	reason string,
	message *string,
) {
	if repos == nil || repos.GPUJobEvents == nil || job == nil || reason == "" {
		return
	}
	_ = repos.GPUJobEvents.Create(ctx, &model.GPUJobEvent{
		ID:         newID("evt"),
		JobID:      job.ID,
		TenantID:   job.TenantID,
		Reason:     reason,
		Message:    message,
		Source:     "scheduler",
		OccurredAt: now,
		CreatedAt:  now,
	})
}
