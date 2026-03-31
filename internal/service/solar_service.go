package service

import (
	"context"
	"errors"

	"github.com/aruncs31s/skvms/internal/dto"
	"github.com/aruncs31s/skvms/internal/model"
	"github.com/aruncs31s/skvms/internal/repository"
	"gorm.io/gorm"
)

var ErrUserNotFound = errors.New("no user found")

type SolarService interface {
	GetAllSolarDevices(
		ctx context.Context,
	) (*[]dto.SolarDeviceView, error)
	CreateASolarDevice(
		ctx context.Context,
		req dto.CreateSolarDeviceDTO,
		createdBy uint,
	) (dto.DeviceView, error)
	GetAllMySolarDevices(
		ctx context.Context,
		userID uint,
	) (*[]dto.SolarDeviceView, error)
}
type solarService struct {
	solarRepo             repository.SolarRepository
	userRepo              repository.UserReader
	deviceRepo            repository.DeviceRepository
	deviceAssignementRepo repository.DeviceAssignmentRepository
	deviceTypeRepo        repository.DeviceTypesRepository
	deviceStateRepo       repository.DeviceStateRepository
	locationRepo          repository.LocationReader
}

func NewSolarService(
	solarRepo repository.SolarRepository,
	user repository.UserReader,
	deviceRepo repository.DeviceRepository,
	deviceStateRepo repository.DeviceStateRepository,
	deviceTypeRepo repository.DeviceTypesRepository,
	locationRepo repository.LocationReader,
) SolarService {
	return &solarService{
		solarRepo:       solarRepo,
		userRepo:        user,
		deviceRepo:      deviceRepo,
		deviceStateRepo: deviceStateRepo,
		deviceTypeRepo:  deviceTypeRepo,
		locationRepo:    locationRepo,
	}
}
func (s *solarService) GetAllSolarDevices(
	ctx context.Context,
) (*[]dto.SolarDeviceView, error) {
	return s.solarRepo.GetAllSolarDevices(ctx)
}
func (s *solarService) mapDeviceToDeviceView(
	d model.Device,
	dt *model.DeviceTypes,
) dto.DeviceView {
	if dt == nil {
		dt = &model.DeviceTypes{
			Name:         "Unknown",
			HardwareType: model.HardwareTypeUnknown,
		}
	}
	return dto.DeviceView{
		ID:              d.ID,
		Name:            d.Name,
		Type:            dt.Name,
		HardwareType:    dt.HardwareType,
		Status:          d.DeviceState.Name,
		IPAddress:       d.Details.IPAddress,
		MACAddress:      d.Details.MACAddress,
		FirmwareVersion: d.Version.Name,
	}
}

func (s *solarService) CreateASolarDevice(
	ctx context.Context,
	req dto.CreateSolarDeviceDTO,
	createdBy uint,
) (dto.DeviceView, error) {
	tx := s.deviceStateRepo.BeginTx()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var userID uint

	if err := tx.Table("users").Select("id").Where("id = ?", createdBy).Scan(&userID).Error; err != nil || userID == 0 {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.DeviceView{}, ErrUserNotFound
		}
		return dto.DeviceView{}, err
	}

	deviceType, err := s.deviceTypeRepo.GetDeviceByID(
		ctx,
		req.DeviceTypeID,
	)

	if err != nil {
		tx.Rollback()
		return dto.DeviceView{}, err
	}
	hwType := deviceType.HardwareType
	if model.HardwareType(hwType) != model.HardwareTypeSolar {
		tx.Rollback()
		return dto.DeviceView{}, errors.New("invalid device type for solar device")
	}

	device := &model.Device{
		Name:         req.Name,
		DeviceTypeID: req.DeviceTypeID,
		CreatedBy:    createdBy,
	}

	initialStateID, err := s.deviceStateRepo.GetInitialDeviceStateID(
		ctx,
	)

	device.CurrentState = initialStateID

	if err := tx.Create(&device).Error; err != nil {
		tx.Rollback()
		return dto.DeviceView{}, err
	}

	// DeviceAssignment is now just LocationID and DeviceID
	var location *model.Location

	if req.LocationID != 0 {
		var err error
		location, err = s.locationRepo.GetByID(ctx, req.LocationID)
		if err != nil {
			tx.Rollback()
			return dto.DeviceView{}, errors.New("invalid location ID")
		}
	}
	if location != nil {
		assignment := &model.DeviceAssignment{
			LocationID: location.ID,
			DeviceID:   device.ID,
		}
		if err := tx.Create(assignment).Error; err != nil {
			tx.Rollback()
			return dto.DeviceView{}, err
		}
	}

	if req.ConnectedMicroControllerID != nil && *req.ConnectedMicroControllerID != 0 {
		if cd, _ := s.deviceRepo.GetDevice(ctx, *req.ConnectedMicroControllerID); cd != nil && cd.ID != 0 {
			_ = s.deviceRepo.AddConnectedDevice(
				ctx,
				device.ID,
				cd.ID,
			)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return dto.DeviceView{}, err
	}

	return s.mapDeviceToDeviceView(
		*device,
		deviceType,
	), nil
}
func (s *solarService) GetAllMySolarDevices(
	ctx context.Context,
	userID uint,
) (*[]dto.SolarDeviceView, error) {

	devices, err := s.solarRepo.GetAllMySolarDevices(ctx, userID)
	return devices, err
}
