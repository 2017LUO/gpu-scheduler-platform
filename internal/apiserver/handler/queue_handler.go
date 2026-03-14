package handler

import (
	"net/http"
	"time"

	"gpu-scheduler-platform/internal/apiserver/dto"
	"gpu-scheduler-platform/internal/apiserver/service"
	model "gpu-scheduler-platform/internal/repo/models"

	"go.uber.org/zap"
)

type QueueHandler struct {
	svc    *service.QueueService
	logger *zap.Logger
}

func NewQueueHandler(svc *service.QueueService, lg *zap.Logger) *QueueHandler {
	if lg == nil {
		lg = zap.NewNop()
	}
	return &QueueHandler{svc: svc, logger: lg}
}

func (h *QueueHandler) Upsert(w http.ResponseWriter, r *http.Request) {
	var req dto.UpsertQueueRequest
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
	writeOK(w, r, toQueueResponse(item))
}

func (h *QueueHandler) Get(w http.ResponseWriter, r *http.Request) {
	tenantID := r.PathValue("tenantID")
	name := r.PathValue("name")

	item, err := h.svc.Get(r.Context(), tenantID, name)
	if err != nil {
		status, code, msg := mapServiceError(err)
		writeError(w, r, status, code, msg)
		return
	}
	writeOK(w, r, toQueueResponse(item))
}

func (h *QueueHandler) List(w http.ResponseWriter, r *http.Request) {
	enabled, err := parseOptionalBool(r.URL.Query().Get("enabled"))
	if err != nil {
		writeError(w, r, http.StatusBadRequest, 400, err.Error())
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

	out := make([]dto.QueueResponse, 0, len(items))
	for i := range items {
		out = append(out, toQueueResponse(&items[i]))
	}
	writeList(w, r, out, page.Limit, page.Offset)
}

func (h *QueueHandler) Delete(w http.ResponseWriter, r *http.Request) {
	tenantID := r.PathValue("tenantID")
	name := r.PathValue("name")

	if err := h.svc.Delete(r.Context(), tenantID, name); err != nil {
		status, code, msg := mapServiceError(err)
		writeError(w, r, status, code, msg)
		return
	}
	writeOK(w, r, dto.Empty{})
}

func toQueueResponse(m *model.Queue) dto.QueueResponse {
	return dto.QueueResponse{
		ID:          m.ID,
		TenantID:    m.TenantID,
		Name:        m.Name,
		Weight:      m.Weight,
		Priority:    m.Priority,
		Enabled:     m.Enabled,
		Description: m.Description,
		CreatedAt:   m.CreatedAt.UTC().Format(time.RFC3339Nano),
		UpdatedAt:   m.UpdatedAt.UTC().Format(time.RFC3339Nano),
	}
}
