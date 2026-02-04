package repository

import (
	"context"

	"github.com/aruncs31s/skvms/internal/dto"
	"github.com/aruncs31s/skvms/internal/model"
)

type solarRepository struct {
	solarDevice SolarReader
}

func NewSolarRepository(solarDevice SolarReader) *solarRepository {
	return &solarRepository{
		solarDevice: solarDevice,
	}
}
func (r *solarRepository) GetAllSolarDevices(
	ctx context.Context,
) (*[]dto.SolarDeviceView, error) {
	devices, err := r.solarDevice.GetDevicesByHardwareType(ctx, model.HardwareTypeSolar)
	if err != nil {
		return nil, err
	}
	var solarDeviceViews []dto.SolarDeviceView
	for _, device := range *devices {
		solarDeviceView := dto.SolarDeviceView{}

		// Get Devices Address , City , IPAddress
		solarDeviceView.Address, solarDeviceView.City, solarDeviceView.ConnectedDeviceIP = r.getAddressAndIP(&device)

		// Get Connected Voltage And Current Meter
	}
}
func (r *solarRepository) getAddressAndIP(
	models *[]model.DeviceView,
	// Address , City , IPAddress
) (string, string, string)

func (r *solarRepository) getConnectedMCs(
	ctx context.Context,
	// Solar
	Solar uint,
) {
	mc := r.solarDevice.GetConnectedDevicesByHardwareType(ctx, Solar, model.HardwareTypeMicroController)
}
