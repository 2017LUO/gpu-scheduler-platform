package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

const (
	HeaderRequestID                = "X-Request-Id"
	ContextRequestIDKey contextKey = "request_id"
)

type contextKey string

func RequestID(next http.Handler) http.Handler {
	if next == nil {
		next = http.NotFoundHandler()
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rid := r.Header.Get(HeaderRequestID)
		if rid == "" {
			rid = uuid.NewString()
		}

		w.Header().Set(HeaderRequestID, rid)
		ctx := context.WithValue(r.Context(), ContextRequestIDKey, rid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequestIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	v, _ := ctx.Value(ContextRequestIDKey).(string)
	return v
}
