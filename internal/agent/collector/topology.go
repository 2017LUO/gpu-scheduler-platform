package collector

import "context"

type GPULink struct {
	From string `json:"from"`
	To   string `json:"to"`
	Type string `json:"type"`
}

type TopologyCollector struct{}

func NewTopologyCollector() *TopologyCollector {
	return &TopologyCollector{}
}

func (c *TopologyCollector) Collect(ctx context.Context) ([]GPULink, error) {
	_ = ctx
	return nil, nil
}
