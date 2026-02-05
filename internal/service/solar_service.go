package service

import (
	"context"
	"errors"

	"github.com/aruncs31s/skvms/internal/dto"
	"github.com/aruncs31s/skvms/internal/logger"
	"github.com/aruncs31s/skvms/internal/model"
	"github.com/aruncs31s/skvms/internal/repository"
	"go.uber.org/zap"
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
	solarRepo       repository.SolarRepository
	userRepo        repository.UserReader
	deviceRepo      repository.DeviceRepository
	deviceTypeRepo  repository.DeviceTypesRepository
	deviceStateRepo repository.DeviceStateRepository
}

func NewSolarService(
	solarRepo repository.SolarRepository,
	user repository.UserReader,
	deviceRepo repository.DeviceRepository,
	deviceStateRepo repository.DeviceStateRepository,
	deviceTypeRepo repository.DeviceTypesRepository,
) SolarService {
	return &solarService{
		solarRepo:       solarRepo,
		userRepo:        user,
		deviceRepo:      deviceRepo,
		deviceStateRepo: deviceStateRepo,
		deviceTypeRepo:  deviceTypeRepo,
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
		Address:         d.Address.Address,
		City:            d.Address.City,
	}
}
func (s *solarService) CreateASolarDevice(
	ctx context.Context,
	req dto.CreateSolarDeviceDTO,
	createdBy uint,

) (dto.DeviceView, error) {

	exists, err := s.userRepo.CheckIfExistsByUserID(ctx, createdBy)

	if err != nil {
		return dto.DeviceView{}, ErrUserNotFound
	}
	if !exists {
		return dto.DeviceView{}, ErrUserNotFound
	}

	deviceType, err := s.deviceTypeRepo.GetDeviceByID(
		ctx,
		req.DeviceTypeID,
	)
	if err != nil {
		return dto.DeviceView{}, err
	}
	hwType := deviceType.HardwareType
	if model.HardwareType(hwType) != model.HardwareTypeSolar {
		return dto.DeviceView{}, errors.New("invalid device type for solar device")
	}

	device := &model.Device{
		Name:         req.Name,
		DeviceTypeID: req.DeviceTypeID,
		CreatedBy:    createdBy,
	}
	details := model.DeviceDetails{} // Empty details for solar devices

	address := model.DeviceAddress{
		Address: req.Address,
		City:    req.City,
	}

	// Get Initial Device State ID
	initialStateID, err := s.deviceStateRepo.GetInitialDeviceStateID(
		ctx,
	)

	if err != nil {
		logger.GetLogger().Error("Failed to get initial device state", zap.Error(err))
		return dto.DeviceView{}, err
	}

	device.CurrentState = initialStateID

	createdDevice, err := s.solarRepo.CreateASolarDevice(
		ctx,
		device,
		req.ConnectedMicroControllerID,
		&details,
		&address,
	)
	if err != nil {
		return dto.DeviceView{}, err
	}

	return s.mapDeviceToDeviceView(
		*createdDevice,
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
