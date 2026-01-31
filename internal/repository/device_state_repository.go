package repository

import (
	"context"

	"github.com/aruncs31s/skvms/internal/model"
	"gorm.io/gorm"
)

type DeviceStateRepository interface {
	BeginTx() *gorm.DB
	ListDeviceStates(ctx context.Context) ([]model.DeviceState, error)
	GetStatesByType(
		ctx context.Context,
		deviceType int,
	) (*[]model.DeviceState, error)
	GetByID(ctx context.Context, id uint) (*model.DeviceState, error)
	Create(ctx context.Context, deviceState *model.DeviceState) error
	Update(ctx context.Context, deviceState *model.DeviceState) error
	Delete(ctx context.Context, id uint) error
	UpdateDeviceState(
		ctx context.Context,
		tx *gorm.DB,
		deviceID uint,
		stateID uint,
	) error
	InsertStateHistory(
		ctx context.Context,
		tx *gorm.DB,
		history *model.DeviceStateHistory,
	) error
}

type deviceStateRepository struct {
	db *gorm.DB
}

func NewDeviceStateRepository(db *gorm.DB) DeviceStateRepository {
	return &deviceStateRepository{db: db}
}

func (r *deviceStateRepository) BeginTx() *gorm.DB {
	return r.db.Begin()
}

func (r *deviceStateRepository) ListDeviceStates(ctx context.Context) ([]model.DeviceState, error) {
	var deviceStates []model.DeviceState
	err := r.db.WithContext(ctx).Find(&deviceStates).Error
	return deviceStates, err
}
func (r *deviceStateRepository) GetStatesByType(
	ctx context.Context,
	deviceType int,
) (*[]model.DeviceState, error) {
	var deviceStates []model.DeviceState
	err := r.db.WithContext(ctx).
		Where("device_type_id = ?", deviceType).
		Find(&deviceStates).Error
	if err != nil {
		return nil, err
	}
	return &deviceStates, nil
}
func (r *deviceStateRepository) GetByID(ctx context.Context, id uint) (*model.DeviceState, error) {
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

func (r *deviceStateRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.DeviceState{}, id).Error
}
func (r *deviceStateRepository) UpdateDeviceState(
	ctx context.Context,
	tx *gorm.DB,
	deviceID uint,
	stateID uint,
) error {
	return tx.WithContext(ctx).
		Model(&model.Device{}).
		Where("id = ?", deviceID).
		Update("device_state", stateID).Error
}

func (r *deviceStateRepository) InsertStateHistory(
	ctx context.Context,
	tx *gorm.DB,
	history *model.DeviceStateHistory,
) error {
	return tx.WithContext(ctx).Create(history).Error
}
