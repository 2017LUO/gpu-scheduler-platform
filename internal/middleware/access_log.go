package middleware

import (
	"net/http"
	"time"

	"gpu-scheduler-platform/internal/observability/logging"

	"go.uber.org/zap"
)

type AccessLogConfig struct {
	Enabled bool
}

type responseRecorder struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (r *responseRecorder) WriteHeader(statusCode int) {
	r.status = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	if r.status == 0 {
		r.status = http.StatusOK
	}
	n, err := r.ResponseWriter.Write(b)
	r.bytes += n
	return n, err
}

func AccessLog(lg *zap.Logger, cfg AccessLogConfig) func(http.Handler) http.Handler {
	logger := logging.LoggerOrNop(lg)

	return func(next http.Handler) http.Handler {
		if !cfg.Enabled {
			return next
		}
		if next == nil {
			next = http.NotFoundHandler()
		}

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rec := &responseRecorder{ResponseWriter: w}

			next.ServeHTTP(rec, r)

			if rec.status == 0 {
				rec.status = http.StatusOK
			}

			logger.Info("http access",
				zap.String("request_id", RequestIDFromContext(r.Context())),
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status", rec.status),
				zap.Duration("latency", time.Since(start)),
				zap.Int("bytes", rec.bytes),
				zap.String("client_ip", clientIP(r)),
				zap.String("user_agent", r.UserAgent()),
			)
		})
	}
}

func clientIP(r *http.Request) string {
	if r == nil {
		return ""
	}
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}
	if xrip := r.Header.Get("X-Real-IP"); xrip != "" {
		return xrip
	}
	return r.RemoteAddr
}
