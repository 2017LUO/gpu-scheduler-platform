package models

var AllModels = []any{
	&Tenant{},
	&GPUJob{},
	&GPUJobEvent{},
	&NodeSnapshot{},
	&GPUDevice{},
	&Allocation{},
	&Reservation{},
	&Binding{},
	&GPUQuota{},
	&AuditLog{},
	&Outbox{},
}
