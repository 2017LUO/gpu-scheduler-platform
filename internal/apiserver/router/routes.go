package router

import (
	"net/http"
	"strings"

	"gpu-scheduler-platform/internal/apiserver/handler"
)

type Routes struct {
	JobHandler           *handler.JobHandler
	InternalAgentHandler *handler.InternalAgentHandler
}

func NewRoutes(jobHandler *handler.JobHandler, internalAgentHandler *handler.InternalAgentHandler) *Routes {
	return &Routes{
		JobHandler:           jobHandler,
		InternalAgentHandler: internalAgentHandler,
	}
}

func (rt *Routes) Register(mux *http.ServeMux) {
	if mux == nil || rt == nil {
		return
	}

	mux.HandleFunc("/api/v1/jobs", rt.handleJobs)
	mux.HandleFunc("/api/v1/jobs/", rt.handleJobSubroutes)

	mux.HandleFunc("/internal/agent/report", rt.handleInternalAgentReport)
}

func (rt *Routes) handleJobs(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		rt.JobHandler.Create(w, r)
	case http.MethodGet:
		rt.JobHandler.List(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (rt *Routes) handleJobSubroutes(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimSuffix(r.URL.Path, "/")

	if strings.HasSuffix(path, "/events") {
		if r.Method == http.MethodGet {
			rt.JobHandler.ListEvents(w, r)
			return
		}
		http.NotFound(w, r)
		return
	}

	if r.Method == http.MethodGet {
		rt.JobHandler.GetByID(w, r)
		return
	}

	http.NotFound(w, r)
}

func (rt *Routes) handleInternalAgentReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}
	rt.InternalAgentHandler.Report(w, r)
}
