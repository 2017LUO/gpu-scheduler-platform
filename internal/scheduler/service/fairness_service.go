package service

import (
	"context"
	"sort"

	model "gpu-scheduler-platform/internal/repo/models"
	repoimpl "gpu-scheduler-platform/internal/repo/mysql"

	"go.uber.org/zap"
)

type FairnessService struct {
	repos  *repoimpl.Repos
	logger *zap.Logger
}

func NewFairnessService(repos *repoimpl.Repos, lg *zap.Logger) *FairnessService {
	if lg == nil {
		lg = zap.NewNop()
	}
	return &FairnessService{
		repos:  repos,
		logger: lg,
	}
}

func (s *FairnessService) Reorder(ctx context.Context, items []model.GPUJob) []model.GPUJob {
	if len(items) <= 1 {
		return items
	}

	type bucketKey struct {
		TenantID string
		Queue    string
	}
	type bucket struct {
		Key    bucketKey
		Weight int
		Items  []model.GPUJob
	}

	grouped := make(map[bucketKey][]model.GPUJob)
	for _, item := range items {
		key := bucketKey{TenantID: item.TenantID, Queue: item.Queue}
		grouped[key] = append(grouped[key], item)
	}

	buckets := make([]bucket, 0, len(grouped))
	for key, jobs := range grouped {
		weight := 1
		if s.repos != nil && s.repos.Queues != nil {
			if q, err := s.repos.Queues.GetByName(ctx, key.TenantID, key.Queue); err == nil && q.Weight > 0 {
				weight = q.Weight
			}
		}
		sort.SliceStable(jobs, func(i, j int) bool {
			pi := EffectivePriority(jobs[i].Priority, jobs[i].CreatedAt, jobs[i].CreatedAt)
			pj := EffectivePriority(jobs[j].Priority, jobs[j].CreatedAt, jobs[j].CreatedAt)
			if pi != pj {
				return pi > pj
			}
			return jobs[i].CreatedAt.Before(jobs[j].CreatedAt)
		})
		buckets = append(buckets, bucket{
			Key:    key,
			Weight: weight,
			Items:  jobs,
		})
	}

	sort.SliceStable(buckets, func(i, j int) bool {
		if buckets[i].Weight != buckets[j].Weight {
			return buckets[i].Weight > buckets[j].Weight
		}
		if buckets[i].Key.TenantID != buckets[j].Key.TenantID {
			return buckets[i].Key.TenantID < buckets[j].Key.TenantID
		}
		return buckets[i].Key.Queue < buckets[j].Key.Queue
	})

	out := make([]model.GPUJob, 0, len(items))
	for {
		progress := false
		for idx := range buckets {
			if len(buckets[idx].Items) == 0 {
				continue
			}
			take := buckets[idx].Weight
			if take <= 0 {
				take = 1
			}
			if take > len(buckets[idx].Items) {
				take = len(buckets[idx].Items)
			}
			out = append(out, buckets[idx].Items[:take]...)
			buckets[idx].Items = buckets[idx].Items[take:]
			progress = true
		}
		if !progress {
			break
		}
	}
	return out
}
