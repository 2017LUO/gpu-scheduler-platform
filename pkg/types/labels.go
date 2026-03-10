package types

func CommonLabels(component, version string) map[string]string {
	return map[string]string{
		LabelManagedBy: "helm",
		LabelPartOf:    SystemName,
		LabelName:      SystemName,
		LabelComponent: component,
		LabelVersion:   version,
	}
}
