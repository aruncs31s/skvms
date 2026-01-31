package service

import (
	"context"

	"github.com/aruncs31s/skvms/internal/dto"
	"github.com/aruncs31s/skvms/internal/repository"
)

type DeviceTypesService interface {
	ListDeviceTypes(
		ctx context.Context,
		limit int,
		offset int,
	) ([]dto.GenericDropdown, error)
}

type deviceTypesService struct {
	deviceTypesRepo repository.DeviceTypesRepository
}

func NewDeviceTypesService(
	deviceTypesRepo repository.DeviceTypesRepository,
) DeviceTypesService {
	return &deviceTypesService{
		deviceTypesRepo: deviceTypesRepo,
	}
}

func (s *deviceTypesService) ListDeviceTypes(
	ctx context.Context,
	limit int,
	offset int,
) ([]dto.GenericDropdown, error) {
	if limit <= 0 {
		limit = 100
	}
	devices, err := s.deviceTypesRepo.ListDeviceTypes(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	var dropdowns []dto.GenericDropdown
	for _, device := range devices {
		dropdown := dto.GenericDropdown{
			ID:   device.ID,
			Name: device.Name,
		}
		dropdowns = append(dropdowns, dropdown)
	}

	return dropdowns, nil
}
