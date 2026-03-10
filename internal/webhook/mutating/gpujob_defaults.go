package mutating

import (
	"encoding/json"
	"net/http"

	webhookserver "gpu-scheduler-platform/internal/webhook/server"

	"go.uber.org/zap"
)

type GPUJobDefaultsHandler struct {
	logger *zap.Logger
}

func NewGPUJobDefaultsHandler(lg *zap.Logger) *GPUJobDefaultsHandler {
	return &GPUJobDefaultsHandler{logger: lg}
}

func (h *GPUJobDefaultsHandler) Handle(w http.ResponseWriter, r *http.Request) {
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

	if _, ok := obj["queue"]; !ok {
		obj["queue"] = "default"
	}
	if _, ok := obj["priority"]; !ok {
		obj["priority"] = "normal"
	}

	reqObj, ok := obj["requirement"].(map[string]any)
	if !ok || reqObj == nil {
		reqObj = map[string]any{}
	}
	if _, ok := reqObj["require_healthy"]; !ok {
		reqObj["require_healthy"] = true
	}
	if _, ok := reqObj["retryable"]; !ok {
		reqObj["retryable"] = true
	}
	if _, ok := reqObj["max_retry"]; !ok {
		reqObj["max_retry"] = float64(0)
	}
	obj["requirement"] = reqObj

	writeAdmission(w, http.StatusOK, webhookserver.AdmissionResponse{
		Allowed: true,
		Message: "defaults applied",
		Patch:   obj,
	})
}

func writeAdmission(w http.ResponseWriter, status int, resp webhookserver.AdmissionResponse) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(resp)
}
