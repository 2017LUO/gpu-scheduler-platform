package handler

import (
	"net/http"
	"time"

	"gpu-scheduler-platform/internal/apiserver/dto"
	"gpu-scheduler-platform/internal/apiserver/service"
	model "gpu-scheduler-platform/internal/repo/models"

	"go.uber.org/zap"
)

type JobHandler struct {
	svc    *service.JobService
	logger *zap.Logger
}

func NewJobHandler(svc *service.JobService, lg *zap.Logger) *JobHandler {
	if lg == nil {
		lg = zap.NewNop()
	}
	return &JobHandler{svc: svc, logger: lg}
}

func (h *JobHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateJobRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, r, http.StatusBadRequest, 400, err.Error())
		return
	}

	job, err := h.svc.Create(r.Context(), req)
	if err != nil {
		status, code, msg := mapServiceError(err)
		writeError(w, r, status, code, msg)
		return
	}

	writeCreated(w, r, toJobResponse(job))
}

func (h *JobHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	job, err := h.svc.Get(r.Context(), id)
	if err != nil {
		status, code, msg := mapServiceError(err)
		writeError(w, r, status, code, msg)
		return
	}

	writeOK(w, r, toJobResponse(job))
}

func (h *JobHandler) List(w http.ResponseWriter, r *http.Request) {
	limit, offset := parsePaging(r, 50, 0)

	items, page, err := h.svc.List(r.Context(), service.ListJobsInput{
		TenantID:  r.URL.Query().Get("tenant_id"),
		Namespace: r.URL.Query().Get("namespace"),
		Queue:     r.URL.Query().Get("queue"),
		Status:    r.URL.Query().Get("status"),
		Limit:     limit,
		Offset:    offset,
	})
	if err != nil {
		status, code, msg := mapServiceError(err)
		writeError(w, r, status, code, msg)
		return
	}

	out := make([]dto.JobResponse, 0, len(items))
	for i := range items {
		out = append(out, toJobResponse(&items[i]))
	}

	writeList(w, r, out, page.Limit, page.Offset)
}

func (h *JobHandler) ListEvents(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	limit, offset := parsePaging(r, 100, 0)

	items, page, err := h.svc.ListEvents(r.Context(), id, limit, offset)
	if err != nil {
		status, code, msg := mapServiceError(err)
		writeError(w, r, status, code, msg)
		return
	}

	out := make([]dto.JobEventResponse, 0, len(items))
	for i := range items {
		out = append(out, toJobEventResponse(&items[i]))
	}

	writeList(w, r, out, page.Limit, page.Offset)
}

func toJobResponse(m *model.GPUJob) dto.JobResponse {
	return dto.JobResponse{
		ID:              m.ID,
		TenantID:        m.TenantID,
		Namespace:       m.Namespace,
		Name:            m.Name,
		Queue:           m.Queue,
		Priority:        m.Priority,
		Status:          m.Status,
		Submitter:       m.Submitter,
		GPUCount:        m.GPUCount,
		GPUMemoryMiB:    m.GPUMemoryMiB,
		GPUModel:        m.GPUModel,
		RequireSameNode: m.RequireSameNode,
		RequireHealthy:  m.RequireHealthy,
		RequireMIG:      m.RequireMIG,
		RetryCount:      m.RetryCount,
		Message:         m.Message,
		ScheduledAt:     formatPtrTime(m.ScheduledAt),
		StartedAt:       formatPtrTime(m.StartedAt),
		FinishedAt:      formatPtrTime(m.FinishedAt),
		CreatedAt:       m.CreatedAt.UTC().Format(time.RFC3339Nano),
		UpdatedAt:       m.UpdatedAt.UTC().Format(time.RFC3339Nano),
	}
}

func toJobEventResponse(m *model.GPUJobEvent) dto.JobEventResponse {
	return dto.JobEventResponse{
		ID:         m.ID,
		JobID:      m.JobID,
		TenantID:   m.TenantID,
		Reason:     m.Reason,
		Message:    m.Message,
		Source:     m.Source,
		OccurredAt: m.OccurredAt.UTC().Format(time.RFC3339Nano),
		CreatedAt:  m.CreatedAt.UTC().Format(time.RFC3339Nano),
	}
}

func formatPtrTime(t *time.Time) *string {
	if t == nil || t.IsZero() {
		return nil
	}
	v := t.UTC().Format(time.RFC3339Nano)
	return &v
}
