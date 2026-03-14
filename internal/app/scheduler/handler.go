package scheduler

import (
	"encoding/json"
	"net/http"
	"time"

	"gpu-scheduler-platform/pkg/version"
)

type healthResponse struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Data    healthDataPayload `json:"data"`
}

type healthDataPayload struct {
	Status  string            `json:"status"`
	Time    string            `json:"time"`
	Version version.Info      `json:"version"`
	Details map[string]string `json:"details,omitempty"`
}

func (a *App) handleRoot(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, healthResponse{
		Code:    0,
		Message: "ok",
		Data: healthDataPayload{
			Status:  "ok",
			Time:    time.Now().UTC().Format(time.RFC3339Nano),
			Version: version.Get(),
			Details: map[string]string{
				"service": a.Config.Service.Name,
				"mode":    "scheduler",
			},
		},
	})
}

func (a *App) handleHealthz(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, healthResponse{
		Code:    0,
		Message: "ok",
		Data: healthDataPayload{
			Status:  "ok",
			Time:    time.Now().UTC().Format(time.RFC3339Nano),
			Version: version.Get(),
		},
	})
}

func (a *App) handleReadyz(w http.ResponseWriter, r *http.Request) {
	ok, details := a.Readiness(r.Context())

	statusCode := http.StatusOK
	status := "ready"
	message := "ok"
	if !ok {
		statusCode = http.StatusServiceUnavailable
		status = "not_ready"
		message = "dependency not ready"
	}

	writeJSON(w, statusCode, healthResponse{
		Code:    0,
		Message: message,
		Data: healthDataPayload{
			Status:  status,
			Time:    time.Now().UTC().Format(time.RFC3339Nano),
			Version: version.Get(),
			Details: details,
		},
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
