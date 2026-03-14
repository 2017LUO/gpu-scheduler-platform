package auth

type Authorizer struct {
	rolePermissions map[string]map[Permission]struct{}
}

func NewAuthorizer() *Authorizer {
	return &Authorizer{
		rolePermissions: map[string]map[Permission]struct{}{
			"admin": {
				PermissionAny: {},
			},
			"tenant-admin": {
				PermissionJobsRead:      {},
				PermissionJobsWrite:     {},
				PermissionTenantsRead:   {},
				PermissionQueuesRead:    {},
				PermissionQueuesWrite:   {},
				PermissionQuotasRead:    {},
				PermissionQuotasWrite:   {},
				PermissionPoliciesRead:  {},
				PermissionPoliciesWrite: {},
				PermissionClusterRead:   {},
			},
			"tenant-viewer": {
				PermissionJobsRead:     {},
				PermissionTenantsRead:  {},
				PermissionQueuesRead:   {},
				PermissionQuotasRead:   {},
				PermissionPoliciesRead: {},
				PermissionClusterRead:  {},
			},
			"node-agent": {
				PermissionAgentReport:    {},
				PermissionAgentHeartbeat: {},
				PermissionClusterRead:    {},
			},
		},
	}
}

func (a *Authorizer) Allowed(sub *Subject, perm Permission) bool {
	if perm == PermissionNone {
		return true
	}
	if sub == nil {
		return false
	}
	if sub.IsSystem {
		return true
	}

	for _, p := range sub.Permissions {
		if Permission(p) == PermissionAny || Permission(p) == perm {
			return true
		}
	}

	for _, role := range sub.Roles {
		ps, ok := a.rolePermissions[role]
		if !ok {
			continue
		}
		if _, ok := ps[PermissionAny]; ok {
			return true
		}
		if _, ok := ps[perm]; ok {
			return true
		}
	}
	return false
}
