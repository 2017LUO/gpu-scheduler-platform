package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gpu-scheduler-platform/internal/apiserver/dto"
	"gpu-scheduler-platform/internal/apiserver/service"
	"gpu-scheduler-platform/internal/domain/job"
	"gpu-scheduler-platform/internal/middleware"

	"go.uber.org/zap"
)

type JobHandler struct {
	svc *service.JobService
	lg  *zap.Logger
}

func NewJobHandler(svc *service.JobService, lg *zap.Logger) *JobHandler {
	return &JobHandler{
		svc: svc,
		lg:  lg,
	}
}

func (h *JobHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, r, http.StatusBadRequest, 400, "invalid json body")
		return
	}

	out, err := h.svc.Create(r.Context(), service.CreateJobRequest{
		TenantID:  req.TenantID,
		Namespace: req.Namespace,
		Name:      req.Name,
		Queue:     job.QueueName(strings.TrimSpace(req.Queue)),
		Priority:  job.PriorityClass(strings.TrimSpace(req.Priority)),
		Requirement: job.Requirement{
			GPUCount:            req.Requirement.GPUCount,
			GPUMemoryMiB:        req.Requirement.GPUMemoryMiB,
			GPUModel:            req.Requirement.GPUModel,
			RequireSameNode:     req.Requirement.RequireSameNode,
			RequireHealthy:      req.Requirement.RequireHealthy,
			RequireMIG:          req.Requirement.RequireMIG,
			MIGProfile:          req.Requirement.MIGProfile,
			RequireNVLink:       req.Requirement.RequireNVLink,
			PreferredNodeLabels: req.Requirement.PreferredNodeLabels,
			PreferredGPULabels:  req.Requirement.PreferredGPULabels,
			Preemptible:         req.Requirement.Preemptible,
			Retryable:           req.Requirement.Retryable,
			MaxRetry:            req.Requirement.MaxRetry,
			ExpectedDuration:    time.Duration(req.Requirement.ExpectedDurationSec) * time.Second,
		},
		Labels:      req.Labels,
		Annotations: req.Annotations,
	})
	if err != nil {
		handleServiceError(w, r, err)
		return
	}

	writeOK(w, r, http.StatusOK, toJobResponse(out))
}

func (h *JobHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	jobID := jobIDFromPath(r.URL.Path)
	if jobID == "" {
		writeErr(w, r, http.StatusBadRequest, 400, "invalid job id")
		return
	}

	out, err := h.svc.GetByID(r.Context(), jobID)
	if err != nil {
		handleServiceError(w, r, err)
		return
	}

	writeOK(w, r, http.StatusOK, toJobResponse(out))
}

func (h *JobHandler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	limit, _ := strconv.Atoi(defaultString(q.Get("limit"), "20"))
	offset, _ := strconv.Atoi(defaultString(q.Get("offset"), "0"))

	items, total, err := h.svc.List(r.Context(), service.ListJobsRequest{
		TenantID:  q.Get("tenant_id"),
		Namespace: q.Get("namespace"),
		Queue:     job.QueueName(strings.TrimSpace(q.Get("queue"))),
		Status:    job.Status(strings.TrimSpace(q.Get("status"))),
		Priority:  job.PriorityClass(strings.TrimSpace(q.Get("priority"))),
		Keyword:   q.Get("keyword"),
		Limit:     limit,
		Offset:    offset,
	})
	if err != nil {
		handleServiceError(w, r, err)
		return
	}

	respItems := make([]dto.JobResponse, 0, len(items))
	for i := range items {
		respItems = append(respItems, toJobResponse(&items[i]))
	}

	writeOK(w, r, http.StatusOK, dto.ListJobsResponse{
		Meta: dto.PageMeta{
			Limit:  limit,
			Offset: offset,
			Total:  total,
		},
		Items: respItems,
	})
}

