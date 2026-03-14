package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"gpu-scheduler-platform/internal/apiserver/dto"
	"gpu-scheduler-platform/internal/apiserver/errcode"
	"gpu-scheduler-platform/internal/apiserver/service"
	"gpu-scheduler-platform/internal/middleware"
	repoimpl "gpu-scheduler-platform/internal/repo/mysql"

	"go.uber.org/zap"
)

type Handlers struct {
	Job           *JobHandler
	InternalAgent *InternalAgentHandler
	Tenant        *TenantHandler
	Queue         *QueueHandler
	Quota         *QuotaHandler
	Policy        *PolicyHandler
	Cluster       *ClusterHandler
}

func NewHandlers(svcs *service.Services, lg *zap.Logger) *Handlers {
	if lg == nil {
		lg = zap.NewNop()
	}
	if svcs == nil {
		return &Handlers{}
	}
	return &Handlers{
		Job:           NewJobHandler(svcs.Job, lg),
		InternalAgent: NewInternalAgentHandler(svcs.InternalAgent, lg),
		Tenant:        NewTenantHandler(svcs.Tenant, lg),
		Queue:         NewQueueHandler(svcs.Queue, lg),
		Quota:         NewQuotaHandler(svcs.Quota, lg),
		Policy:        NewPolicyHandler(svcs.Policy, lg),
		Cluster:       NewClusterHandler(svcs.Cluster, lg),
	}
}

func writeOK[T any](w http.ResponseWriter, r *http.Request, data T) {
	writeJSON(w, http.StatusOK, dto.Response[T]{
		Code:      0,
		Message:   "ok",
		RequestID: middleware.RequestIDFromContext(r.Context()),
		Data:      data,
	})
}

func writeCreated[T any](w http.ResponseWriter, r *http.Request, data T) {
	writeJSON(w, http.StatusCreated, dto.Response[T]{
		Code:      0,
		Message:   "ok",
		RequestID: middleware.RequestIDFromContext(r.Context()),
		Data:      data,
	})
}

func writeList[T any](w http.ResponseWriter, r *http.Request, items []T, limit, offset int) {
	writeJSON(w, http.StatusOK, dto.Response[dto.ListData[T]]{
		Code:      0,
		Message:   "ok",
		RequestID: middleware.RequestIDFromContext(r.Context()),
		Data: dto.ListData[T]{
			Items: items,
			Meta: dto.ListMeta{
				Limit:  limit,
				Offset: offset,
				Count:  len(items),
			},
		},
	})
}

func writeError(w http.ResponseWriter, r *http.Request, status int, code int, message string) {
	writeJSON(w, status, map[string]any{
		"code":       code,
		"error_code": errcode.String(code),
		"message":    message,
		"request_id": middleware.RequestIDFromContext(r.Context()),
		"path":       r.URL.Path,
	})
}

func decodeJSON(r *http.Request, v any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func parsePaging(r *http.Request, defaultLimit, defaultOffset int) (int, int) {
	limit := defaultLimit
	offset := defaultOffset

	if s := r.URL.Query().Get("limit"); s != "" {
		if n, err := strconv.Atoi(s); err == nil {
			limit = n
		}
	}
	if s := r.URL.Query().Get("offset"); s != "" {
		if n, err := strconv.Atoi(s); err == nil {
			offset = n
		}
	}
	return limit, offset
}

func parseOptionalBool(s string) (*bool, error) {
	if s == "" {
		return nil, nil
	}
	v, err := strconv.ParseBool(s)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func mapServiceError(err error) (int, int, string) {
	switch {
	case err == nil:
		return http.StatusOK, errcode.CodeOK, "ok"
	case errors.Is(err, repoimpl.ErrInvalidArgument):
		return http.StatusBadRequest, errcode.CodeInvalidArgument, err.Error()
	case errors.Is(err, repoimpl.ErrNotFound):
		return http.StatusNotFound, errcode.CodeNotFound, err.Error()
	case errors.Is(err, repoimpl.ErrConflict):
		return http.StatusConflict, errcode.CodeConflict, err.Error()
	case errors.Is(err, service.ErrQuotaExceeded):
		return http.StatusConflict, errcode.CodeQuotaExceeded, err.Error()
	default:
		return http.StatusInternalServerError, errcode.CodeInternal, err.Error()
	}
}
