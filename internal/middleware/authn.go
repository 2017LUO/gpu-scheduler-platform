package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"gpu-scheduler-platform/internal/apiserver/errcode"
	authpkg "gpu-scheduler-platform/internal/auth"
)

type AuthNConfig struct {
	Enabled          bool
	Manager          *authpkg.JWTManager
	PublicExactPaths []string
	PublicPrefixes   []string
}

func AuthN(cfg AuthNConfig) func(http.Handler) http.Handler {
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
				ctx := authpkg.WithSubject(r.Context(), authpkg.SystemSubject())
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			authz := r.Header.Get("Authorization")
			if authz == "" {
				writeAuthError(w, r, http.StatusUnauthorized, errcode.CodeUnauthenticated, "missing authorization header")
				return
			}

			parts := strings.SplitN(authz, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || strings.TrimSpace(parts[1]) == "" {
				writeAuthError(w, r, http.StatusUnauthorized, errcode.CodeUnauthenticated, "invalid authorization header")
				return
			}

			sub, err := cfg.Manager.ParseToken(strings.TrimSpace(parts[1]))
			if err != nil {
				writeAuthError(w, r, http.StatusUnauthorized, errcode.CodeUnauthenticated, "invalid token")
				return
			}

			ctx := authpkg.WithSubject(r.Context(), sub)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func writeAuthError(w http.ResponseWriter, r *http.Request, status, code int, message string) {
	writeJSONError(w, status, map[string]any{
		"code":       code,
		"error_code": errcode.String(code),
		"message":    message,
		"request_id": RequestIDFromContext(r.Context()),
		"path":       r.URL.Path,
	})
}

func writeJSONError(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func shouldSkipPath(path string, exact []string, prefixes []string) bool {
	for _, p := range exact {
		if path == p {
			return true
		}
	}
	for _, p := range prefixes {
		if p != "" && strings.HasPrefix(path, p) {
			return true
		}
	}
	return false
}

func withContext(r *http.Request, ctx context.Context) *http.Request {
	return r.WithContext(ctx)
}
