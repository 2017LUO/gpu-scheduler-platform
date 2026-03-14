package auth

import "strings"

type Permission string

const (
	PermissionAny            Permission = "*"
	PermissionNone           Permission = ""
	PermissionJobsRead       Permission = "jobs:read"
	PermissionJobsWrite      Permission = "jobs:write"
	PermissionTenantsRead    Permission = "tenants:read"
	PermissionTenantsWrite   Permission = "tenants:write"
	PermissionQueuesRead     Permission = "queues:read"
	PermissionQueuesWrite    Permission = "queues:write"
	PermissionQuotasRead     Permission = "quotas:read"
	PermissionQuotasWrite    Permission = "quotas:write"
	PermissionPoliciesRead   Permission = "policies:read"
	PermissionPoliciesWrite  Permission = "policies:write"
	PermissionClusterRead    Permission = "cluster:read"
	PermissionAgentReport    Permission = "agent:report"
	PermissionAgentHeartbeat Permission = "agent:heartbeat"
)

func PermissionForRequest(method, path string) Permission {
	switch {
	case strings.HasPrefix(path, "/api/v1/jobs"):
		if method == "GET" {
			return PermissionJobsRead
		}
		return PermissionJobsWrite

	case strings.HasPrefix(path, "/api/v1/tenants"):
		if method == "GET" {
			return PermissionTenantsRead
		}
		return PermissionTenantsWrite

	case strings.HasPrefix(path, "/api/v1/queues"):
		if method == "GET" {
			return PermissionQueuesRead
		}
		return PermissionQueuesWrite

	case strings.HasPrefix(path, "/api/v1/quotas"):
		if method == "GET" {
			return PermissionQuotasRead
		}
		return PermissionQuotasWrite

	case strings.HasPrefix(path, "/api/v1/policies"):
		if method == "GET" {
			return PermissionPoliciesRead
		}
		return PermissionPoliciesWrite

	case strings.HasPrefix(path, "/api/v1/clusters"):
		return PermissionClusterRead

	case path == "/internal/agent/report":
		return PermissionAgentReport

	case path == "/internal/agent/heartbeat":
		return PermissionAgentHeartbeat

	default:
		return PermissionNone
	}
}

func ActionForRequest(method, path string) string {
	switch {
	case strings.HasPrefix(path, "/api/v1/jobs"):
		if method == "GET" {
			return "job.read"
		}
		return "job.write"

	case strings.HasPrefix(path, "/api/v1/tenants"):
		if method == "GET" {
			return "tenant.read"
		}
		return "tenant.write"

	case strings.HasPrefix(path, "/api/v1/queues"):
		if method == "GET" {
			return "queue.read"
		}
		return "queue.write"

	case strings.HasPrefix(path, "/api/v1/quotas"):
		if method == "GET" {
			return "quota.read"
		}
		return "quota.write"

	case strings.HasPrefix(path, "/api/v1/policies"):
		if method == "GET" {
			return "policy.read"
		}
		return "policy.write"

	case strings.HasPrefix(path, "/api/v1/clusters"):
		return "cluster.read"

	case path == "/internal/agent/report":
		return "agent.report"

	case path == "/internal/agent/heartbeat":
		return "agent.heartbeat"

	default:
		return strings.ToLower(method) + ".request"
	}
}

func ResourceTypeForPath(path string) string {
	switch {
	case strings.HasPrefix(path, "/api/v1/jobs"):
		return "gpu_job"
	case strings.HasPrefix(path, "/api/v1/tenants"):
		return "tenant"
	case strings.HasPrefix(path, "/api/v1/queues"):
		return "queue"
	case strings.HasPrefix(path, "/api/v1/quotas"):
		return "gpu_quota"
	case strings.HasPrefix(path, "/api/v1/policies"):
		return "gpu_policy"
	case strings.HasPrefix(path, "/api/v1/clusters"):
		return "cluster_node"
	case strings.HasPrefix(path, "/internal/agent/"):
		return "node"
	default:
		return "request"
	}
}
