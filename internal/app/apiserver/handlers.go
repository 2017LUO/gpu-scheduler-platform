package apiserver

import (
	"encoding/json"
	"net/http"
	"time"

	"gpu-scheduler-platform/internal/middleware"
	"gpu-scheduler-platform/pkg/version"
)

type healthResponse struct {
	Code      int               `json:"code"`
	Message   string            `json:"message"`
	RequestID string            `json:"request_id,omitempty"`
	Data      healthDataPayload `json:"data"`
}

type healthDataPayload struct {
	Status  string            `json:"status"`
	Time    string            `json:"time"`
	Version version.Info      `json:"version"`
	Details map[string]string `json:"details,omitempty"`
}

func (a *App) handleRoot(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, healthResponse{
		Code:      0,
		Message:   "ok",
		RequestID: middleware.RequestIDFromContext(r.Context()),
		Data: healthDataPayload{
			Status:  "ok",
			Time:    time.Now().UTC().Format(time.RFC3339Nano),
			Version: version.Get(),
			Details: map[string]string{
				"service": a.Config.Service.Name,
			},
		},
	})
}

func (a *App) handleHealthz(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, healthResponse{
		Code:      0,
		Message:   "ok",
		RequestID: middleware.RequestIDFromContext(r.Context()),
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
		Code:      0,
		Message:   message,
		RequestID: middleware.RequestIDFromContext(r.Context()),
		Data: healthDataPayload{
			Status:  status,
			Time:    time.Now().UTC().Format(time.RFC3339Nano),
			Version: version.Get(),
			Details: details,
		},
	})
}

func (a *App) handleNotImplementedAPI(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusNotImplemented, map[string]any{
		"code":       501,
		"message":    "api not implemented yet",
		"request_id": middleware.RequestIDFromContext(r.Context()),
		"path":       r.URL.Path,
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
