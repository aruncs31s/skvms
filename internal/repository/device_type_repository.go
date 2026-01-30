package repository

import (
	"context"

	"github.com/aruncs31s/skvms/internal/model"
	"gorm.io/gorm"
)

type DeviceTypesRepository interface {
	ListDeviceTypes(
		ctx context.Context,
		limit int,
		offset int,
	) ([]model.DeviceTypes, error)
}
type deviceTypesRepository struct {
	db *gorm.DB
}

func NewDeviceTypesRepository(db *gorm.DB) DeviceTypesRepository {
	return &deviceTypesRepository{db: db}
}
func (r *deviceTypesRepository) ListDeviceTypes(
	ctx context.Context,
	limit int,
	offset int,
) ([]model.DeviceTypes, error) {
	var deviceTypes []model.DeviceTypes
	err := r.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Find(&deviceTypes).Error
	if err != nil {
		return nil, err
	}
	return deviceTypes, nil
}
