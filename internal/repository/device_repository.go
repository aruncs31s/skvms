package repository

import (
	"context"

	"github.com/aruncs31s/skvms/internal/dto"
	"github.com/aruncs31s/skvms/internal/model"
	"github.com/aruncs31s/skvms/utils"
	"gorm.io/gorm"
)

type DeviceRepository interface {
	ListDevices(
		ctx context.Context,
	) (
		[]model.DeviceView,
		error,
	)
	ListDevicesByUser(
		ctx context.Context,
		userID uint,
	) ([]dto.DeviceView, error)
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
		assignment *model.DeviceAssignment,
	) (*model.Device, error)
	UpdateDevice(
		ctx context.Context,
		device *model.Device,
		details *model.DeviceDetails,
		assignment *model.DeviceAssignment,
	) error
	DeleteDevice(
		ctx context.Context,
		id uint,
	) error
	FindVersionByVersion(
		ctx context.Context,
		version string,
	) (*model.Version, error)
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
	SearchDevices(
		ctx context.Context,
		query string,
	) ([]dto.GenericDropdown, error)
	SearchDevicesByHardwareType(
		ctx context.Context,
		query string,
		hardwareType model.HardwareType,
	) ([]dto.GenericDropdown, error)
	SearchDevicesByHardwareTypes(
		ctx context.Context,
		query string,
		hardwareType []model.HardwareType,
	) ([]dto.GenericDropdown, error)
	ListDevicesByHardwareTypes(
		ctx context.Context,
		hardwareTypes []model.HardwareType,
	) ([]model.DeviceView, error)
	DeviceReader
	IsParent(ctx context.Context, parentID uint, childID uint) (bool, error)
	RemoveConnectedDevice(
		ctx context.Context,
		parentID uint,
		childID uint,
	) error
	GetRecentlyCreatedDevices(
		ctx context.Context,
		limit,
		offset int,
	) (*[]model.DeviceView, error)
	GetTotalDeviceCountByType(
		ctx context.Context,
		deviceType model.HardwareType,
	) (int64, error)
	GetDevicesByState(
		ctx context.Context,
		state string,
	) ([]model.DeviceView, error)
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
	GetConnectedDevicesByIDs(
		ctx context.Context,
		parents []uint,
	) (map[uint]model.ConnectedDeviceReadings, error)
	// I want to get a connected devices by hardware type
	// Example : i want to get my microcontroller connected to the solar charger
	GetConnectedDevicesByHardwareType(
		ctx context.Context,
		// This will be the solar charger id
		parentID uint,
		hardwareType model.HardwareType,
		// Right now only supports connecting one device.
	) (dto.DeviceView, error)
	GetDevicesByHardwareTypeAndUserID(
		ctx context.Context,
		hardwareType model.HardwareType,
		userID uint,
	) (*[]model.DeviceView, error)
	// Solar Devices Only
	GetDevicesByLocationID(
		ctx context.Context,
		locationID uint,
	) ([]model.DeviceView, error)
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
		HardwareType:    d.DeviceType.HardwareType,
		Status:          d.DeviceState.Name,
		IPAddress:       d.Details.IPAddress,
		MACAddress:      d.Details.MACAddress,
		FirmwareVersion: d.Version.Name,
	}
}

