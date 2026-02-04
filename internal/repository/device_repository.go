package repository

import (
	"context"

	"github.com/aruncs31s/skvms/internal/dto"
	"github.com/aruncs31s/skvms/internal/model"
	"gorm.io/gorm"
)

type DeviceRepository interface {
	ListDevices(
		ctx context.Context,
	) (
		[]model.DeviceView,
		error,
	)
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
	AddConnectedDevice(
		ctx context.Context,
		parentID,
		childID uint,
	) error
	Count(ctx context.Context) (int64, error)
	CountActive(ctx context.Context) (int64, error)
	DeviceReader
}
type SolarReader interface {
	DeviceReader
}
type DeviceReader interface {
	GetDevicesByHardwareType(
		ctx context.Context,
		hardwareType model.HardwareType,
	) (*[]model.DeviceView, error)
	GetUsersDevicesByHardwareType(
		ctx context.Context,
		hardwareType model.HardwareType,
		userID uint,
	) (*[]model.DeviceView, error)
	GetConnectedDevices(
		ctx context.Context,
		parentID uint,
	) ([]dto.DeviceView, error)
	// I want to get a connected devices by hardware type
	// Example : i want to get my microcontroller connected to the solar charger
	GetConnectedDevicesByHardwareType(
		ctx context.Context,
		// This will be the solar charger id
		parentID uint,
		hardwareType model.HardwareType,
		// Right now only supports connecting one device.
	) (dto.DeviceView, error)
}

type deviceRepository struct {
	db *gorm.DB
}

// Optimize
func NewSolarReader(db *gorm.DB) SolarReader {
	return &deviceRepository{db: db}
}
func NewDeviceRepository(db *gorm.DB) DeviceRepository {
	return &deviceRepository{db: db}
}

func (r *deviceRepository) mapDeviceToDeviceView(d model.Device) dto.DeviceView {
	return dto.DeviceView{
		ID:              d.ID,
		Name:            d.Name,
		Type:            d.DeviceType.Name,
		HardwareType:    d.DeviceType.HardwareType.String(),
		Status:          d.DeviceState.Name,
		IPAddress:       d.Details.IPAddress,
		MACAddress:      d.Details.MACAddress,
		FirmwareVersion: d.Version.Name,
		Address:         d.Address.Address,
		City:            d.Address.City,
	}
}

func (r *deviceRepository) ListDevices(
	ctx context.Context,
) (
	[]model.DeviceView,
	error,
) {
	var devices []model.DeviceView

	query := r.db.Table("devices as d")
	query.Select([]string{
		"d.id",
		"d.name",
		"COALESCE(dt.name, 'Unknown') AS type",
		"dt.hardware_type as hardware_type",
		"details.ip_address",
		"details.mac_address",
		"v.name  as firmware_version",
		"device_address.address",
		"device_address.city",
		"ds.name  as current_state",
	}).
		Joins("JOIN device_details details ON details.device_id = d.id").
		Joins("LEFT JOIN device_address ON device_address.device_id = d.id").
		Joins("JOIN device_types dt ON dt.id = d.device_type").
		Joins("JOIN device_states ds ON d.current_state  = ds.id").
		Joins("JOIN versions v ON v.id = d.version_id")

	err := query.Scan(&devices).Error
	if err != nil {
		return nil, err
	}
	return devices, nil
}

func (r *deviceRepository) ListDevicesByUser(ctx context.Context, userID uint) ([]dto.DeviceView, error) {
	var devices []model.Device
	err := r.db.WithContext(ctx).
		Where("created_by = ?", userID).
		Preload("DeviceType").
		Preload("Version").
		Preload("Details").
		Preload("Address").
		Preload("DeviceState").
		Find(&devices).Error
	if err != nil {
		return nil, err
	}
	var dtoDevices []dto.DeviceView
	for _, d := range devices {
		dtoDevices = append(dtoDevices, r.mapDeviceToDeviceView(d))
	}
	return dtoDevices, nil
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
			dv := r.mapDeviceToDeviceView(*device)
			devices = append(devices, dv)
		}
	}
	return devices, nil
}

func (r *deviceRepository) GetConnectedDevicesByHardwareType(
	ctx context.Context,
	// This will be the solar charger id
	parentID uint,
	// Right now only supports connecting one device.
) (dto.DeviceView, error) {
	var connected model.ConnectedDevice
	err := r.
		db.
		WithContext(ctx).
		Preload("Device").
		Where("parent_id = ? , ", parentID).
		First(&connected).Error

	if err != nil {
		return dto.DeviceView{}, err
	}

	device, err := r.GetDevice(ctx, connected.ChildID)
	var dv dto.DeviceView
	if device != nil {
		dv = r.mapDeviceToDeviceView(*device)
	}
	return dv, nil
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

type DeviceView struct {
	ID              uint         `gorm:"column:id"`
	Name            string       `gorm:"column:name"`
	Type            string       `gorm:"column:type"`
	HardwareType    HardwareType `gorm:"column:hardware_type"`
	IPAddress       string       `gorm:"column:ip_address"`
	MACAddress      string       `gorm:"column:mac_address"`
	FirmwareVersion string       `gorm:"column:firmware_version"`
	Address         string       `gorm:"column:address"`
	City            string       `gorm:"column:city"`
	DeviceState     string       `gorm:"column:current_state"`
}

func (r *deviceRepository) GetDevicesByHardwareType(
	ctx context.Context,
	hardwareType model.HardwareType,
) (*[]model.DeviceView, error) {
	var devices []model.DeviceView

	query := r.db.Table("devices as d")
	query.Select([]string{
		"d.id",
		"d.name",
		"COALESCE(dt.name, 'Unknown') AS type",
		"dt.hardware_type as hardware_type",
		"details.ip_address",
		"details.mac_address",
		"v.name  as firmware_version",
		"device_address.address",
		"device_address.city",
		"ds.name  as current_state",
	}).
		Joins("JOIN device_details details ON details.device_id = d.id").
		Joins("LEFT JOIN device_address ON device_address.device_id = d.id").
		Joins("JOIN device_types dt ON dt.id = d.device_type").
		Joins("JOIN device_states ds ON d.current_state  = ds.id").
		Joins("JOIN versions v ON v.id = d.version_id")

	query.Where(
		"dt.harware_type ? =", hardwareType,
	)
	err := query.Scan(&devices).Error
	if err != nil {
		return nil, err
	}
	return &devices, nil
}

func (r *deviceRepository) GetUsersDevicesByHardwareType(
	ctx context.Context,
	hardwareType model.HardwareType,
	userIDs uint,
) (*[]model.DeviceView, error) {
	var devices []model.DeviceView

	query := r.db.Table("devices as d")
	query.Select([]string{
		"d.id",
		"d.name",
		"COALESCE(dt.name, 'Unknown') AS type",
		"dt.hardware_type as hardware_type",
		"details.ip_address",
		"details.mac_address",
		"v.name  as firmware_version",
		"device_address.address",
		"device_address.city",
		"ds.name  as current_state",
	}).
		Joins("JOIN device_details details ON details.device_id = d.id").
		Joins("LEFT JOIN device_address ON device_address.device_id = d.id").
		Joins("JOIN device_types dt ON dt.id = d.device_type").
		Joins("JOIN device_states ds ON d.current_state  = ds.id").
		Joins("JOIN versions v ON v.id = d.version_id")

	query.Where(
		"dt.harware_type ? =", hardwareType,
	)
	query.Where("d.created_by = ?", userIDs)

	err := query.Scan(&devices).Error
	if err != nil {
		return nil, err
	}
	return &devices, nil
}
