package mysql

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gpu-scheduler-platform/internal/domain/allocation"
	"gpu-scheduler-platform/internal/domain/cluster"
	"gpu-scheduler-platform/internal/domain/event"
	"gpu-scheduler-platform/internal/domain/job"
	"gpu-scheduler-platform/internal/domain/policy"
	"gpu-scheduler-platform/internal/repo/models"
	"gpu-scheduler-platform/internal/util"

	"gorm.io/gorm"
)

func dbFromContext(ctx context.Context, fallback *gorm.DB) *gorm.DB {
	if fallback == nil {
		return nil
	}
	if ctx == nil {
		return fallback
	}
	if tx, ok := ctx.Value(txContextKey{}).(*gorm.DB); ok && tx != nil {
		return tx
	}
	return fallback
}

func parseJSONMap(s string) map[string]string {
	if s == "" {
		return nil
	}
	var out map[string]string
	if err := json.Unmarshal([]byte(s), &out); err != nil {
		return nil
	}
	return out
}

func mustJSON(v any) string {
	if v == nil {
		return ""
	}
	b, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(b)
}

func timePtr(t *time.Time) *time.Time {
	if t == nil {
		return nil
	}
	v := *t
	return &v
}

func copyStringSlice(in []string) []string {
	if len(in) == 0 {
		return nil
	}
	out := make([]string, len(in))
	copy(out, in)
	return out
}

func toJobModel(j *job.Job) *models.GPUJob {
	if j == nil {
		return nil
	}

	return &models.GPUJob{
		ID:                      j.ID,
		TenantID:                j.TenantID,
		Namespace:               j.Namespace,
		Name:                    j.Name,
		Queue:                   string(j.Queue),
		Priority:                string(j.Priority),
		Status:                  string(j.Status),
		GPUCount:                j.Requirement.GPUCount,
		GPUMemoryMiB:            j.Requirement.GPUMemoryMiB,
		GPUModel:                j.Requirement.GPUModel,
		RequireSameNode:         j.Requirement.RequireSameNode,
		RequireHealthy:          j.Requirement.RequireHealthy,
		RequireMIG:              j.Requirement.RequireMIG,
		MIGProfile:              j.Requirement.MIGProfile,
		RequireNVLink:           j.Requirement.RequireNVLink,
		Preemptible:             j.Requirement.Preemptible,
		Retryable:               j.Requirement.Retryable,
		MaxRetry:                j.Requirement.MaxRetry,
		ExpectedDuration:        int64(j.Requirement.ExpectedDuration.Seconds()),
		PreferredNodeLabelsJSON: mustJSON(j.Requirement.PreferredNodeLabels),
		PreferredGPULabelsJSON:  mustJSON(j.Requirement.PreferredGPULabels),
		LabelsJSON:              mustJSON(j.Labels),
		AnnotationsJSON:         mustJSON(j.Annotations),
		RetryCount:              j.RetryCount,
		Message:                 j.Message,
		ScheduledAt:             timePtr(j.ScheduledAt),
		StartedAt:               timePtr(j.StartedAt),
		FinishedAt:              timePtr(j.FinishedAt),
		CreatedAt:               j.CreatedAt,
		UpdatedAt:               j.UpdatedAt,
	}
}

func toJobDomain(m *models.GPUJob) *job.Job {
	if m == nil {
		return nil
	}

	return &job.Job{
		ID:        m.ID,
		TenantID:  m.TenantID,
		Namespace: m.Namespace,
		Name:      m.Name,
		Queue:     job.QueueName(m.Queue),
		Priority:  job.PriorityClass(m.Priority),
		Status:    job.Status(m.Status),
		Requirement: job.Requirement{
			GPUCount:            m.GPUCount,
			GPUMemoryMiB:        m.GPUMemoryMiB,
			GPUModel:            m.GPUModel,
			RequireSameNode:     m.RequireSameNode,
			RequireHealthy:      m.RequireHealthy,
			RequireMIG:          m.RequireMIG,
			MIGProfile:          m.MIGProfile,
			RequireNVLink:       m.RequireNVLink,
			PreferredNodeLabels: parseJSONMap(m.PreferredNodeLabelsJSON),
			PreferredGPULabels:  parseJSONMap(m.PreferredGPULabelsJSON),
			Preemptible:         m.Preemptible,
			Retryable:           m.Retryable,
			MaxRetry:            m.MaxRetry,
			ExpectedDuration:    time.Duration(m.ExpectedDuration) * time.Second,
		},
		Labels:      parseJSONMap(m.LabelsJSON),
		Annotations: parseJSONMap(m.AnnotationsJSON),
		RetryCount:  m.RetryCount,
		Message:     m.Message,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		ScheduledAt: timePtr(m.ScheduledAt),
		StartedAt:   timePtr(m.StartedAt),
		FinishedAt:  timePtr(m.FinishedAt),
	}
}

func toEventModel(e *event.Event) *models.GPUJobEvent {
	if e == nil {
		return nil
	}
	return &models.GPUJobEvent{
		ID:         e.ID,
		JobID:      e.JobID,
		TenantID:   e.TenantID,
		Reason:     string(e.Reason),
		Message:    e.Message,
		Source:     e.Source,
		OccurredAt: e.OccurredAt,
		CreatedAt:  e.OccurredAt,
	}
}

