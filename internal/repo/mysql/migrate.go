package mysql

import (
	"fmt"
	model "gpu-scheduler-platform/internal/repo/models"

	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	if db == nil {
		return ErrNilDB
	}
	if err := db.AutoMigrate(model.AllModels()...); err != nil {
		return fmt.Errorf("auto migrate: %w", err)
	}
	return nil
}
