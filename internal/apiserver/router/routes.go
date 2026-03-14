package router

import (
	"net/http"

	apihandler "gpu-scheduler-platform/internal/apiserver/handler"
	appcfg "gpu-scheduler-platform/internal/config"
)

type Routes struct {
	handlers *apihandler.Handlers
	features appcfg.APIFeaturesConfig
}

func NewRoutes(handlers *apihandler.Handlers, features appcfg.APIFeaturesConfig) *Routes {
	return &Routes{
		handlers: handlers,
		features: features,
	}
}

func (r *Routes) Register(root *http.ServeMux) {
	if r == nil || root == nil || r.handlers == nil {
		return
	}

	if r.handlers.Job != nil {
		root.HandleFunc("POST /api/v1/jobs", r.handlers.Job.Create)
		root.HandleFunc("GET /api/v1/jobs", r.handlers.Job.List)
		root.HandleFunc("GET /api/v1/jobs/{id}", r.handlers.Job.Get)
		root.HandleFunc("GET /api/v1/jobs/{id}/events", r.handlers.Job.ListEvents)
	}

	if r.handlers.InternalAgent != nil {
		root.HandleFunc("POST /internal/agent/heartbeat", r.handlers.InternalAgent.Heartbeat)
		root.HandleFunc("POST /internal/agent/report", r.handlers.InternalAgent.Report)
	}

	if r.features.EnableTenantAPI && r.handlers.Tenant != nil {
		root.HandleFunc("POST /api/v1/tenants", r.handlers.Tenant.Create)
		root.HandleFunc("GET /api/v1/tenants", r.handlers.Tenant.List)
		root.HandleFunc("GET /api/v1/tenants/{id}", r.handlers.Tenant.Get)
		root.HandleFunc("PUT /api/v1/tenants/{id}", r.handlers.Tenant.Update)
		root.HandleFunc("DELETE /api/v1/tenants/{id}", r.handlers.Tenant.Delete)
	}

	if r.features.EnableQueueAPI && r.handlers.Queue != nil {
		root.HandleFunc("PUT /api/v1/queues", r.handlers.Queue.Upsert)
		root.HandleFunc("GET /api/v1/queues", r.handlers.Queue.List)
		root.HandleFunc("GET /api/v1/queues/{tenantID}/{name}", r.handlers.Queue.Get)
		root.HandleFunc("DELETE /api/v1/queues/{tenantID}/{name}", r.handlers.Queue.Delete)
	}

	if r.features.EnableQuotaAPI && r.handlers.Quota != nil {
		root.HandleFunc("PUT /api/v1/quotas", r.handlers.Quota.Upsert)
		root.HandleFunc("GET /api/v1/quotas/{tenantID}", r.handlers.Quota.ListByTenant)
		root.HandleFunc("GET /api/v1/quotas/{tenantID}/{namespace}", r.handlers.Quota.Get)
		root.HandleFunc("DELETE /api/v1/quotas/{tenantID}/{namespace}", r.handlers.Quota.Delete)
	}

	if r.features.EnablePolicyAPI && r.handlers.Policy != nil {
		root.HandleFunc("PUT /api/v1/policies", r.handlers.Policy.Upsert)
		root.HandleFunc("GET /api/v1/policies", r.handlers.Policy.List)
		root.HandleFunc("GET /api/v1/policies/{tenantID}/{name}", r.handlers.Policy.Get)
		root.HandleFunc("DELETE /api/v1/policies/{tenantID}/{name}", r.handlers.Policy.Delete)
	}

	if r.features.EnableClusterAPI && r.handlers.Cluster != nil {
		root.HandleFunc("GET /api/v1/clusters/nodes", r.handlers.Cluster.ListNodes)
		root.HandleFunc("GET /api/v1/clusters/nodes/{nodeName}", r.handlers.Cluster.GetNode)
	}
}