func toEventDomain(m *models.GPUJobEvent) event.Event {
	return event.Event{
		ID:         m.ID,
		JobID:      m.JobID,
		TenantID:   m.TenantID,
		Reason:     event.Reason(m.Reason),
		Message:    m.Message,
		Source:     m.Source,
		OccurredAt: m.OccurredAt,
	}
}

func toAllocationModel(a *allocation.Allocation) *models.Allocation {
	if a == nil {
		return nil
	}
	return &models.Allocation{
		ID:          a.ID,
		JobID:       a.JobID,
		TenantID:    a.TenantID,
		NodeName:    a.NodeName,
		GPUIDsJSON:  mustJSON(a.GPUIDs),
		Status:      string(a.Status),
		Message:     a.Message,
		CommittedAt: timePtr(a.CommittedAt),
		ReleasedAt:  timePtr(a.ReleasedAt),
		CreatedAt:   a.CreatedAt,
		UpdatedAt:   a.UpdatedAt,
	}
}

func toAllocationDomain(m *models.Allocation) *allocation.Allocation {
	if m == nil {
		return nil
	}
	var gpuIDs []string
	_ = json.Unmarshal([]byte(m.GPUIDsJSON), &gpuIDs)

	return &allocation.Allocation{
		ID:          m.ID,
		JobID:       m.JobID,
		TenantID:    m.TenantID,
		NodeName:    m.NodeName,
		GPUIDs:      copyStringSlice(gpuIDs),
		Status:      allocation.Status(m.Status),
		Message:     m.Message,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		CommittedAt: timePtr(m.CommittedAt),
		ReleasedAt:  timePtr(m.ReleasedAt),
	}
}

func toReservationModel(r *allocation.Reservation) *models.Reservation {
	if r == nil {
		return nil
	}
	return &models.Reservation{
		ID:         r.ID,
		JobID:      r.JobID,
		NodeName:   r.NodeName,
		GPUIDsJSON: mustJSON(r.GPUIDs),
		ExpireAt:   r.ExpireAt,
		CreatedAt:  r.CreatedAt,
	}
}

func toReservationDomain(m *models.Reservation) *allocation.Reservation {
	if m == nil {
		return nil
	}
	var gpuIDs []string
	_ = json.Unmarshal([]byte(m.GPUIDsJSON), &gpuIDs)

	return &allocation.Reservation{
		ID:        m.ID,
		JobID:     m.JobID,
		NodeName:  m.NodeName,
		GPUIDs:    copyStringSlice(gpuIDs),
		ExpireAt:  m.ExpireAt,
		CreatedAt: m.CreatedAt,
	}
}

func toBindingModel(b *allocation.Binding) *models.Binding {
	if b == nil {
		return nil
	}
	return &models.Binding{
		ID:         b.ID,
		JobID:      b.JobID,
		NodeName:   b.NodeName,
		GPUIDsJSON: mustJSON(b.GPUIDs),
		PodName:    b.PodName,
		Namespace:  b.Namespace,
		CreatedAt:  b.CreatedAt,
	}
}

func toBindingDomain(m *models.Binding) *allocation.Binding {
	if m == nil {
		return nil
	}
	var gpuIDs []string
	_ = json.Unmarshal([]byte(m.GPUIDsJSON), &gpuIDs)

	return &allocation.Binding{
		ID:        m.ID,
		JobID:     m.JobID,
		NodeName:  m.NodeName,
		GPUIDs:    copyStringSlice(gpuIDs),
		PodName:   m.PodName,
		Namespace: m.Namespace,
		CreatedAt: m.CreatedAt,
	}
}

func toQuotaModel(q *policy.Quota) *models.GPUQuota {
	if q == nil {
		return nil
	}
	return &models.GPUQuota{
		ID:              q.ID,
		TenantID:        q.TenantID,
		Namespace:       q.Namespace,
		MaxGPUCount:     q.MaxGPUCount,
		MaxRunningJobs:  q.MaxRunningJobs,
		MaxQueuedJobs:   q.MaxQueuedJobs,
		MaxGPUMemoryMiB: q.MaxGPUMemoryMiB,
		Enabled:         q.Enabled,
	}
}

func toQuotaDomain(m *models.GPUQuota) *policy.Quota {
	if m == nil {
		return nil
	}
	return &policy.Quota{
		ID:              m.ID,
		TenantID:        m.TenantID,
		Namespace:       m.Namespace,
		MaxGPUCount:     m.MaxGPUCount,
		MaxRunningJobs:  m.MaxRunningJobs,
		MaxQueuedJobs:   m.MaxQueuedJobs,
		MaxGPUMemoryMiB: m.MaxGPUMemoryMiB,
		Enabled:         m.Enabled,
	}
}

func isNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, util.ErrNotFound)
}

func wrapNotFound(err error, msg string) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("%s: %w", msg, util.ErrNotFound)
	}
	return fmt.Errorf("%s: %w", msg, err)
}

type nodeSnapshotPayload struct {
	Labels      map[string]string   `json:"labels,omitempty"`
	Annotations map[string]string   `json:"annotations,omitempty"`
	Topology    cluster.Topology    `json:"topology"`
	GPUs        []cluster.GPU       `json:"gpus,omitempty"`
	MIGs        []cluster.MIGDevice `json:"migs,omitempty"`
}
