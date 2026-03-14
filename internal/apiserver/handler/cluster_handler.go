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

type ClusterHandler struct {
	svc    *service.ClusterService
	logger *zap.Logger
}

func NewClusterHandler(svc *service.ClusterService, lg *zap.Logger) *ClusterHandler {
	if lg == nil {
		lg = zap.NewNop()
	}
	return &ClusterHandler{svc: svc, logger: lg}
}

func (h *ClusterHandler) ListNodes(w http.ResponseWriter, r *http.Request) {
	schedulable, err := parseOptionalBool(r.URL.Query().Get("schedulable"))
	if err != nil {
		writeError(w, r, http.StatusBadRequest, 400, err.Error())
		return
	}
	limit, offset := parsePaging(r, 100, 0)

	items, page, err := h.svc.ListNodes(r.Context(), service.ListNodesInput{
		ClusterName: r.URL.Query().Get("cluster_name"),
		State:       r.URL.Query().Get("state"),
		Schedulable: schedulable,
		Limit:       limit,
		Offset:      offset,
	})
	if err != nil {
		status, code, msg := mapServiceError(err)
		writeError(w, r, status, code, msg)
		return
	}

	out := make([]dto.ClusterNodeSummaryResponse, 0, len(items))
	for i := range items {
		out = append(out, toClusterNodeSummary(&items[i]))
	}
	writeList(w, r, out, page.Limit, page.Offset)
}

func (h *ClusterHandler) GetNode(w http.ResponseWriter, r *http.Request) {
	nodeName := r.PathValue("nodeName")

	item, err := h.svc.GetNode(r.Context(), nodeName)
	if err != nil {
		status, code, msg := mapServiceError(err)
		writeError(w, r, status, code, msg)
		return
	}

	writeOK(w, r, toClusterNodeDetail(item))
}

func toClusterNodeSummary(m *model.Node) dto.ClusterNodeSummaryResponse {
	return dto.ClusterNodeSummaryResponse{
		NodeName:          m.NodeName,
		ClusterName:       m.ClusterName,
		Source:            m.Source,
		State:             m.State,
		Schedulable:       m.Schedulable,
		GPUCount:          m.GPUCount,
		HealthyGPUCount:   m.HealthyGPUCount,
		TotalMemoryMiB:    m.TotalMemoryMiB,
		FreeMemoryMiB:     m.FreeMemoryMiB,
		Labels:            jsonToStringMap(m.LabelsJSON),
		Annotations:       jsonToStringMap(m.AnnotationsJSON),
		LastReportTime:    formatPtrTime(m.LastReportTime),
		LastHeartbeatTime: formatPtrTime(m.LastHeartbeatTime),
		CreatedAt:         m.CreatedAt.UTC().Format(time.RFC3339Nano),
		UpdatedAt:         m.UpdatedAt.UTC().Format(time.RFC3339Nano),
	}
}

