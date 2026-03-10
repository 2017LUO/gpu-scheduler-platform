package allocation

import "time"

type Reservation struct {
	ID        string
	JobID     string
	NodeName  string
	GPUIDs    []string
	ExpireAt  time.Time
	CreatedAt time.Time
}
