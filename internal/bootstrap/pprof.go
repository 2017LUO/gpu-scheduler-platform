package bootstrap

import (
	"net/http"

	appcfg "gpu-scheduler-platform/internal/config"
	obsprofiling "gpu-scheduler-platform/internal/observability/profiling"
)

// RegisterPprof 将 pprof 能力接入指定 mux。
// 这里保留 bootstrap 层薄封装，方便 app 层统一通过 bootstrap 调用。
func RegisterPprof(mux *http.ServeMux, cfg appcfg.PProfConfig) {
	if mux == nil {
		return
	}
	obsprofiling.RegisterPprofRoutes(mux, cfg)
}
