package reserve

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	model "gpu-scheduler-platform/internal/repo/models"
	repoimpl "gpu-scheduler-platform/internal/repo/mysql"
	schedcache "gpu-scheduler-platform/internal/scheduler/cache"
	schedframework "gpu-scheduler-platform/internal/scheduler/framework"

	"go.uber.org/zap"
	"gorm.io/datatypes"
)

type Reservation struct {
	repos   *repoimpl.Repos
	cache   *schedcache.ReservationCache
	logger  *zap.Logger
	ttl     time.Duration
	nowFunc func() time.Time
}

func NewReservation(repos *repoimpl.Repos, cache *schedcache.ReservationCache, ttl time.Duration, lg *zap.Logger) *Reservation {
	if lg == nil {
		lg = zap.NewNop()
	}
	if ttl <= 0 {
		ttl = 45 * time.Second
	}
	return &Reservation{
		repos:   repos,
		cache:   cache,
		logger:  lg,
		ttl:     ttl,
		nowFunc: func() time.Time { return time.Now().UTC() },
	}
}

func (p *Reservation) Name() string { return "Reservation" }

func (p *Reservation) Reserve(cs *schedframework.CycleState, job *model.GPUJob, node *model.Node) *schedframework.Status {
	if job == nil || node == nil {
		return schedframework.NewStatus(schedframework.CodeError, "job or node is nil")
	}
	if p.repos == nil || p.repos.Reservations == nil {
		return schedframework.NewStatus(schedframework.CodeError, "reservation repo is nil")
	}

	gpuUUIDs, ok := schedframework.ReadSelectedGPUUUIDs(cs)
	if !ok || len(gpuUUIDs) == 0 {
		return schedframework.NewStatus(schedframework.CodeError, "selected gpu uuids are missing")
	}

	now := p.nowFunc()
	reservationID := fmt.Sprintf("resv-%d", now.UnixNano())

	raw, _ := json.Marshal(gpuUUIDs)

	_ = p.repos.Reservations.DeleteByJobID(context.Background(), job.ID)

	rec := &model.Reservation{
		ID:         reservationID,
		JobID:      job.ID,
		NodeName:   node.NodeName,
		GPUIDsJSON: datatypes.JSON(raw),
		ExpireAt:   now.Add(p.ttl),
		CreatedAt:  now,
	}
	if err := p.repos.Reservations.Create(context.Background(), rec); err != nil {
		return schedframework.NewStatus(schedframework.CodeError, err.Error())
	}

	if p.cache != nil {
		p.cache.Set(&schedcache.Reservation{
			JobID:    job.ID,
			NodeName: node.NodeName,
			GPUIDs:   append([]string(nil), gpuUUIDs...),
		})
	}
	cs.Write(schedframework.StateKeyReservationID, reservationID)
	return nil
}

func (p *Reservation) Unreserve(cs *schedframework.CycleState, job *model.GPUJob, node *model.Node) {
	if job == nil {
		return
	}
	if p.cache != nil {
		p.cache.Delete(job.ID)
	}
	if p.repos != nil && p.repos.Reservations != nil {
		_ = p.repos.Reservations.DeleteByJobID(context.Background(), job.ID)
	}
}
