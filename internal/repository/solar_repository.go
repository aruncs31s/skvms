package repository

import (
	"context"

	"github.com/aruncs31s/skvms/internal/dto"
	"github.com/aruncs31s/skvms/internal/logger"
	"github.com/aruncs31s/skvms/internal/model"
	"go.uber.org/zap"
)

type SolarRepository interface {
	GetAllSolarDevices(ctx context.Context) (*[]dto.SolarDeviceView, error)
	CreateASolarDevice(
		ctx context.Context,
		device *model.Device,
		createdBy uint,
	) (dto.DeviceView, error)
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
		solarDeviceView := dto.SolarDeviceView{}

		// Get Devices Address , City , IPAddress
		solarDeviceView.Address, solarDeviceView.City, solarDeviceView.ConnectedDeviceIP = r.getAddressAndIP(&device)

		// Get Connected Voltage And Current Meter
		mc, err := r.getConnectedMCs(ctx, device.ID)
		if err != nil {
			return nil, err
		}
		voltageMeter, err := r.getConnectedVoltageMeters(ctx, mc.ID)
		if err != nil {
			return nil, err
		}
		// Get Voltage And Current Readings
		readings, err := r.getVoltageAndCurrent(ctx, voltageMeter.ID)
		if err != nil {
			return nil, err
		}
		if len(*readings) > 0 {
			solarDeviceView.BatteryVoltage = (*readings)[0].Voltage
			solarDeviceView.ChargingCurrent = (*readings)[0].Current
		}

		solarDeviceView.ID = device.ID
		solarDeviceView.Name = device.Name
		// TODO: GetThe LED Status from device details or another source
		solarDeviceView.LedStatus = "Green" // Placeholder for LED status

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
	Solar uint,
) (*dto.DeviceView, error) {
	mc, err := r.deviceRepo.GetConnectedDevicesByHardwareType(ctx, Solar, model.HardwareTypeMicroController)
	if err != nil {
		logger.GetLogger().Warn(
			"Failed to get connected microcontrollers by hardware type",
			zap.String(
				"solar_id", string(Solar),
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
	voltageMeter, err := r.deviceRepo.GetConnectedDevicesByHardwareType(ctx, microcontroller, model.HardwareTypeVoltageMeter)
	if err != nil {
		logger.GetLogger().Error(
			"Failed to get connected voltage meters by hardware type",
			zap.String(
				"microcontroller_id", string(microcontroller),
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
) (*[]dto.EssentialReadingRequest, error) {
	readings, err := r.readingsRepo.ListByDevice(
		ctx,
		voltageMeterID,
		1, //return only one reading
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
	return r.mapToReadingsDTO(&readings)

}

func (r *solarRepository) mapToReadingsDTO(
	readings *[]model.Reading,
) (*[]dto.EssentialReadingRequest, error) {
	var readingsDTO []dto.EssentialReadingRequest
	for _, reading := range *readings {
		readingDTO := dto.EssentialReadingRequest{
			Voltage: reading.Voltage,
			Current: reading.Current,
		}
		readingsDTO = append(readingsDTO, readingDTO)
	}
	return &readingsDTO, nil
}
func (r *solarRepository) CreateASolarDevice(
	ctx context.Context,
	device *model.Device,
	createdBy uint,
) (dto.DeviceView, error) {
	r.deviceRepo.CreateDevice(
		ctx,
		device,
		&device.Details,
		&device.Address,
	)

	return dto.DeviceView{
		Name:         device.Name,
		HardwareType: model.HardwareTypeMap[device.DeviceType.HardwareType],
		Type:         device.DeviceType.Name,
		Status:       device.DeviceState.Name,
	}, nil
}
