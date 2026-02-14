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
	CreateDeviceType(
		ctx context.Context,
		deviceType *model.DeviceTypes,
	) error
	GetDeviceByID(
		ctx context.Context,
		id uint,
	) (*model.DeviceTypes, error)
	GetDeviceTypeByDeviceID(
		ctx context.Context,
		deviceID uint,
	) (map[string]interface{}, error)
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
func (r *deviceTypesRepository) CreateDeviceType(
	ctx context.Context,
	deviceType *model.DeviceTypes,
) error {
	return r.db.WithContext(ctx).Create(deviceType).Error
}
func (r *deviceTypesRepository) GetDeviceByID(
	ctx context.Context,
	id uint,
) (*model.DeviceTypes, error) {
	var deviceType model.DeviceTypes
	err := r.db.WithContext(ctx).First(&deviceType, id).Error
	if err != nil {
		return nil, err
	}
	return &deviceType, nil
}
func (r *deviceTypesRepository) GetDeviceTypeByDeviceID(
	ctx context.Context,
	deviceID uint,
) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := r.db.WithContext(ctx).
		Table("device_types").
		Select(
			[]string{
				"device_types.name",
				"device_types.hardware_type",
			},
		).
		Joins("join devices on devices.device_type = device_types.id").
		Where("devices.id = ?", deviceID).
		Take(&result).Error
	if err != nil {
		return nil, nil
	}
	return result, nil
}
