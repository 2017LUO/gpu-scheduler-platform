package validating

import (
	"encoding/json"
	"fmt"
	"net/http"

	webhookserver "gpu-scheduler-platform/internal/webhook/server"

	"go.uber.org/zap"
)

type GPUPolicyValidateHandler struct {
	logger *zap.Logger
}

func NewGPUPolicyValidateHandler(lg *zap.Logger) *GPUPolicyValidateHandler {
	return &GPUPolicyValidateHandler{logger: lg}
}

func (h *GPUPolicyValidateHandler) Handle(w http.ResponseWriter, r *http.Request) {
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

	if err := validateGPUPolicy(obj); err != nil {
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

func validateGPUPolicy(obj map[string]any) error {
	if weights, ok := obj["scoring_weights"].(map[string]any); ok && weights != nil {
		for k, v := range weights {
			x, ok := asInt(v)
			if !ok || x < 0 {
				return fmt.Errorf("invalid scoring weight for %s", k)
			}
		}
	}
	return nil
}
