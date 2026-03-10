package types

const (
	HeaderRequestID = "X-Request-Id"

	LabelManagedBy = "app.kubernetes.io/managed-by"
	LabelComponent = "app.kubernetes.io/component"
	LabelPartOf    = "app.kubernetes.io/part-of"
	LabelName      = "app.kubernetes.io/name"
	LabelInstance  = "app.kubernetes.io/instance"
	LabelVersion   = "app.kubernetes.io/version"

	SystemName          = "gpu-scheduler-platform"
	ComponentAPIServer  = "api-server"
	ComponentScheduler  = "scheduler"
	ComponentController = "controller"
	ComponentWebhook    = "webhook"
	ComponentAgent      = "agent"
)
