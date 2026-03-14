package handler

import (
	"net/http"
	"time"

	"gpu-scheduler-platform/internal/apiserver/dto"
	"gpu-scheduler-platform/internal/apiserver/service"
	model "gpu-scheduler-platform/internal/repo/models"

	"go.uber.org/zap"
)

type TenantHandler struct {
	svc    *service.TenantService
	logger *zap.Logger
}

func NewTenantHandler(svc *service.TenantService, lg *zap.Logger) *TenantHandler {
	if lg == nil {
		lg = zap.NewNop()
	}
	return &TenantHandler{svc: svc, logger: lg}
}

func (h *TenantHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateTenantRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, r, http.StatusBadRequest, 400, err.Error())
		return
	}

	item, err := h.svc.Create(r.Context(), req)
	if err != nil {
		status, code, msg := mapServiceError(err)
		writeError(w, r, status, code, msg)
		return
	}

	writeCreated(w, r, toTenantResponse(item))
}

func (h *TenantHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	item, err := h.svc.Get(r.Context(), id)
	if err != nil {
		status, code, msg := mapServiceError(err)
		writeError(w, r, status, code, msg)
		return
	}
	writeOK(w, r, toTenantResponse(item))
}

func (h *TenantHandler) List(w http.ResponseWriter, r *http.Request) {
	enabled, err := parseOptionalBool(r.URL.Query().Get("enabled"))
	if err != nil {
		writeError(w, r, http.StatusBadRequest, 400, err.Error())
		return
	}
	limit, offset := parsePaging(r, 50, 0)

	items, page, err := h.svc.List(r.Context(), enabled, limit, offset)
	if err != nil {
		status, code, msg := mapServiceError(err)
		writeError(w, r, status, code, msg)
		return
	}

	out := make([]dto.TenantResponse, 0, len(items))
	for i := range items {
		out = append(out, toTenantResponse(&items[i]))
	}
	writeList(w, r, out, page.Limit, page.Offset)
}

func (h *TenantHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var req dto.UpdateTenantRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, r, http.StatusBadRequest, 400, err.Error())
		return
	}

	item, err := h.svc.Update(r.Context(), id, req)
	if err != nil {
		status, code, msg := mapServiceError(err)
		writeError(w, r, status, code, msg)
		return
	}
	writeOK(w, r, toTenantResponse(item))
}

func (h *TenantHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.svc.Delete(r.Context(), id); err != nil {
		status, code, msg := mapServiceError(err)
		writeError(w, r, status, code, msg)
		return
	}
	writeOK(w, r, dto.Empty{})
}

func toTenantResponse(m *model.Tenant) dto.TenantResponse {
	return dto.TenantResponse{
		ID:          m.ID,
		Name:        m.Name,
		Enabled:     m.Enabled,
		Description: m.Description,
		CreatedAt:   m.CreatedAt.UTC().Format(time.RFC3339Nano),
		UpdatedAt:   m.UpdatedAt.UTC().Format(time.RFC3339Nano),
	}
}
