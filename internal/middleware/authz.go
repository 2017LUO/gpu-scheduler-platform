package middleware

import (
	"net/http"

	"gpu-scheduler-platform/internal/apiserver/errcode"
	authpkg "gpu-scheduler-platform/internal/auth"
)

type AuthZConfig struct {
	Enabled          bool
	Authorizer       *authpkg.Authorizer
	PublicExactPaths []string
	PublicPrefixes   []string
}

func AuthZ(cfg AuthZConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		if next == nil {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
		}
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if shouldSkipPath(r.URL.Path, cfg.PublicExactPaths, cfg.PublicPrefixes) {
				next.ServeHTTP(w, r)
				return
			}
			if !cfg.Enabled {
				next.ServeHTTP(w, r)
				return
			}

			perm := authpkg.PermissionForRequest(r.Method, r.URL.Path)
			if perm == authpkg.PermissionNone {
				next.ServeHTTP(w, r)
				return
			}

			sub, ok := authpkg.SubjectFromContext(r.Context())
			if !ok || sub == nil {
				writeAuthError(w, r, http.StatusUnauthorized, errcode.CodeUnauthenticated, "subject is not authenticated")
				return
			}

			if cfg.Authorizer == nil || !cfg.Authorizer.Allowed(sub, perm) {
				writeAuthError(w, r, http.StatusForbidden, errcode.CodePermissionDenied, "permission denied")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
