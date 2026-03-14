package handler

import (
	"net/http"

	"gpu-scheduler-platform/internal/apiserver/dto"
	"gpu-scheduler-platform/internal/apiserver/service"

	"go.uber.org/zap"
)

type InternalAgentHandler struct {
	svc    *service.InternalAgentService
	logger *zap.Logger
}

func NewInternalAgentHandler(svc *service.InternalAgentService, lg *zap.Logger) *InternalAgentHandler {
	if lg == nil {
		lg = zap.NewNop()
	}
	return &InternalAgentHandler{svc: svc, logger: lg}
}

func (h *InternalAgentHandler) Heartbeat(w http.ResponseWriter, r *http.Request) {
	var req dto.AgentHeartbeatRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, r, http.StatusBadRequest, 400, err.Error())
		return
	}

	seenAt, err := h.svc.HandleHeartbeat(r.Context(), req)
	if err != nil {
		status, code, msg := mapServiceError(err)
		writeError(w, r, status, code, msg)
		return
	}

	writeOK(w, r, dto.AgentHeartbeatResponse{
		NodeName: req.NodeName,
		Status:   req.Status,
		SeenAt:   seenAt.UTC().Format("2006-01-02T15:04:05.999999999Z07:00"),
	})
}

func (h *InternalAgentHandler) Report(w http.ResponseWriter, r *http.Request) {
	var req dto.AgentReportRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, r, http.StatusBadRequest, 400, err.Error())
		return
	}

	snapshotID, reportTime, err := h.svc.HandleReport(r.Context(), req)
	if err != nil {
		status, code, msg := mapServiceError(err)
		writeError(w, r, status, code, msg)
		return
	}

	writeOK(w, r, dto.AgentReportResponse{
		NodeName:   req.NodeName,
		SnapshotID: snapshotID,
		ReportTime: reportTime.UTC().Format("2006-01-02T15:04:05.999999999Z07:00"),
	})
}
