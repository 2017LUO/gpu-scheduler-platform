package validating

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	webhookserver "gpu-scheduler-platform/internal/webhook/server"

	"go.uber.org/zap"
)

type GPUJobValidateHandler struct {
	logger *zap.Logger
}

func NewGPUJobValidateHandler(lg *zap.Logger) *GPUJobValidateHandler {
	return &GPUJobValidateHandler{logger: lg}
}

func (h *GPUJobValidateHandler) Handle(w http.ResponseWriter, r *http.Request) {
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
		writeAdmission(w, http.StatusOK, webhookserver.AdmissionResponse{
			Allowed: false,
			Message: "object is required",
		})
		return
	}

	if err := validateGPUJob(obj); err != nil {
		writeAdmission(w, http.StatusOK, webhookserver.AdmissionResponse{
			Allowed: false,
			Message: err.Error(),
		})
		return
	}

	writeAdmission(w, http.StatusOK, webhookserver.AdmissionResponse{
		Allowed: true,
		Message: "validated",
	})
}

func validateGPUJob(obj map[string]any) error {
	tenantID, _ := obj["tenant_id"].(string)
	namespace, _ := obj["namespace"].(string)
	name, _ := obj["name"].(string)
	queue, _ := obj["queue"].(string)
	priority, _ := obj["priority"].(string)

	if strings.TrimSpace(tenantID) == "" {
		return fmt.Errorf("tenant_id is required")
	}
	if strings.TrimSpace(namespace) == "" {
		return fmt.Errorf("namespace is required")
	}
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("name is required")
	}

	switch strings.TrimSpace(queue) {
	case "", "default", "inference", "training", "batch":
	default:
		return fmt.Errorf("invalid queue")
	}

	switch strings.TrimSpace(priority) {
	case "", "critical", "high", "normal", "low":
	default:
		return fmt.Errorf("invalid priority")
	}

	reqObj, _ := obj["requirement"].(map[string]any)
	if reqObj == nil {
		return fmt.Errorf("requirement is required")
	}

	gpuCount, ok := asInt(reqObj["gpu_count"])
	if !ok || gpuCount <= 0 {
		return fmt.Errorf("requirement.gpu_count must be > 0")
	}
	memMiB, ok := asInt64(reqObj["gpu_memory_mib"])
	if !ok || memMiB < 0 {
		return fmt.Errorf("requirement.gpu_memory_mib must be >= 0")
	}

	maxRetry, ok := asInt(reqObj["max_retry"])
	if ok && maxRetry < 0 {
		return fmt.Errorf("requirement.max_retry must be >= 0")
	}

	return nil
}

func asInt(v any) (int, bool) {
	switch t := v.(type) {
	case int:
		return t, true
	case int64:
		return int(t), true
	case float64:
		return int(t), true
	default:
		return 0, false
	}
}

func asInt64(v any) (int64, bool) {
	switch t := v.(type) {
	case int:
		return int64(t), true
	case int64:
		return t, true
	case float64:
		return int64(t), true
	default:
		return 0, false
	}
}

func writeAdmission(w http.ResponseWriter, status int, resp webhookserver.AdmissionResponse) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(resp)
}
