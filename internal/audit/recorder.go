package audit

import (
	"context"
	"encoding/json"
	"fmt"
	model "gpu-scheduler-platform/internal/repo/models"
	repoimpl "gpu-scheduler-platform/internal/repo/mysql"
	"time"

	"go.uber.org/zap"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Record struct {
	TenantID     string
	Actor        string
	Action       string
	ResourceType string
	ResourceID   string
	ResourceName string
	Status       string
	RequestID    string
	Detail       map[string]any
}

type Recorder struct {
	repo    *repoimpl.AuditLogRepo
	logger  *zap.Logger
	nowFunc func() time.Time
}

func NewRecorder(db *gorm.DB, lg *zap.Logger) (*Recorder, error) {
	repo, err := repoimpl.NewAuditLogRepo(db)
	if err != nil {
		return nil, err
	}
	if lg == nil {
		lg = zap.NewNop()
	}
	return &Recorder{
		repo:    repo,
		logger:  lg,
		nowFunc: func() time.Time { return time.Now().UTC() },
	}, nil
}

func (r *Recorder) Record(ctx context.Context, rec Record) error {
	if r == nil || r.repo == nil {
		return nil
	}
	if rec.Actor == "" {
		rec.Actor = "anonymous"
	}
	if rec.Action == "" {
		rec.Action = "request"
	}
	if rec.ResourceType == "" {
		rec.ResourceType = "request"
	}
	if rec.ResourceID == "" {
		rec.ResourceID = "unknown"
	}

	var detail datatypes.JSON = datatypes.JSON([]byte("{}"))
	if len(rec.Detail) > 0 {
		if b, err := json.Marshal(rec.Detail); err == nil && len(b) > 0 {
			detail = datatypes.JSON(b)
		}
	}

	return r.repo.Create(ctx, &model.AuditLog{
		TenantID:     rec.TenantID,
		Actor:        rec.Actor,
		Action:       rec.Action,
		ResourceType: rec.ResourceType,
		ResourceID:   rec.ResourceID,
		ResourceName: rec.ResourceName,
		Status:       rec.Status,
		RequestID:    rec.RequestID,
		DetailJSON:   detail,
		CreatedAt:    r.nowFunc(),
	})
}

func (r *Recorder) MustRecord(ctx context.Context, rec Record) {
	if err := r.Record(ctx, rec); err != nil && r.logger != nil {
		r.logger.Warn("record audit log failed", zap.Error(fmt.Errorf("audit record: %w", err)))
	}
}
