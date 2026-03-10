package job

type PriorityClass string

const (
	PriorityCritical PriorityClass = "critical"
	PriorityHigh     PriorityClass = "high"
	PriorityNormal   PriorityClass = "normal"
	PriorityLow      PriorityClass = "low"
)

func PriorityValue(p PriorityClass) int {
	switch p {
	case PriorityCritical:
		return 400
	case PriorityHigh:
		return 300
	case PriorityNormal:
		return 200
	case PriorityLow:
		return 100
	default:
		return 0
	}
}
