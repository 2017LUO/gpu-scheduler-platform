package cache

import "sync"

type Reservation struct {
	JobID    string
	NodeName string
	GPUIDs   []string
}

type ReservationCache struct {
	mu    sync.RWMutex
	items map[string]*Reservation
}

func NewReservationCache() *ReservationCache {
	return &ReservationCache{
		items: make(map[string]*Reservation),
	}
}

func (c *ReservationCache) Set(r *Reservation) {
	if r == nil || r.JobID == "" {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	cp := *r
	if r.GPUIDs != nil {
		cp.GPUIDs = append([]string(nil), r.GPUIDs...)
	}
	c.items[r.JobID] = &cp
}

func (c *ReservationCache) Get(jobID string) (*Reservation, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	v, ok := c.items[jobID]
	if !ok {
		return nil, false
	}
	cp := *v
	if v.GPUIDs != nil {
		cp.GPUIDs = append([]string(nil), v.GPUIDs...)
	}
	return &cp, true
}

func (c *ReservationCache) Delete(jobID string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, jobID)
}

func (c *ReservationCache) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]*Reservation)
}
