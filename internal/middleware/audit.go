package middleware

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	auditpkg "gpu-scheduler-platform/internal/audit"
	authpkg "gpu-scheduler-platform/internal/auth"
)

type AuditConfig struct {
	Enabled         bool
	Recorder        *auditpkg.Recorder
	SkipExactPaths  []string
	SkipPrefixes    []string
	OnlyMutatingAPI bool
}

type auditResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *auditResponseWriter) WriteHeader(statusCode int) {
	w.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func Audit(cfg AuditConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		if next == nil {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
		}
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !cfg.Enabled || cfg.Recorder == nil {
				next.ServeHTTP(w, r)
				return
			}
			if shouldSkipPath(r.URL.Path, cfg.SkipExactPaths, cfg.SkipPrefixes) {
				next.ServeHTTP(w, r)
				return
			}
			if cfg.OnlyMutatingAPI {
				if !(r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch || r.Method == http.MethodDelete) {
					next.ServeHTTP(w, r)
					return
				}
				if !strings.HasPrefix(r.URL.Path, "/api/v1/") {
					next.ServeHTTP(w, r)
					return
				}
			}

			start := time.Now()
			aw := &auditResponseWriter{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(aw, r)

			sub, _ := authpkg.SubjectFromContext(r.Context())
			actor := "anonymous"
			tenantID := ""
			if sub != nil {
				if sub.Name != "" {
					actor = sub.Name
				} else if sub.SubjectID != "" {
					actor = sub.SubjectID
				}
				tenantID = sub.TenantID
			}

			resourceID := firstNonEmpty(
				r.PathValue("id"),
				r.PathValue("tenantID"),
				r.PathValue("namespace"),
				r.PathValue("name"),
				r.PathValue("nodeName"),
				r.URL.Query().Get("tenant_id"),
				r.URL.Path,
			)
			resourceName := firstNonEmpty(
				r.PathValue("name"),
				r.PathValue("nodeName"),
				r.PathValue("namespace"),
				resourceID,
			)

			status := "SUCCESS"
			if aw.status >= 400 {
				status = "FAILED"
			}

			cfg.Recorder.MustRecord(r.Context(), auditpkg.Record{
				TenantID:     tenantID,
				Actor:        actor,
				Action:       authpkg.ActionForRequest(r.Method, r.URL.Path),
				ResourceType: authpkg.ResourceTypeForPath(r.URL.Path),
				ResourceID:   resourceID,
				ResourceName: resourceName,
				Status:       status,
				RequestID:    RequestIDFromContext(r.Context()),
				Detail: map[string]any{
					"method":      r.Method,
					"path":        r.URL.Path,
					"status_code": aw.status,
					"duration_ms": strconv.FormatInt(time.Since(start).Milliseconds(), 10),
					"remote_addr": r.RemoteAddr,
				},
			})
		})
	}
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}
