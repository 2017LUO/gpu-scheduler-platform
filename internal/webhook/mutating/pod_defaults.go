package mutating

import (
	"encoding/json"
	"net/http"

	webhookserver "gpu-scheduler-platform/internal/webhook/server"

	"go.uber.org/zap"
)

type PodDefaultsHandler struct {
	logger *zap.Logger
}

func NewPodDefaultsHandler(lg *zap.Logger) *PodDefaultsHandler {
	return &PodDefaultsHandler{logger: lg}
}

func (h *PodDefaultsHandler) Handle(w http.ResponseWriter, r *http.Request) {
	var req webhookserver.AdmissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeAdmission(w, http.StatusBadRequest, webhookserver.AdmissionResponse{
			Allowed: false,
			Message: "invalid json body",
		})
		return
	}

	obj := req.Object
	if obj == nil {
		obj = map[string]any{}
	}

	meta, _ := obj["metadata"].(map[string]any)
	if meta == nil {
		meta = map[string]any{}
	}
	annotations, _ := meta["annotations"].(map[string]any)
	if annotations == nil {
		annotations = map[string]any{}
	}
	if _, ok := annotations["gpu.scheduler.io/managed"]; !ok {
		annotations["gpu.scheduler.io/managed"] = "true"
	}
	meta["annotations"] = annotations
	obj["metadata"] = meta

	writeAdmission(w, http.StatusOK, webhookserver.AdmissionResponse{
		Allowed: true,
		Message: "defaults applied",
		Patch:   obj,
	})
}
