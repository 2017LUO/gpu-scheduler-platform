package handler

import (
	"net/http"
	"time"

	"gpu-scheduler-platform/internal/apiserver/dto"
	"gpu-scheduler-platform/internal/apiserver/service"
	model "gpu-scheduler-platform/internal/repo/models"

	"go.uber.org/zap"
)

type QuotaHandler struct {
	svc    *service.QuotaService
	logger *zap.Logger
}

func NewQuotaHandler(svc *service.QuotaService, lg *zap.Logger) *QuotaHandler {
	if lg == nil {
		lg = zap.NewNop()
	}
	return &QuotaHandler{svc: svc, logger: lg}
}

func (h *QuotaHandler) Upsert(w http.ResponseWriter, r *http.Request) {
	var req dto.UpsertQuotaRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, r, http.StatusBadRequest, 400, err.Error())
		return
	}

	item, err := h.svc.Upsert(r.Context(), req)
	if err != nil {
		status, code, msg := mapServiceError(err)
		writeError(w, r, status, code, msg)
		return
	}
	writeOK(w, r, toQuotaResponse(item))
}

func (h *QuotaHandler) Get(w http.ResponseWriter, r *http.Request) {
	tenantID := r.PathValue("tenantID")
	namespace := r.PathValue("namespace")

	item, err := h.svc.Get(r.Context(), tenantID, namespace)
	if err != nil {
		status, code, msg := mapServiceError(err)
		writeError(w, r, status, code, msg)
		return
	}
	writeOK(w, r, toQuotaResponse(item))
}

func (h *QuotaHandler) ListByTenant(w http.ResponseWriter, r *http.Request) {
	tenantID := r.PathValue("tenantID")
	items, err := h.svc.ListByTenant(r.Context(), tenantID)
	if err != nil {
		status, code, msg := mapServiceError(err)
		writeError(w, r, status, code, msg)
		return
	}

	out := make([]dto.QuotaResponse, 0, len(items))
	for i := range items {
		out = append(out, toQuotaResponse(&items[i]))
	}
	writeList(w, r, out, len(out), 0)
}

func (h *QuotaHandler) Delete(w http.ResponseWriter, r *http.Request) {
	tenantID := r.PathValue("tenantID")
	namespace := r.PathValue("namespace")

	if err := h.svc.Delete(r.Context(), tenantID, namespace); err != nil {
		status, code, msg := mapServiceError(err)
		writeError(w, r, status, code, msg)
		return
	}
	writeOK(w, r, dto.Empty{})
}

func toQuotaResponse(m *model.GPUQuota) dto.QuotaResponse {
	return dto.QuotaResponse{
		ID:              m.ID,
		TenantID:        m.TenantID,
		Namespace:       m.Namespace,
		MaxGPUCount:     m.MaxGPUCount,
		MaxRunningJobs:  m.MaxRunningJobs,
		MaxQueuedJobs:   m.MaxQueuedJobs,
		MaxGPUMemoryMiB: m.MaxGPUMemoryMiB,
		Enabled:         m.Enabled,
		CreatedAt:       m.CreatedAt.UTC().Format(time.RFC3339Nano),
		UpdatedAt:       m.UpdatedAt.UTC().Format(time.RFC3339Nano),
	}
}