func (r *deviceRepository) ListDevices(
	ctx context.Context,
) (
	[]model.DeviceView,
	error,
) {
	var devices []model.DeviceView

	query := r.db.
		Table("devices as d")
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

func (r *deviceRepository) CreateDevice(
	ctx context.Context,
	device *model.Device,
	details *model.DeviceDetails,
	assignment *model.DeviceAssignment,
) (*model.Device, error) {

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(device).Error; err != nil {
			return err
		}

		detailsStruct := &model.DeviceDetails{
			DeviceID:   device.ID,
			IPAddress:  details.IPAddress,
			MACAddress: details.MACAddress,
		}
		if err := tx.Create(detailsStruct).Error; err != nil {
			return err
		}
		if assignment != nil {
			assignmentStruct := &model.DeviceAssignment{
				DeviceID:   device.ID,
				LocationID: assignment.LocationID,
			}
			if err := tx.Create(assignmentStruct).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	device.Details = *details
	if assignment != nil {
		device.Assignment = *assignment
	}

	return device, nil
}

func (r *deviceRepository) UpdateDevice(ctx context.Context, device *model.Device, details *model.DeviceDetails, assignment *model.DeviceAssignment) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(device).Error; err != nil {
			return err
		}
		if err := tx.Where("device_id = ?", device.ID).Save(details).Error; err != nil {
			return err
		}
		if assignment != nil {
			if err := tx.Where("device_id = ?", device.ID).Save(assignment).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *deviceRepository) DeleteDevice(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("device_id = ?", id).Delete(&model.DeviceDetails{}).Error; err != nil {
			return err
		}
		if err := tx.Where("device_id = ?", id).Delete(&model.DeviceAssignment{}).Error; err != nil {
			return err
		}
		if err := tx.Delete(&model.Device{}, id).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *deviceRepository) FindVersionByVersion(
	ctx context.Context,
	version string,
) (*model.Version, error) {
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
	hardwareType model.HardwareType,
) (dto.DeviceView, error) {
	var connected model.ConnectedDevice
	err := r.
		db.
		WithContext(ctx).
		Joins("JOIN devices d ON connected_devices.child_id = d.id").
		Joins("JOIN device_types dt ON d.device_type = dt.id").
		Where("connected_devices.parent_id = ? AND dt.hardware_type = ?", parentID, int(hardwareType)).
		First(&connected).Error

	if err != nil {
		return dto.DeviceView{}, err
	}

	device, err := r.GetDevice(ctx, connected.ChildID)
	if err != nil {
		return dto.DeviceView{}, err
	}
	if device == nil {
		return dto.DeviceView{}, gorm.ErrRecordNotFound
	}
	dv := r.mapDeviceToDeviceView(*device)
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
		"locations.name as address",
		"locations.city",
		"ds.name  as current_state",
	}).
		Joins("JOIN device_details details ON details.device_id = d.id").
		Joins("LEFT JOIN device_assignment ON device_assignment.device_id = d.id").
		Joins("LEFT JOIN locations ON locations.id = device_assignment.location_id").
		Joins("JOIN device_types dt ON dt.id = d.device_type").
		Joins("JOIN device_states ds ON d.current_state  = ds.id").
		Joins("JOIN versions v ON v.id = d.version_id")

	query.Where(
		"dt.hardware_type  = ?", hardwareType,
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
		Joins("LEFT JOIN device_assignment ON device_assignment.device_id = d.id").
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

func (r *deviceRepository) SearchDevices(
	ctx context.Context,
	query string,
) ([]dto.GenericDropdown, error) {
	var results []dto.GenericDropdown
	err := r.db.WithContext(ctx).Model(&model.Device{}).
		Select("id, name").
		Where("name LIKE ?", "%"+query+"%").
		Scan(&results).Error
	return results, err
}

func (r *deviceRepository) SearchDevicesByHardwareType(
	ctx context.Context,
	query string,
	hardwareType model.HardwareType,
) ([]dto.GenericDropdown, error) {
	var results []dto.GenericDropdown
	err := r.db.
		WithContext(ctx).
		Table(model.Device{}.TableName()+" d").
		Select("d.id, d.name").
		Joins("JOIN device_types dt ON dt.id = d.device_type").
		Where("d.name LIKE ? AND dt.hardware_type = ?", "%"+query+"%", hardwareType).
		Scan(&results).Error
	return results, err
}
func (r *deviceRepository) SearchDevicesByHardwareTypes(
	ctx context.Context,
	query string,
	hardwareType []model.HardwareType,
) ([]dto.GenericDropdown, error) {
	var results []dto.GenericDropdown
	err := r.db.
		WithContext(ctx).
		Table(model.Device{}.TableName()+" d").
		Select("d.id, d.name").
		Joins("JOIN device_types dt ON dt.id = d.device_type").
		Where("d.name LIKE ? AND dt.hardware_type IN ?", "%"+query+"%", hardwareType).
		Scan(&results).Error
	return results, err
}

func (r *deviceRepository) ListDevicesByHardwareTypes(
	ctx context.Context,
	hardwareTypes []model.HardwareType,
) ([]model.DeviceView, error) {
	var devices []model.DeviceView

	query := r.db.
		WithContext(ctx).
		Table("devices as d")
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
		Joins("JOIN versions v ON v.id = d.version_id").
		Where("dt.hardware_type IN ?",
			func() []int {
				var types []int
				for _, ht := range hardwareTypes {
					types = append(types, int(ht))
				}
				return types
			}(),
		)

	err := query.Scan(&devices).Error
	if err != nil {
		return nil, err
	}
	return devices, nil
}
func (r *deviceRepository) GetDevicesByHardwareTypeAndUserID(
	ctx context.Context,
	hardwareType model.HardwareType,
	userID uint,
) (*[]model.DeviceView, error) {
	var devices []model.DeviceView

	query := r.db.
		WithContext(ctx).
		Table("devices as d")
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
		Joins("JOIN versions v ON v.id = d.version_id").
		Where("dt.hardware_type = ? AND d.created_by = ?", int(hardwareType), userID)

	err := query.Scan(&devices).Error
	if err != nil {
		return nil, err
	}
	return &devices, nil
}
func (r *deviceRepository) IsParent(ctx context.Context, parentID uint, childID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.ConnectedDevice{}).
		Where("parent_id = ? AND child_id = ?", parentID, childID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
func (r *deviceRepository) RemoveConnectedDevice(
	ctx context.Context,
	parentID uint,
	childID uint,
) error {
	err := r.db.WithContext(ctx).
		Where("parent_id = ? AND child_id = ?", parentID, childID).
		Delete(&model.ConnectedDevice{}).Error
	return err
}

func (r *deviceRepository) GetRecentlyCreatedDevices(
	ctx context.Context,
	limit,
	offset int,
) (*[]model.DeviceView, error) {
	var devices []model.DeviceView

	query := r.db.
		WithContext(ctx).
		Table("devices as d")
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
		Joins("JOIN versions v ON v.id = d.version_id").
		Order("d.created_at DESC").
		Limit(limit).
		Offset(offset)

	err := query.Scan(&devices).Error
	if err != nil {
		return nil, err
	}
	return &devices, nil
}

func (r *deviceRepository) GetTotalDeviceCountByType(
	ctx context.Context,
	hardwareType model.HardwareType,
) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Table("devices as d").
		Joins("JOIN device_types dt ON dt.id = d.device_type").
		Where("dt.hardware_type = ?", int(hardwareType)).
		Count(&count).Error
	return count, err
}
func (r *deviceRepository) GetDevicesByState(
	ctx context.Context,
	state string,
) ([]model.DeviceView, error) {

	var devices []model.DeviceView

	query := r.db.
		WithContext(ctx).
		Table("devices as d")
	query.Select([]string{
		"d.id",
		"d.name",
		"COALESCE(dt.name, 'Unknown') AS type",
		"dt.hardware_type as hardware_type",
		"details.ip_address",
		"details.mac_address",
		"v.name  as firmware_version",
		"locations.name as address",
		"locations.city",
		"ds.name  as current_state",
	}).
		Joins("JOIN device_details details ON details.device_id = d.id").
		Joins("LEFT JOIN device_assignment ON device_assignment.device_id = d.id").
		Joins("LEFT JOIN locations ON locations.id = device_assignment.location_id").
		Joins("JOIN device_types dt ON dt.id = d.device_type").
		Joins("JOIN device_states ds ON d.current_state  = ds.id").
		Joins("JOIN versions v ON v.id = d.version_id").
		Order("d.created_at DESC").
		Where("ds.name = ?", state)

	err := query.Scan(&devices).Error
	if err != nil {
		return nil, err
	}
	return devices, nil
}

// TODO Optmize, Not using Every Fields.
func (r *deviceRepository) GetConnectedDevicesByIDs(
	ctx context.Context,
	parents []uint,
	// Key is parent ID
) (map[uint]model.ConnectedDeviceReadings, error) {

	today := utils.GetBeginningOfDay()
	night := utils.GetEndOfDay()
	query := `
		SELECT
		cd.parent_id,
		cd.child_id,
			r.created_at,
			r.voltage,
			r.current,
			AVG(r.current) OVER (
				ORDER BY r.created_at
				ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW
			) AS avg_current,
			AVG(r.voltage) OVER (
					ORDER BY r.created_at
					ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW
			) AS avg_voltage,
			80 / AVG(r.current) OVER (
				ORDER BY r.created_at
				ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW
			) AS estimated_remaining_hours
		FROM readings r
		JOIN connected_devices cd ON r.device_id = cd.child_id
		WHERE r.created_at BETWEEN ? AND ? AND cd.parent_id IN ?
		ORDER BY r.created_at DESC
		LIMIT 1000;
	`
	var readings []model.ConnectedDeviceReadings
	err := r.db.WithContext(ctx).Raw(query, today, night, parents).Scan(&readings).Error
	if err != nil {
		return nil, err
	}

	result := make(map[uint]model.ConnectedDeviceReadings)
	for _, reading := range readings {
		result[reading.ParentDevice] = reading
	}
	return result, nil
}
func (r *deviceRepository) GetDevicesByLocationID(ctx context.Context, locationID uint) ([]model.DeviceView, error) {
	var devices []model.DeviceView

	query := r.db.
		WithContext(ctx).
		Table("devices as d")
	query.Select([]string{
		"d.id",
		"d.name",
		"COALESCE(dt.name, 'Unknown') AS type",
		"dt.hardware_type as hardware_type",
		"details.ip_address",
		"details.mac_address",
		"v.name  as firmware_version",
		"ds.name  as current_state",
	}).
		Joins("JOIN device_details details ON details.device_id = d.id").
		Joins("JOIN device_types dt ON dt.id = d.device_type").
		Joins("JOIN device_states ds ON d.current_state  = ds.id").
		Joins("JOIN versions v ON v.id = d	.version_id").
		Joins("JOIN device_assignment da ON da.device_id = d.id").
		Where("da.location_id = ? AND dt.hardware_type = ?", locationID, int(model.HardwareTypeSolar))

	err := query.Scan(&devices).Error
	if err != nil {
		return nil, err
	}
	return devices, nil
}
