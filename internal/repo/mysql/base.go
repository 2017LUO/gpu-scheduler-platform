package mysql

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type BaseRepo struct {
	db *gorm.DB
}

func NewBaseRepo(db *gorm.DB) (*BaseRepo, error) {
	if db == nil {
		return nil, ErrNilDB
	}
	return &BaseRepo{db: db}, nil
}

func (r *BaseRepo) DB(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}

func (r *BaseRepo) Tx(ctx context.Context, fn func(tx *gorm.DB) error) error {
	if r == nil || r.db == nil {
		return ErrNilDB
	}
	if fn == nil {
		return ErrInvalidArgument
	}
	return r.db.WithContext(ctx).Transaction(fn)
}

func mapDBError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrNotFound
	}
	return err
}
