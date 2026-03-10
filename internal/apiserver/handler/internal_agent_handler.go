package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"gpu-scheduler-platform/internal/apiserver/dto"
	"gpu-scheduler-platform/internal/apiserver/service"

	"go.uber.org/zap"
)

type InternalAgentHandler struct {
	svc *service.AgentIngestService
	lg  *zap.Logger
}

func NewInternalAgentHandler(svc *service.AgentIngestService, lg *zap.Logger) *InternalAgentHandler {
	return &InternalAgentHandler{
		svc: svc,
		lg:  lg,
	}
}

func (h *InternalAgentHandler) Report(w http.ResponseWriter, r *http.Request) {
	var raw map[string]any
	if err := json.NewDecoder(r.Body).Decode(&raw); err != nil {
		writeErr(w, r, http.StatusBadRequest, 400, "invalid json body")
		return
	}

	if typ, _ := raw["type"].(string); strings.EqualFold(strings.TrimSpace(typ), "heartbeat") {
		nodeName, _ := raw["node_name"].(string)
		writeOK(w, r, http.StatusOK, map[string]any{
			"accepted":  true,
			"type":      "heartbeat",
			"node_name": nodeName,
		})
		return
	}

	body, err := json.Marshal(raw)
	if err != nil {
		writeErr(w, r, http.StatusBadRequest, 400, "invalid body")
		return
	}

	var req dto.AgentReportRequest
	if err := json.Unmarshal(body, &req); err != nil {
		writeErr(w, r, http.StatusBadRequest, 400, "invalid report payload")
		return
	}

	ts := time.Now().UTC()
	if strings.TrimSpace(req.Timestamp) != "" {
		if parsed, err := time.Parse(time.RFC3339Nano, req.Timestamp); err == nil {
			ts = parsed.UTC()
		}
	}

	gpus := make([]service.IngestGPUInfo, 0, len(req.GPUs))
	for _, g := range req.GPUs {
		gpus = append(gpus, service.IngestGPUInfo{
			ID:            g.ID,
			UUID:          g.UUID,
			Index:         g.Index,
			Model:         g.Model,
			MemoryMiB:     g.MemoryMiB,
			FreeMemoryMiB: g.FreeMemoryMiB,
			Healthy:       g.Healthy,
			Health:        g.Health,
		})
	}

	migs := make([]service.IngestMIGInfo, 0, len(req.MIGs))
	for _, m := range req.MIGs {
		migs = append(migs, service.IngestMIGInfo{
			ID:         m.ID,
			ParentUUID: m.ParentUUID,
			Profile:    m.Profile,
			MemoryMiB:  m.MemoryMiB,
		})
	}

	links := make([]service.IngestGPULink, 0, len(req.Topology))
	for _, l := range req.Topology {
		links = append(links, service.IngestGPULink{
			From: l.From,
			To:   l.To,
			Type: l.Type,
		})
	}

	pods := make([]service.IngestPodGPUInfo, 0, len(req.PodBindings))
	for _, p := range req.PodBindings {
		pods = append(pods, service.IngestPodGPUInfo{
			PodName:   p.PodName,
			Namespace: p.Namespace,
			GPUIDs:    append([]string(nil), p.GPUIDs...),
		})
	}

	out, err := h.svc.IngestReport(r.Context(), service.IngestAgentReportRequest{
		NodeName:    req.NodeName,
		Timestamp:   ts,
		GPUs:        gpus,
		MIGs:        migs,
		Topology:    links,
		PodBindings: pods,
	})
	if err != nil {
		handleServiceError(w, r, err)
		return
	}

	writeOK(w, r, http.StatusOK, dto.AgentReportResponse{
		Accepted:        out.Accepted,
		SnapshotVersion: out.SnapshotVersion,
		NodeName:        out.NodeName,
		GPUCount:        out.GPUCount,
	})
}
