package syncer

import "context"

type ClusterSyncer struct{}

func NewClusterSyncer() *ClusterSyncer {
	return &ClusterSyncer{}
}

func (s *ClusterSyncer) Sync(ctx context.Context) error {
	_ = ctx
	return nil
}
