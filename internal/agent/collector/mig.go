package collector

import "context"

type MIGInfo struct {
	ID         string `json:"id"`
	ParentUUID string `json:"parent_uuid"`
	Profile    string `json:"profile"`
	MemoryMiB  int64  `json:"memory_mib"`
}

type MIGCollector struct{}

func NewMIGCollector() *MIGCollector {
	return &MIGCollector{}
}

func (c *MIGCollector) Collect(ctx context.Context) ([]MIGInfo, error) {
	_ = ctx
	return nil, nil
}
