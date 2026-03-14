package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"gpu-scheduler-platform/internal/apiserver/dto"
	"gpu-scheduler-platform/internal/apiserver/service"
	model "gpu-scheduler-platform/internal/repo/models"

	"go.uber.org/zap"
)

type PolicyHandler struct {
	svc    *service.PolicyService
	logger *zap.Logger
}

func NewPolicyHandler(svc *service.PolicyService, lg *zap.Logger) *PolicyHandler {
	if lg == nil {
		lg = zap.NewNop()
	}
	return &PolicyHandler{svc: svc, logger: lg}
}

func (h *PolicyHandler) Upsert(w http.ResponseWriter, r *http.Request) {
	var req dto.UpsertPolicyRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, r, http.StatusBadRequest, 1001, err.Error())
		return
	}

	item, err := h.svc.Upsert(r.Context(), req)
	if err != nil {
		status, code, msg := mapServiceError(err)
		writeError(w, r, status, code, msg)
		return
	}
	writeOK(w, r, toPolicyResponse(item))
}

func (h *PolicyHandler) Get(w http.ResponseWriter, r *http.Request) {
	tenantID := r.PathValue("tenantID")
	name := r.PathValue("name")

	item, err := h.svc.Get(r.Context(), tenantID, name)
	if err != nil {
		status, code, msg := mapServiceError(err)
		writeError(w, r, status, code, msg)
		return
	}
	writeOK(w, r, toPolicyResponse(item))
}

func (h *PolicyHandler) List(w http.ResponseWriter, r *http.Request) {
	enabled, err := parseOptionalBool(r.URL.Query().Get("enabled"))
	if err != nil {
		writeError(w, r, http.StatusBadRequest, 1001, err.Error())
		return
	}
	tenantID := r.URL.Query().Get("tenant_id")
	limit, offset := parsePaging(r, 50, 0)

	items, page, err := h.svc.List(r.Context(), tenantID, enabled, limit, offset)
	if err != nil {
		status, code, msg := mapServiceError(err)
		writeError(w, r, status, code, msg)
		return
	}

	out := make([]dto.PolicyResponse, 0, len(items))
	for i := range items {
		out = append(out, toPolicyResponse(&items[i]))
	}
	writeList(w, r, out, page.Limit, page.Offset)
}

func (h *PolicyHandler) Delete(w http.ResponseWriter, r *http.Request) {
	tenantID := r.PathValue("tenantID")
	name := r.PathValue("name")

	if err := h.svc.Delete(r.Context(), tenantID, name); err != nil {
		status, code, msg := mapServiceError(err)
		writeError(w, r, status, code, msg)
		return
	}
	writeOK(w, r, dto.Empty{})
}

func toPolicyResponse(m *model.GPUPolicy) dto.PolicyResponse {
	return dto.PolicyResponse{
		ID:                 m.ID,
		TenantID:           m.TenantID,
		Name:               m.Name,
		Queue:              m.Queue,
		Priority:           m.Priority,
		Enabled:            m.Enabled,
		Preemptible:        m.Preemptible,
		RequireHealthy:     m.RequireHealthy,
		RequireMIG:         m.RequireMIG,
		MaxGPUCount:        m.MaxGPUCount,
		RequiredGPUModel:   m.RequiredGPUModel,
		RequiredNodeLabels: jsonToAnyMapBytes(m.RequiredNodeLabelsJSON),
		Selector:           jsonToAnyMapBytes(m.SelectorJSON),
		Description:        m.Description,
		CreatedAt:          m.CreatedAt.UTC().Format(time.RFC3339Nano),
		UpdatedAt:          m.UpdatedAt.UTC().Format(time.RFC3339Nano),
	}
}

func jsonToAnyMapBytes(v []byte) map[string]any {
	if len(v) == 0 {
		return nil
	}
	var out map[string]any
	if err := json.Unmarshal(v, &out); err != nil {
		return nil
	}
	return out
}
