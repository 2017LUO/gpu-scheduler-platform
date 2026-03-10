package syncer

import "context"

type MetricsSyncer struct{}

func NewMetricsSyncer() *MetricsSyncer {
	return &MetricsSyncer{}
}

func (s *MetricsSyncer) Sync(ctx context.Context) error {
	_ = ctx
	return nil
}
