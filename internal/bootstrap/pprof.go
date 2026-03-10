package bootstrap

import (
	"net/http"
	"net/http/pprof"
	"strings"

	appcfg "gpu-scheduler-platform/internal/config"
)

func MountPProf(mux *http.ServeMux, cfg appcfg.PProfConfig) {
	if mux == nil || !cfg.Enabled {
		return
	}

	prefix := strings.TrimRight(cfg.PathPrefix, "/")
	if prefix == "" {
		prefix = "/debug/pprof"
	}

	mux.HandleFunc(prefix+"/", pprof.Index)
	mux.HandleFunc(prefix+"/cmdline", pprof.Cmdline)
	mux.HandleFunc(prefix+"/profile", pprof.Profile)
	mux.HandleFunc(prefix+"/symbol", pprof.Symbol)
	mux.HandleFunc(prefix+"/trace", pprof.Trace)
	mux.Handle(prefix+"/allocs", pprof.Handler("allocs"))
	mux.Handle(prefix+"/block", pprof.Handler("block"))
	mux.Handle(prefix+"/goroutine", pprof.Handler("goroutine"))
	mux.Handle(prefix+"/heap", pprof.Handler("heap"))
	mux.Handle(prefix+"/mutex", pprof.Handler("mutex"))
	mux.Handle(prefix+"/threadcreate", pprof.Handler("threadcreate"))
}
