package framework

import "sync"

type StateKey string

type StateData interface{}

type CycleState struct {
	mu   sync.RWMutex
	data map[StateKey]StateData
}

func NewCycleState() *CycleState {
	return &CycleState{
		data: make(map[StateKey]StateData),
	}
}

func (c *CycleState) Write(key StateKey, val StateData) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = val
}

func (c *CycleState) Read(key StateKey) (StateData, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.data[key]
	return v, ok
}

func (c *CycleState) Clone() *CycleState {
	c.mu.RLock()
	defer c.mu.RUnlock()

	out := make(map[StateKey]StateData, len(c.data))
	for k, v := range c.data {
		out[k] = v
	}
	return &CycleState{data: out}
}
