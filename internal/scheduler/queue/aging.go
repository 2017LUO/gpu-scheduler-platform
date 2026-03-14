package queue

import (
	"strings"
	"time"
)

func EffectivePriority(priority string, createdAt time.Time, now time.Time) int {
	base := normalizePriority(priority)
	if createdAt.IsZero() || now.Before(createdAt) {
		return base
	}

	waitMinutes := int(now.Sub(createdAt).Minutes())
	boost := waitMinutes / 10
	if boost > 20 {
		boost = 20
	}
	return base + boost
}

func normalizePriority(p string) int {
	switch strings.ToUpper(strings.TrimSpace(p)) {
	case "CRITICAL":
		return 100
	case "HIGH":
		return 80
	case "MEDIUM":
		return 50
	case "NORMAL":
		return 40
	case "LOW":
		return 20
	default:
		return 10
	}
}