func toClusterNodeDetail(d *service.NodeDetail) dto.ClusterNodeDetailResponse {
	resp := dto.ClusterNodeDetailResponse{
		Node: toClusterNodeSummary(d.Node),
	}

	if d.Heartbeat != nil {
		resp.Heartbeat = &dto.ClusterHeartbeatResponse{
			Status:     d.Heartbeat.Status,
			Message:    d.Heartbeat.Message,
			LastSeenAt: d.Heartbeat.LastSeenAt.UTC().Format(time.RFC3339Nano),
			UpdatedAt:  d.Heartbeat.UpdatedAt.UTC().Format(time.RFC3339Nano),
		}
	}

	if d.LatestSnapshot != nil {
		resp.LatestSnapshot = &dto.ClusterSnapshotResponse{
			ID:              d.LatestSnapshot.ID,
			Version:         d.LatestSnapshot.Version,
			AgentVersion:    d.LatestSnapshot.AgentVersion,
			ClusterName:     d.LatestSnapshot.ClusterName,
			NodeName:        d.LatestSnapshot.NodeName,
			Source:          d.LatestSnapshot.Source,
			NodeState:       d.LatestSnapshot.NodeState,
			Schedulable:     d.LatestSnapshot.Schedulable,
			GPUCount:        d.LatestSnapshot.GPUCount,
			HealthyGPUCount: d.LatestSnapshot.HealthyGPUCount,
			TotalMemoryMiB:  d.LatestSnapshot.TotalMemoryMiB,
			FreeMemoryMiB:   d.LatestSnapshot.FreeMemoryMiB,
			Labels:          jsonToStringMap(d.LatestSnapshot.LabelsJSON),
			Annotations:     jsonToStringMap(d.LatestSnapshot.AnnotationsJSON),
			Topology:        jsonToAnyMap(d.LatestSnapshot.TopologyJSON),
			ReportTime:      d.LatestSnapshot.ReportTime.UTC().Format(time.RFC3339Nano),
			CreatedAt:       d.LatestSnapshot.CreatedAt.UTC().Format(time.RFC3339Nano),
		}
	}

	for i := range d.GPUDevices {
		item := d.GPUDevices[i]
		resp.GPUDevices = append(resp.GPUDevices, dto.ClusterGPUDeviceResponse{
			ID:                item.ID,
			SnapshotID:        item.SnapshotID,
			NodeName:          item.NodeName,
			UUID:              item.UUID,
			GPUIndex:          item.GPUIndex,
			Model:             item.Model,
			Vendor:            item.Vendor,
			Type:              item.Type,
			MemoryMiB:         item.MemoryMiB,
			FreeMemoryMiB:     item.FreeMemoryMiB,
			Healthy:           item.Healthy,
			Health:            item.Health,
			MIGEnabled:        item.MIGEnabled,
			MIGProfile:        item.MIGProfile,
			UtilizationGPU:    item.UtilizationGPU,
			UtilizationMemory: item.UtilizationMemory,
			Temperature:       item.Temperature,
			PowerWatts:        item.PowerWatts,
			Labels:            jsonToStringMap(item.LabelsJSON),
			Annotations:       jsonToStringMap(item.AnnotationsJSON),
			Allocated:         item.Allocated,
			Reserved:          item.Reserved,
			CreatedAt:         item.CreatedAt.UTC().Format(time.RFC3339Nano),
		})
	}

	for i := range d.MIGDevices {
		item := d.MIGDevices[i]
		resp.MIGDevices = append(resp.MIGDevices, dto.ClusterMIGDeviceResponse{
			ID:            item.ID,
			SnapshotID:    item.SnapshotID,
			NodeName:      item.NodeName,
			ParentGPUUUID: item.ParentGPUUUID,
			MIGUUID:       item.MIGUUID,
			Profile:       item.Profile,
			MemoryMiB:     item.MemoryMiB,
			Healthy:       item.Healthy,
			Allocated:     item.Allocated,
			Reserved:      item.Reserved,
			CreatedAt:     item.CreatedAt.UTC().Format(time.RFC3339Nano),
		})
	}

	for i := range d.RuntimeBindings {
		item := d.RuntimeBindings[i]
		resp.RuntimeBindings = append(resp.RuntimeBindings, dto.ClusterRuntimeBindingResponse{
			ID:         item.ID,
			SnapshotID: item.SnapshotID,
			NodeName:   item.NodeName,
			Namespace:  item.Namespace,
			PodName:    item.PodName,
			GPUIDs:     jsonToStringSlice(item.GPUIDsJSON),
			CreatedAt:  item.CreatedAt.UTC().Format(time.RFC3339Nano),
		})
	}

	return resp
}

func jsonToStringMap(v []byte) map[string]string {
	if len(v) == 0 {
		return nil
	}
	var out map[string]string
	if err := json.Unmarshal(v, &out); err != nil {
		return nil
	}
	return out
}

func jsonToAnyMap(v []byte) map[string]any {
	if len(v) == 0 {
		return nil
	}
	var out map[string]any
	if err := json.Unmarshal(v, &out); err != nil {
		return nil
	}
	return out
}

func jsonToStringSlice(v []byte) []string {
	if len(v) == 0 {
		return nil
	}
	var out []string
	if err := json.Unmarshal(v, &out); err != nil {
		return nil
	}
	return out
}