func (h *JobHandler) ListEvents(w http.ResponseWriter, r *http.Request) {
	jobID := jobIDFromEventsPath(r.URL.Path)
	if jobID == "" {
		writeErr(w, r, http.StatusBadRequest, 400, "invalid job id")
		return
	}

	q := r.URL.Query()
	limit, _ := strconv.Atoi(defaultString(q.Get("limit"), "20"))
	offset, _ := strconv.Atoi(defaultString(q.Get("offset"), "0"))

	items, total, err := h.svc.ListEvents(r.Context(), jobID, limit, offset)
	if err != nil {
		handleServiceError(w, r, err)
		return
	}

	respItems := make([]dto.JobEventResponse, 0, len(items))
	for i := range items {
		respItems = append(respItems, dto.JobEventResponse{
			ID:         items[i].ID,
			JobID:      items[i].JobID,
			TenantID:   items[i].TenantID,
			Reason:     string(items[i].Reason),
			Message:    items[i].Message,
			Source:     items[i].Source,
			OccurredAt: items[i].OccurredAt.UTC().Format(time.RFC3339Nano),
		})
	}

	writeOK(w, r, http.StatusOK, dto.ListJobEventsResponse{
		Meta: dto.PageMeta{
			Limit:  limit,
			Offset: offset,
			Total:  total,
		},
		Items: respItems,
	})
}

func toJobResponse(j *job.Job) dto.JobResponse {
	return dto.JobResponse{
		ID:        j.ID,
		TenantID:  j.TenantID,
		Namespace: j.Namespace,
		Name:      j.Name,
		Queue:     string(j.Queue),
		Priority:  string(j.Priority),
		Status:    string(j.Status),
		Requirement: dto.JobRequirementDTO{
			GPUCount:            j.Requirement.GPUCount,
			GPUMemoryMiB:        j.Requirement.GPUMemoryMiB,
			GPUModel:            j.Requirement.GPUModel,
			RequireSameNode:     j.Requirement.RequireSameNode,
			RequireHealthy:      j.Requirement.RequireHealthy,
			RequireMIG:          j.Requirement.RequireMIG,
			MIGProfile:          j.Requirement.MIGProfile,
			RequireNVLink:       j.Requirement.RequireNVLink,
			PreferredNodeLabels: j.Requirement.PreferredNodeLabels,
			PreferredGPULabels:  j.Requirement.PreferredGPULabels,
			Preemptible:         j.Requirement.Preemptible,
			Retryable:           j.Requirement.Retryable,
			MaxRetry:            j.Requirement.MaxRetry,
			ExpectedDurationSec: int64(j.Requirement.ExpectedDuration.Seconds()),
		},
		Labels:      j.Labels,
		Annotations: j.Annotations,
		RetryCount:  j.RetryCount,
		Message:     j.Message,
		CreatedAt:   j.CreatedAt.UTC().Format(time.RFC3339Nano),
		UpdatedAt:   j.UpdatedAt.UTC().Format(time.RFC3339Nano),
		ScheduledAt: formatTimePtr(j.ScheduledAt),
		StartedAt:   formatTimePtr(j.StartedAt),
		FinishedAt:  formatTimePtr(j.FinishedAt),
	}
}

func formatTimePtr(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.UTC().Format(time.RFC3339Nano)
	return &s
}

func jobIDFromPath(path string) string {
	path = strings.TrimSuffix(path, "/")
	parts := strings.Split(path, "/")
	if len(parts) < 1 {
		return ""
	}
	return parts[len(parts)-1]
}

func jobIDFromEventsPath(path string) string {
	path = strings.TrimSuffix(path, "/")
	parts := strings.Split(path, "/")
	if len(parts) < 2 {
		return ""
	}
	if parts[len(parts)-1] != "events" {
		return ""
	}
	return parts[len(parts)-2]
}

func defaultString(v, def string) string {
	if strings.TrimSpace(v) == "" {
		return def
	}
	return v
}

func writeOK(w http.ResponseWriter, r *http.Request, status int, data any) {
	rid := middleware.RequestIDFromContext(r.Context())
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(dto.OK(rid, data))
}

func writeErr(w http.ResponseWriter, r *http.Request, status int, code int, message string) {
	rid := middleware.RequestIDFromContext(r.Context())
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(dto.Err(rid, code, message))
}

func handleServiceError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case service.IsInvalidArgument(err):
		writeErr(w, r, http.StatusBadRequest, 400, err.Error())
	case service.IsNotFound(err):
		writeErr(w, r, http.StatusNotFound, 404, err.Error())
	case service.IsUnavailable(err):
		writeErr(w, r, http.StatusServiceUnavailable, 503, err.Error())
	default:
		writeErr(w, r, http.StatusInternalServerError, 500, "internal server error")
	}
}
