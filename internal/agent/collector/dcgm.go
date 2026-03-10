package collector

import "context"

type DCGMCollector struct{}

func NewDCGMCollector() *DCGMCollector {
	return &DCGMCollector{}
}

func (c *DCGMCollector) Collect(ctx context.Context) ([]GPUInfo, error) {
	return NewNvidiaSMICollector().Collect(ctx)
}
