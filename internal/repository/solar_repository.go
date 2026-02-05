package repository

import (
	"context"

	"github.com/aruncs31s/skvms/internal/dto"
	"github.com/aruncs31s/skvms/internal/logger"
	"github.com/aruncs31s/skvms/internal/model"
	"go.uber.org/zap"
)

type SolarRepository interface {
	GetAllSolarDevices(
		ctx context.Context,
	) (*[]dto.SolarDeviceView, error)
	CreateASolarDevice(
		ctx context.Context,
		device *model.Device,
		connectedID *uint,
		details *model.DeviceDetails,
		address *model.DeviceAddress,
	) (*model.Device, error)
	GetAllMySolarDevices(
		ctx context.Context,
		userID uint,
	) (*[]dto.SolarDeviceView, error)
}

type solarRepository struct {
	deviceRepo            DeviceRepository
	readingsRepo          ReadingRepository
	deviceStateRepository deviceStateRepository
}

func NewSolarRepository(
	deviceRepo DeviceRepository,
	readingsRepo ReadingRepository,
) *solarRepository {
	return &solarRepository{
		deviceRepo:   deviceRepo,
		readingsRepo: readingsRepo,
	}
}
func (r *solarRepository) GetAllSolarDevices(
	ctx context.Context,
) (*[]dto.SolarDeviceView, error) {
	devices, err := r.deviceRepo.GetDevicesByHardwareType(ctx, model.HardwareTypeSolar)
	if err != nil {
		return nil, err
	}

	var solarDeviceViews []dto.SolarDeviceView
	for _, device := range *devices {
		solarDeviceView, err := r.mapDeviceToSolarDeviceView(ctx, device)
		if err != nil {
			return nil, err
		}
		solarDeviceViews = append(solarDeviceViews, solarDeviceView)
	}
	return &solarDeviceViews, nil
}

func (r *solarRepository) getAddressAndIP(
	m *model.DeviceView,
	// Address , City , IPAddress
) (string, string, string) {
	if m == nil {
		return "", "", ""
	}
	return m.Address, m.City, m.IPAddress
}

func (r *solarRepository) getConnectedMCs(
	ctx context.Context,
	// Solar
	solar uint,
) (*dto.DeviceView, error) {
	mc, err := r.deviceRepo.GetConnectedDevicesByHardwareType(
		ctx,
		solar,
		model.HardwareTypeMicroController,
	)
	if err != nil {
		logger.GetLogger().Warn(
			"Failed to get connected microcontrollers by hardware type",
			zap.String(
				"solar_id", string(solar),
			),
			zap.Error(err),
			zap.Any(
				"type", model.HardwareTypeMap[model.HardwareTypeMicroController],
			),
		)
		return nil, err
	}
	return &mc, nil
}
func (r *solarRepository) getConnectedVoltageMeters(
	ctx context.Context,
	// Solar
	microcontroller uint,
) (*dto.DeviceView, error) {
	voltageMeter, err := r.deviceRepo.GetConnectedDevicesByHardwareType(
		ctx,
		microcontroller,
		model.HardwareTypeVoltageMeter,
	)
	if err != nil {
		logger.GetLogger().Error(
			"Failed to get connected voltage meters by hardware type",
			zap.Int(
				"microcontroller_id", int(microcontroller),
			),
			zap.Error(err),
			zap.Any(
				"type", model.HardwareTypeMap[model.HardwareTypeVoltageMeter],
			),
		)
		return nil, err
	}
	return &voltageMeter, nil
}

// This voltage an current is of the connected voltage meter
func (r *solarRepository) getVoltageAndCurrent(
	ctx context.Context,
	voltageMeterID uint,
) (*dto.EssentialReadingRequest, error) {
	reading, err := r.readingsRepo.GetLastReading(
		ctx,
		voltageMeterID,
	)
	if err != nil {
		logger.GetLogger().Error(
			"Failed to get voltage and current readings",
			zap.String(
				"voltage_meter_id", string(voltageMeterID),
			),
			zap.Error(err),
		)
		return nil, err
	}
	return r.mapToReadingsDTO(reading)

}

func (r *solarRepository) mapToReadingsDTO(
	reading *model.Reading,
) (*dto.EssentialReadingRequest, error) {
	return &dto.EssentialReadingRequest{
		Voltage: reading.Voltage,
		Current: reading.Current,
	}, nil
}
func (r *solarRepository) CreateASolarDevice(
	ctx context.Context,
	device *model.Device,
	connectedID *uint,
	details *model.DeviceDetails,
	address *model.DeviceAddress,
) (*model.Device, error) {
	createdDevice, err := r.deviceRepo.CreateDevice(
		ctx,
		device,
		details,
		address,
	)

	if err != nil {
		return nil, err
	}
	if connectedID != nil {
		err = r.deviceRepo.AddConnectedDevice(ctx, createdDevice.ID, *connectedID)
		if err != nil {
			logger.GetLogger().Warn("Failed to add connected device", zap.Error(err))
			// Continue, as device is created
		}
	}

	return createdDevice, nil
}

func (r *solarRepository) mapDeviceToSolarDeviceView(
	ctx context.Context,
	device model.DeviceView,
) (dto.SolarDeviceView, error) {
	solarDeviceView := dto.SolarDeviceView{}

	// Get Devices Address , City , IPAddress
	solarDeviceView.Address, solarDeviceView.City, solarDeviceView.ConnectedDeviceIP = r.getAddressAndIP(&device)

	// Get Connected Voltage And Current Meter
	mc, err := r.getConnectedMCs(ctx, device.ID)

	if err == nil && mc != nil {
		voltageMeter, err := r.getConnectedVoltageMeters(ctx, mc.ID)
		if err != nil {
			logger.GetLogger().Warn(
				"Failed to get connected voltage meter for device",
				zap.Uint("device_id", device.ID),
				zap.Error(err),
			)
			// Continue, as we can still return the device info without voltage and current readings
			solarDeviceView.BatteryVoltage = 0
			solarDeviceView.ChargingCurrent = 0
		} else {
			// Get Voltage And Current Readings
			reading, err := r.getVoltageAndCurrent(ctx, voltageMeter.ID)
			if err != nil {
				return dto.SolarDeviceView{}, err
			}
			solarDeviceView.BatteryVoltage = reading.Voltage
			solarDeviceView.ChargingCurrent = reading.Current
		}
	}

	solarDeviceView.ID = device.ID
	solarDeviceView.Name = device.Name
	// TODO: GetThe LED Status from device details or another source
	solarDeviceView.LedStatus = "Green" // Placeholder for LED status

	return solarDeviceView, nil
}

func (r *solarRepository) GetAllMySolarDevices(
	ctx context.Context,
	userID uint,
) (*[]dto.SolarDeviceView, error) {
	devices, err := r.deviceRepo.GetDevicesByHardwareTypeAndUserID(
		ctx, model.HardwareTypeSolar,
		userID,
	)

	if err != nil {
		return &[]dto.SolarDeviceView{}, err
	}
	var solarDeviceViews []dto.SolarDeviceView
	for _, device := range *devices {
		solarDeviceView, err := r.mapDeviceToSolarDeviceView(ctx, device)
		if err != nil {
			return &[]dto.SolarDeviceView{}, err
		}
		solarDeviceViews = append(solarDeviceViews, solarDeviceView)
	}
	return &solarDeviceViews, nil
}
