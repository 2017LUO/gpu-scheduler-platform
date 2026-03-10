package syncer

import "context"

type JobStatusSyncer struct{}

func NewJobStatusSyncer() *JobStatusSyncer {
	return &JobStatusSyncer{}
}

func (s *JobStatusSyncer) Sync(ctx context.Context) error {
	_ = ctx
	return nil
}
