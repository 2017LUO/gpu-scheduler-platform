package logging

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

func RequestID(v string) zap.Field {
	return zap.String("request_id", v)
}

func TenantID(v string) zap.Field {
	return zap.String("tenant_id", v)
}

func JobID(v string) zap.Field {
	return zap.String("job_id", v)
}

func NodeName(v string) zap.Field {
	return zap.String("node_name", v)
}

func GPUCount(v int) zap.Field {
	return zap.Int("gpu_count", v)
}

func Path(v string) zap.Field {
	return zap.String("path", v)
}

func Method(v string) zap.Field {
	return zap.String("method", v)
}

func StatusCode(v int) zap.Field {
	return zap.Int("status", v)
}

func Duration(v time.Duration) zap.Field {
	return zap.Duration("duration", v)
}

func RemoteIP(v string) zap.Field {
	return zap.String("remote_ip", v)
}

func UserAgent(v string) zap.Field {
	return zap.String("user_agent", v)
}

func HTTPRequest(r *http.Request) []zap.Field {
	if r == nil {
		return nil
	}
	return []zap.Field{
		Method(r.Method),
		Path(r.URL.Path),
		RemoteIP(r.RemoteAddr),
		UserAgent(r.UserAgent()),
	}
}
