package profiling

import (
	"net/http"
	"net/http/pprof"
	"strings"

	appcfg "gpu-scheduler-platform/internal/config"
)

// RegisterPprofRoutes 根据配置将 pprof 路由注册到指定 mux。
// 当 cfg.Enabled=false 时，不注册任何路由。
func RegisterPprofRoutes(mux *http.ServeMux, cfg appcfg.PProfConfig) {
	if mux == nil || !cfg.Enabled {
		return
	}

	prefix := normalizePrefix(cfg.PathPrefix)
	base := strings.TrimSuffix(prefix, "/")

	// 根入口：兼容 /debug/pprof 和 /debug/pprof/
	mux.HandleFunc(prefix, pprof.Index)
	if !strings.HasSuffix(prefix, "/") {
		mux.HandleFunc(prefix+"/", pprof.Index)
	}

	// 标准端点
	mux.HandleFunc(base+"/cmdline", pprof.Cmdline)
	mux.HandleFunc(base+"/profile", pprof.Profile)
	mux.HandleFunc(base+"/symbol", pprof.Symbol)
	mux.HandleFunc(base+"/trace", pprof.Trace)

	// 常见 profile
	mux.Handle(base+"/allocs", pprof.Handler("allocs"))
	mux.Handle(base+"/block", pprof.Handler("block"))
	mux.Handle(base+"/goroutine", pprof.Handler("goroutine"))
	mux.Handle(base+"/heap", pprof.Handler("heap"))
	mux.Handle(base+"/mutex", pprof.Handler("mutex"))
	mux.Handle(base+"/threadcreate", pprof.Handler("threadcreate"))
}

// normalizePrefix 统一规范 pprof 路径前缀：
// 1. 为空时默认 /debug/pprof
// 2. 必须以 / 开头
// 3. 去掉尾部 /
// 4. 如果传入 "/"，回退为默认值
func normalizePrefix(prefix string) string {
	prefix = strings.TrimSpace(prefix)
	if prefix == "" {
		return "/debug/pprof"
	}
	if !strings.HasPrefix(prefix, "/") {
		prefix = "/" + prefix
	}
	if len(prefix) > 1 {
		prefix = strings.TrimSuffix(prefix, "/")
	}
	if prefix == "/" {
		return "/debug/pprof"
	}
	return prefix
}
