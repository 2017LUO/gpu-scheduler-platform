package profiling

import (
	"net/http"
	"net/http/pprof"
	"strings"
)

func Mount(mux *http.ServeMux, pathPrefix string) {
	if mux == nil {
		return
	}

	prefix := strings.TrimRight(strings.TrimSpace(pathPrefix), "/")
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
