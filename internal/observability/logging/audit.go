package logging

import "go.uber.org/zap"

func AuditFields(actor, action, resourceType, resourceID string) []zap.Field {
	return []zap.Field{
		zap.String("audit_actor", actor),
		zap.String("audit_action", action),
		zap.String("audit_resource_type", resourceType),
		zap.String("audit_resource_id", resourceID),
	}
}
