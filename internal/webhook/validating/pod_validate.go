package validating

import (
	"encoding/json"
	"net/http"

	webhookserver "gpu-scheduler-platform/internal/webhook/server"

	"go.uber.org/zap"
)

type PodValidateHandler struct {
	logger *zap.Logger
}

func NewPodValidateHandler(lg *zap.Logger) *PodValidateHandler {
	return &PodValidateHandler{logger: lg}
}

func (h *PodValidateHandler) Handle(w http.ResponseWriter, r *http.Request) {
	var req webhookserver.AdmissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeAdmission(w, http.StatusBadRequest, webhookserver.AdmissionResponse{
			Allowed: false,
			Message: "invalid json body",
		})
		return
	}

	// 当前先做占位放行。
	writeAdmission(w, http.StatusOK, webhookserver.AdmissionResponse{
		Allowed: true,
		Message: "validated",
	})
}
