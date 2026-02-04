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
}
type solarService struct {
	solarRepo       repository.SolarRepository
	userRepo        repository.UserReader
	deviceRepo      repository.DeviceRepository
	deviceStateRepo repository.DeviceStateRepository
}

func NewSolarService(
	solarRepo repository.SolarRepository,
	user repository.UserReader,
	deviceRepo repository.DeviceRepository,
	deviceStateRepo repository.DeviceStateRepository,
) SolarService {
	return &solarService{
		solarRepo:       solarRepo,
		userRepo:        user,
		deviceRepo:      deviceRepo,
		deviceStateRepo: deviceStateRepo,
	}
}
func (s *solarService) GetAllSolarDevices(
	ctx context.Context,
) (*[]dto.SolarDeviceView, error) {
	return s.solarRepo.GetAllSolarDevices(ctx)
}
func (s *solarService) CreateASolarDevice(
	ctx context.Context,
	req dto.CreateSolarDeviceDTO,
	createdBy uint,

) (dto.DeviceView, error) {

	exists, err := s.userRepo.CheckIfExistsByUserID(ctx, createdBy)
	if err != nil {
		return dto.DeviceView{}, err
	}
	if !exists {
		return dto.DeviceView{}, ErrUserNotFound
	}

	device := &model.Device{
		Name:         req.Name,
		DeviceTypeID: req.DeviceTypeID,
	}
	connectedDevice, err := s.deviceRepo.GetDevice(
		ctx, *req.ConnectedMicroControllerID,
	)
	if err != nil {
		logger.GetLogger().Warn(
			"Connected Device Error",
			zap.Error(err),
		)
	}
	if connectedDevice != nil && connectedDevice.ID != 0 {
		device.ConnectedDevices = []model.ConnectedDevice{
			{
				ParentID: *req.ConnectedMicroControllerID,
				ChildID:  connectedDevice.ID,
			},
		}
	}

	device.Address = model.DeviceAddress{
		Address: req.Address,
		City:    req.City,
	}

	device.Version = model.Version{
		Name: "V1.0.0",
	}
	// Get Initial Device State ID
	initialStateID, _ := s.deviceStateRepo.GetInitialDeviceStateID(ctx)

	device.CurrentState = initialStateID
	solarDevice, err := s.solarRepo.CreateASolarDevice(ctx, device, createdBy)
	return solarDevice, err
}
