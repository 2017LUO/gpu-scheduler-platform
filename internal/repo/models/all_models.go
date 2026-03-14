package model

func AllModels() []any {
	return []any{
		&Tenant{},
		&Queue{},
		&Node{},
		&GPUJob{},
		&GPUJobEvent{},
		&NodeSnapshot{},
		&GPUDevice{},
		&GPUMIGDevice{},
		&PodGPUBindingRuntime{},
		&NodeHeartbeat{},
		&Reservation{},
		&Allocation{},
		&Binding{},
		&GPUQuota{},
		&GPUPolicy{},
		&SchedulingAttempt{},
		&JobRetry{},
		&AuditLog{},
		&Outbox{},
	}
}
