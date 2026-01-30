package repository

import (
	"context"

	"github.com/aruncs31s/skvms/internal/dto"
	"github.com/aruncs31s/skvms/internal/model"
	"gorm.io/gorm"
)

type DeviceRepository interface {
	ListDevices(ctx context.Context) ([]dto.DeviceView, error)
	GetDevice(ctx context.Context, id uint) (*dto.DeviceView, error)
	CreateDevice(ctx context.Context, device *model.Device, details *model.DeviceDetails, address *model.DeviceAddress) error
	UpdateDevice(ctx context.Context, device *model.Device, details *model.DeviceDetails, address *model.DeviceAddress) error
	DeleteDevice(ctx context.Context, id uint) error
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
		Select("d.id, d.name, dt.type_name AS type, device_details.ip_address, device_details.mac_address, device_details.firmware_version, device_address.address, device_address.city").
		Joins("LEFT JOIN device_details ON device_details.device_id = d.id").
		Joins("LEFT JOIN device_address ON device_address.device_id = d.id").
		Joins("LEFT JOIN device_types dt ON dt.id = d.device_type").
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
		Select("devices.id, devices.name, dt.type_name AS type, device_details.ip_address, device_details.mac_address, device_details.firmware_version, device_address.address, device_address.city").
		Joins("LEFT JOIN device_details ON device_details.device_id = devices.id").
		Joins("LEFT JOIN device_address ON device_address.device_id = devices.id").
		Joins("LEFT JOIN device_types dt ON dt.id = devices.device_type").
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

func (r *deviceRepository) CreateDevice(ctx context.Context, device *model.Device, details *model.DeviceDetails, address *model.DeviceAddress) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(device).Error; err != nil {
			return err
		}
		details.DeviceID = device.ID
		if err := tx.Create(details).Error; err != nil {
			return err
		}
		address.DeviceID = device.ID
		if err := tx.Create(address).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *deviceRepository) UpdateDevice(ctx context.Context, device *model.Device, details *model.DeviceDetails, address *model.DeviceAddress) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(device).Error; err != nil {
			return err
		}
		if err := tx.Where("device_id = ?", device.ID).Save(details).Error; err != nil {
			return err
		}
		if err := tx.Where("device_id = ?", device.ID).Save(address).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *deviceRepository) DeleteDevice(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("device_id = ?", id).Delete(&model.DeviceDetails{}).Error; err != nil {
			return err
		}
		if err := tx.Where("device_id = ?", id).Delete(&model.DeviceAddress{}).Error; err != nil {
			return err
		}
		if err := tx.Delete(&model.Device{}, id).Error; err != nil {
			return err
		}
		return nil
	})
}
