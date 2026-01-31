package repository

import (
	"context"

	"github.com/aruncs31s/skvms/internal/model"
	"gorm.io/gorm"
)

type DeviceStateRepository interface {
	ListDeviceStates(ctx context.Context) ([]model.DeviceState, error)
	GetByID(ctx context.Context, id int) (*model.DeviceState, error)
	Create(ctx context.Context, deviceState *model.DeviceState) error
	Update(ctx context.Context, deviceState *model.DeviceState) error
	Delete(ctx context.Context, id int) error
}

type deviceStateRepository struct {
	db *gorm.DB
}

func NewDeviceStateRepository(db *gorm.DB) DeviceStateRepository {
	return &deviceStateRepository{db: db}
}

func (r *deviceStateRepository) ListDeviceStates(ctx context.Context) ([]model.DeviceState, error) {
	var deviceStates []model.DeviceState
	err := r.db.WithContext(ctx).Find(&deviceStates).Error
	return deviceStates, err
}

func (r *deviceStateRepository) GetByID(ctx context.Context, id int) (*model.DeviceState, error) {
	var deviceState model.DeviceState
	err := r.db.WithContext(ctx).First(&deviceState, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &deviceState, nil
}

func (r *deviceStateRepository) Create(ctx context.Context, deviceState *model.DeviceState) error {
	return r.db.WithContext(ctx).Create(deviceState).Error
}

func (r *deviceStateRepository) Update(ctx context.Context, deviceState *model.DeviceState) error {
	return r.db.WithContext(ctx).Save(deviceState).Error
}

func (r *deviceStateRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&model.DeviceState{}, id).Error
}