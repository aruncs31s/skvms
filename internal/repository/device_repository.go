package repository

import (
	"context"

	"github.com/aruncs31s/skvms/internal/dto"
	"github.com/aruncs31s/skvms/internal/model"
	"gorm.io/gorm"
)

type DeviceRepository interface {
	ListDevices(ctx context.Context) ([]dto.DeviceView, error)
	ListDevicesByUser(ctx context.Context, userID uint) ([]dto.DeviceView, error)
	GetDeviceForUpdate(
		ctx context.Context,
		tx *gorm.DB,
		id uint) (*model.Device, error)
	GetDevice(
		ctx context.Context,
		id uint,
	) (*model.Device, error)
	CreateDevice(
		ctx context.Context,
		device *model.Device,
		details *model.DeviceDetails,
		address *model.DeviceAddress,
	) (*model.Device, error)
	UpdateDevice(
		ctx context.Context,
		device *model.Device,
		details *model.DeviceDetails,
		address *model.DeviceAddress,
	) error
	DeleteDevice(ctx context.Context, id uint) error
	FindVersionByVersion(ctx context.Context, version string) (*model.Version, error)
	FindVersionByID(
		ctx context.Context,
		id uint,
	) (
		*model.Version,
		error,
	)
	AddConnectedDevice(ctx context.Context, parentID, childID uint) error
	GetConnectedDevices(ctx context.Context, parentID uint) ([]dto.DeviceView, error)
	Count(ctx context.Context) (int64, error)
	CountActive(ctx context.Context) (int64, error)
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
		Select("d.id, d.name, COALESCE(dt.type_name, 'Unknown') AS type, device_details.ip_address, device_details.mac_address, device_details.firmware_version, d.version_id AS version_id, device_address.address, device_address.city, d.device_state AS device_state").
		Joins("LEFT JOIN device_details ON device_details.device_id = d.id").
		Joins("LEFT JOIN device_address ON device_address.device_id = d.id").
		Joins("LEFT JOIN device_types dt ON dt.id = d.device_type").
		Joins("LEFT JOIN device_states ds ON d.device_state  = ds.id").
		Scan(&devices).Error
	if err != nil {
		return nil, err
	}
	return devices, nil
}

func (r *deviceRepository) ListDevicesByUser(ctx context.Context, userID uint) ([]dto.DeviceView, error) {
	var devices []dto.DeviceView
	err := r.db.WithContext(ctx).
		Table("devices d").
		Select("d.id, d.name, COALESCE(dt.type_name, 'Unknown') AS type, device_details.ip_address, device_details.mac_address, device_details.firmware_version, d.version_id AS version_id, device_address.address, device_address.city, d.device_state AS device_state").
		Joins("LEFT JOIN device_details ON device_details.device_id = d.id").
		Joins("LEFT JOIN device_address ON device_address.device_id = d.id").
		Joins("LEFT JOIN device_types dt ON dt.id = d.device_type").
		Where("d.created_by = ?", userID).
		Scan(&devices).Error
	if err != nil {
		return nil, err
	}
	return devices, nil
}

func (r *deviceRepository) GetDeviceForUpdate(
	ctx context.Context,
	tx *gorm.DB,
	id uint,
) (*model.Device, error) {
	var device model.Device
	err := tx.
		WithContext(ctx).
		Preload("DeviceState").
		Preload("Details").
		Preload("Address").
		Preload("DeviceType").
		Preload("Version").
		First(&device, id).Error
	if err != nil {
		return nil, err
	}
	if device.ID == 0 {
		return nil, nil
	}
	return &device, nil
}

func (r *deviceRepository) GetDevice(ctx context.Context, id uint) (*model.Device, error) {
	var device model.Device
	err := r.db.WithContext(ctx).
		Preload("DeviceState").
		Preload("Details").
		Preload("Address").
		Preload("DeviceType").
		Preload("Version").
		First(&device, id).Error
	if err != nil {
		return nil, err
	}
	if device.ID == 0 {
		return nil, nil
	}
	return &device, nil
}

func (r *deviceRepository) CreateDevice(ctx context.Context, device *model.Device, details *model.DeviceDetails, address *model.DeviceAddress) (*model.Device, error) {

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
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
	if err != nil {
		return nil, err
	}
	return device, nil
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

func (r *deviceRepository) FindVersionByVersion(ctx context.Context, version string) (*model.Version, error) {
	var v model.Version
	err := r.db.WithContext(ctx).Where("version = ?", version).First(&v).Error
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *deviceRepository) FindVersionByID(ctx context.Context, id uint) (*model.Version, error) {
	var v model.Version
	err := r.db.WithContext(ctx).First(&v, id).Error
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *deviceRepository) AddConnectedDevice(ctx context.Context, parentID, childID uint) error {
	connected := &model.ConnectedDevice{
		ParentID: parentID,
		ChildID:  childID,
	}
	return r.db.WithContext(ctx).Create(connected).Error
}

func (r *deviceRepository) GetConnectedDevices(ctx context.Context, parentID uint) ([]dto.DeviceView, error) {
	var connected []model.ConnectedDevice
	err := r.db.WithContext(ctx).Where("parent_id = ?", parentID).Find(&connected).Error
	if err != nil {
		return nil, err
	}

	var devices []dto.DeviceView
	for _, c := range connected {
		device, err := r.GetDevice(ctx, c.ChildID)
		if err != nil {
			return nil, err
		}
		if device != nil {
			dv := dto.DeviceView{
				ID:              device.ID,
				Name:            device.Name,
				Type:            device.DeviceType.Name,
				FirmwareVersion: device.Version.Version,
				IPAddress:       device.Details.IPAddress,
				MACAddress:      device.Details.MACAddress,
				Address:         device.Address.Address,
				City:            device.Address.City,
				DeviceState:     device.DeviceState.Name,
			}
			devices = append(devices, dv)
		}
	}
	return devices, nil
}

func (r *deviceRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Device{}).Count(&count).Error
	return count, err
}

func (r *deviceRepository) CountActive(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Device{}).Where("current_state = ?", "active").Count(&count).Error
	return count, err
}
