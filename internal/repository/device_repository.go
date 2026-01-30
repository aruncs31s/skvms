package repository

import (
	"context"

	"github.com/aruncs31s/skvms/internal/dto"
	"gorm.io/gorm"
)

type DeviceRepository interface {
	ListDevices(ctx context.Context) ([]dto.DeviceView, error)
	GetDevice(ctx context.Context, id uint) (*dto.DeviceView, error)
}

type deviceRepository struct {
	db *gorm.DB
}

func NewDeviceRepository(db *gorm.DB) DeviceRepository {
	return &deviceRepository{db: db}
}

func (r *deviceRepository) ListDevices(ctx context.Context) ([]dto.DeviceView, error) {
	var devices []dto.DeviceView
	err := r.db.WithContext(ctx).
		Table("devices d").
		Select("d.id, d.name, d.device_type AS type, device_details.ip_address, device_details.mac_address, device_details.firmware_version, device_address.address, device_address.city").
		Joins("LEFT JOIN device_details ON device_details.device_id = d.id").
		Joins("LEFT JOIN device_address ON device_address.device_id = d.id").
		Scan(&devices).Error
	if err != nil {
		return nil, err
	}
	return devices, nil
}

func (r *deviceRepository) GetDevice(ctx context.Context, id uint) (*dto.DeviceView, error) {
	var device dto.DeviceView
	err := r.db.WithContext(ctx).
		Table("devices").
		Select("devices.id, devices.name, devices.device_type AS type, device_details.ip_address, device_details.mac_address, device_details.firmware_version, device_address.address, device_address.city").
		Joins("LEFT JOIN device_details ON device_details.device_id = devices.id").
		Joins("LEFT JOIN device_address ON device_address.device_id = devices.id").
		Where("devices.id = ?", id).
		Scan(&device).Error
	if err != nil {
		return nil, err
	}
	if device.ID == 0 {
		return nil, nil
	}
	return &device, nil
}
