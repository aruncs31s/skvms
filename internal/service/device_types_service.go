package service

import (
	"context"

	"github.com/aruncs31s/skvms/internal/dto"
	"github.com/aruncs31s/skvms/internal/model"
	"github.com/aruncs31s/skvms/internal/repository"
)

type DeviceTypesService interface {
	ListDeviceTypes(
		ctx context.Context,
		limit int,
		offset int,
	) ([]dto.GenericDropdown, error)
	// GetHardwareTypeByID(
	// 	ctx context.Context,
	// 	id uint,
	// ) (*dto.GenericDropdownWithFeatures, error)
	CreateDeviceType(
		ctx context.Context,
		req dto.CreateDeviceTypeRequest,
		userID uint,
	) error
	GetAllHardwareTypes(ctx context.Context) ([]dto.GenericDropdown, error)
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
func (s *deviceTypesService) GetHardwareTypeByID(
	ctx context.Context,
	id uint,
) (*dto.GenericDropdownWithFeatures, error) {
	m := make(map[string]bool)
	hwType := model.HardwareType(id)
	if hwType.CanControl() {
		m["can_control"] = true
	} else {
		m["can_control"] = false
	}
	dropdown := &dto.GenericDropdownWithFeatures{
		ID:       uint(hwType),
		Name:     hwType.String(),
		Features: m,
	}
	return dropdown, nil
}
func (s *deviceTypesService) CreateDeviceType(
	ctx context.Context,
	req dto.CreateDeviceTypeRequest,
	userID uint,
) error {

	if val, ok := model.HardwareTypeMap[model.HardwareType(req.HardwareType)]; !ok || val == "Unknown" {
		req.HardwareType = uint(model.HardwareTypeUnknown)
	}
	deviceType := &model.DeviceTypes{
		Name:         req.Name,
		HardwareType: model.HardwareType(req.HardwareType),
		CreatedBy:    userID,
		UpdatedBy:    userID,
	}
	return s.deviceTypesRepo.CreateDeviceType(ctx, deviceType)
}
func (s *deviceTypesService) GetAllHardwareTypes(ctx context.Context) ([]dto.GenericDropdown, error) {
	var dropdowns []dto.GenericDropdown
	for key, hwType := range model.HardwareTypeMap {
		dropdown := dto.GenericDropdown{
			ID:   uint(key),
			Name: hwType,
		}
		dropdowns = append(dropdowns, dropdown)
	}
	return dropdowns, nil
}
