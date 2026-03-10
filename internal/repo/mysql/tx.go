package mysql

import (
	"context"
	"fmt"

	"gpu-scheduler-platform/internal/repo"

	"gorm.io/gorm"
)

type txContextKey struct{}

type txWrapper struct {
	tx *gorm.DB
}

func (t *txWrapper) Commit() error {
	if t == nil || t.tx == nil {
		return nil
	}
	return t.tx.Commit().Error
}

func (t *txWrapper) Rollback() error {
	if t == nil || t.tx == nil {
		return nil
	}
	return t.tx.Rollback().Error
}

type TxManager struct {
	db *gorm.DB
}

func NewTxManager(db *gorm.DB) *TxManager {
	return &TxManager{db: db}
}

func (m *TxManager) Begin(ctx context.Context) (repo.Tx, error) {
	if m == nil || m.db == nil {
		return nil, fmt.Errorf("tx manager db is nil")
	}
	tx := dbFromContext(ctx, m.db).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &txWrapper{tx: tx}, nil
}

func (m *TxManager) WithinTx(ctx context.Context, fn func(ctx context.Context) error) error {
	if m == nil || m.db == nil {
		return fmt.Errorf("tx manager db is nil")
	}
	return dbFromContext(ctx, m.db).WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txCtx := context.WithValue(ctx, txContextKey{}, tx)
		return fn(txCtx)
	})
}
