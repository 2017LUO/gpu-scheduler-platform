package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"gpu-scheduler-platform/internal/observability/logging"

	"go.uber.org/zap"
)

func Recovery(lg *zap.Logger) func(http.Handler) http.Handler {
	logger := logging.LoggerOrNop(lg)

	return func(next http.Handler) http.Handler {
		if next == nil {
			next = http.NotFoundHandler()
		}

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					logger.Error("panic recovered",
						zap.Any("panic", rec),
						zap.ByteString("stack", debug.Stack()),
						zap.String("request_id", RequestIDFromContext(r.Context())),
						zap.String("method", r.Method),
						zap.String("path", r.URL.Path),
					)

					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					_, _ = w.Write([]byte(fmt.Sprintf(`{"code":500,"message":"internal server error","request_id":"%s"}`, RequestIDFromContext(r.Context()))))
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
