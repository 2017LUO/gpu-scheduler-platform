package allocation

import "time"

type Binding struct {
	ID        string
	JobID     string
	NodeName  string
	GPUIDs    []string
	PodName   string
	Namespace string
	CreatedAt time.Time
}
