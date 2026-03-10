package validating

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	webhookserver "gpu-scheduler-platform/internal/webhook/server"

	"go.uber.org/zap"
)

type GPUQuotaValidateHandler struct {
	logger *zap.Logger
}

func NewGPUQuotaValidateHandler(lg *zap.Logger) *GPUQuotaValidateHandler {
	return &GPUQuotaValidateHandler{logger: lg}
}

func (h *GPUQuotaValidateHandler) Handle(w http.ResponseWriter, r *http.Request) {
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

	if err := validateGPUQuota(obj); err != nil {
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

func validateGPUQuota(obj map[string]any) error {
	tenantID, _ := obj["tenant_id"].(string)
	if strings.TrimSpace(tenantID) == "" {
		return fmt.Errorf("tenant_id is required")
	}

	if v, ok := asInt(obj["max_gpu_count"]); ok && v < 0 {
		return fmt.Errorf("max_gpu_count must be >= 0")
	}
	if v, ok := asInt(obj["max_running_jobs"]); ok && v < 0 {
		return fmt.Errorf("max_running_jobs must be >= 0")
	}
	if v, ok := asInt(obj["max_queued_jobs"]); ok && v < 0 {
		return fmt.Errorf("max_queued_jobs must be >= 0")
	}
	if v, ok := asInt64(obj["max_gpu_memory_mib"]); ok && v < 0 {
		return fmt.Errorf("max_gpu_memory_mib must be >= 0")
	}

	return nil
}
