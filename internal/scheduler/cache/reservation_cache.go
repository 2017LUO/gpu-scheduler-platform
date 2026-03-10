package cache

import (
	"sync"
	"time"
)

type reservationItem struct {
	jobID    string
	expireAt time.Time
}

type ReservationCache struct {
	mu    sync.RWMutex
	items map[string]reservationItem
}

func NewReservationCache() *ReservationCache {
	return &ReservationCache{
		items: make(map[string]reservationItem),
	}
}

func (c *ReservationCache) Put(jobID string, expireAt time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[jobID] = reservationItem{
		jobID:    jobID,
		expireAt: expireAt,
	}
}

func (c *ReservationCache) Exists(jobID string, now time.Time) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	it, ok := c.items[jobID]
	if !ok {
		return false
	}
	return now.Before(it.expireAt)
}

func (c *ReservationCache) Delete(jobID string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, jobID)
}
