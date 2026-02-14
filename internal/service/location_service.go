package service

import (
	"context"

	"github.com/aruncs31s/skvms/internal/dto"
	"github.com/aruncs31s/skvms/internal/model"
	"github.com/aruncs31s/skvms/internal/repository"
)

type LocationService interface {
	LocationWriter
	LocationReader
}

type LocationWriter interface {
	Create(
		ctx context.Context,
		location dto.CreateLocationRequest,
	) error
	Update(
		ctx context.Context,
		id uint,
		location dto.UpdateLocationRequest,
	) error
	Delete(
		ctx context.Context,
		id uint,
	) error
}
type LocationReader interface {
	List(
		ctx context.Context,
	) ([]dto.LocationResponse, error)
	GetByID(
		ctx context.Context,
		id uint,
	) (*dto.LocationResponse, error)
	Search(
		ctx context.Context,
		query string,
	) ([]dto.LocationResponse, error)
	GetByCode(
		ctx context.Context,
		code string,
	) (*dto.LocationResponse, error)
	ListDevicesInLocation(
		ctx context.Context,
		locationID uint,
	) ([]dto.DeviceView, error)
}
type locationService struct {
	repo repository.LocationRepository
	dr   repository.DeviceRepository
}

func NewLocationService(
	repo repository.LocationRepository,
	deviceRepo repository.DeviceRepository,
) LocationService {
	return &locationService{
		repo: repo,
		dr:   deviceRepo,
	}
}

func (s *locationService) List(ctx context.Context) ([]dto.LocationResponse, error) {
	locations, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	dc := map[uint]int{}
	uc := map[uint]int{}
	for _, loc := range locations {
		c, err := s.repo.GetConnectedDevicesCount(
			ctx,
			loc.ID,
			[]model.HardwareType{
				model.HardwareTypeSolar,
			}...,
		)
		if err != nil {
			dc[loc.ID] = 0
		}
		dc[loc.ID] = c
		c, err = s.repo.GetUserCount(ctx, loc.ID)
		if err != nil {
			uc[loc.ID] = 0
		}
		uc[loc.ID] = c

	}
	var responses []dto.LocationResponse
	for _, loc := range locations {
		responses = append(responses, dto.LocationResponse{
			ID:                    loc.ID,
			Code:                  loc.Code,
			Name:                  loc.Name,
			City:                  loc.City,
			State:                 loc.State,
			PinCode:               loc.PinCode,
			ConnectedDevicesCount: dc[loc.ID],
			UserCount:             uc[loc.ID],
		})
	}
	if len(responses) == 0 {
		return []dto.LocationResponse{}, nil
	}
	return responses, nil
}

func (s *locationService) GetByID(ctx context.Context, id uint) (*dto.LocationResponse, error) {
	location, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &dto.LocationResponse{
		ID:      location.ID,
		Code:    location.Code,
		Name:    location.Name,
		State:   location.State,
		City:    location.City,
		PinCode: location.PinCode,
	}, nil
}

func (s *locationService) Search(ctx context.Context, query string) ([]dto.LocationResponse, error) {
	locations, err := s.repo.Search(ctx, query)
	if err != nil {
		return nil, err
	}
	var responses []dto.LocationResponse
	for _, loc := range locations {
		responses = append(responses, dto.LocationResponse{
			ID:   loc.ID,
			Code: loc.Code,
			Name: loc.Name,
		})
	}
	return responses, nil
}

func (s *locationService) GetByCode(ctx context.Context, code string) (*dto.LocationResponse, error) {
	location, err := s.repo.GetByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	return &dto.LocationResponse{
		ID:   location.ID,
		Code: location.Code,
		Name: location.Name,
	}, nil
}

func (s *locationService) Create(ctx context.Context, location dto.CreateLocationRequest) error {
	return s.repo.Create(ctx, &model.Location{
		Code: location.Code,
		Name: location.Name,
	})
}

func (s *locationService) Update(ctx context.Context, locationID uint, location dto.UpdateLocationRequest) error {
	existing, err := s.repo.GetByID(ctx, locationID)
	if err != nil {
		return err
	}
	existing.Code = location.Code
	existing.Name = location.Name
	existing.City = location.City
	existing.State = location.State
	existing.PinCode = location.PinCode
	return s.repo.Update(ctx, existing)
}

func (s *locationService) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
func (s *locationService) ListDevicesInLocation(ctx context.Context, locationID uint) ([]dto.DeviceView, error) {
	devices, err := s.dr.GetDevicesByLocationID(ctx, locationID)
	if err != nil {
		return nil, err
	}
	var deviceViews []dto.DeviceView
	for _, d := range devices {
		deviceViews = append(deviceViews, dto.DeviceView{
			ID:              d.ID,
			Name:            d.Name,
			Type:            d.Type,
			FirmwareVersion: d.FirmwareVersion,
			IPAddress:       d.IPAddress,
			MACAddress:      d.MACAddress,
			Status:          d.DeviceState,
		})
	}
	if len(deviceViews) == 0 {
		return []dto.DeviceView{}, nil
	}
	return deviceViews, nil
}
