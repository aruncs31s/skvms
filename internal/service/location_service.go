package service

import (
	"context"
	"sync"

	"github.com/aruncs31s/skvms/internal/dto"
	"github.com/aruncs31s/skvms/internal/model"
	"github.com/aruncs31s/skvms/internal/repository"
	"golang.org/x/sync/errgroup"
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
	) ([]dto.GenericDropdown, error)
	GetByCode(
		ctx context.Context,
		code string,
	) (*dto.LocationResponse, error)
	ListDevicesInLocation(
		ctx context.Context,
		locationID uint,
	) ([]dto.DeviceView, error)
	SevenDaysReadings(
		ctx context.Context,
		locationID uint,
	) ([]dto.ReadingsResponse, error)
}
type locationService struct {
	repo repository.LocationRepository
	dr   repository.DeviceRepository
	rr   repository.ReadingRepository
}

func NewLocationService(
	repo repository.LocationRepository,
	deviceRepo repository.DeviceRepository,
	readingRepo repository.ReadingRepository,
) LocationService {
	return &locationService{
		repo: repo,
		dr:   deviceRepo,
		rr:   readingRepo,
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
			ID:          loc.ID,
			Code:        loc.Code,
			Name:        loc.Name,
			Description: loc.Description,
			City:        loc.City,
			State:       loc.State,

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

func (s *locationService) GetByID(
	ctx context.Context,
	id uint,
) (*dto.LocationResponse, error) {
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

func (s *locationService) Search(ctx context.Context, query string) ([]dto.GenericDropdown, error) {
	locations, err := s.repo.Search(ctx, query)
	if err != nil {
		return nil, err
	}
	var responses []dto.GenericDropdown
	for _, loc := range locations {
		responses = append(responses, dto.GenericDropdown{
			ID:   loc.ID,
			Name: loc.Name + " - (" + loc.Code + ")",
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

func (s *locationService) Create(
	ctx context.Context,
	location dto.CreateLocationRequest,
) error {
	err := s.repo.Create(ctx, &model.Location{
		Code:        location.Code,
		Name:        location.Name,
		Description: location.Description,
		State:       location.State,
		City:        location.City,
		PinCode:     location.PinCode,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *locationService) Update(
	ctx context.Context,
	locationID uint,
	location dto.UpdateLocationRequest,
) error {
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

func (s *locationService) Delete(
	ctx context.Context,
	id uint,
) error {
	return s.repo.Delete(ctx, id)
}
func (s *locationService) ListDevicesInLocation(
	ctx context.Context,
	locationID uint,
) ([]dto.DeviceView, error) {
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

func (s *locationService) SevenDaysReadings(
	ctx context.Context,
	locationID uint,
) ([]dto.ReadingsResponse, error) {

	if locationID == 0 {
		return []dto.ReadingsResponse{}, nil
	}

	deviceInThatLocation, err := s.dr.GetDevicesByLocationID(
		ctx,
		locationID,
	)
	if err != nil {
		return []dto.ReadingsResponse{}, err
	}
	if len(deviceInThatLocation) == 0 {
		return []dto.ReadingsResponse{}, nil
	}
	if len(deviceInThatLocation) == 1 {
		responses, err := s.sevenDaysReadingFromSingleMicroController(ctx, locationID, deviceInThatLocation[0].ID)
		if err != nil {
			return nil, err
		}
		return responses, nil
	}
	var deviceIDs []uint
	for _, d := range deviceInThatLocation {
		deviceIDs = append(deviceIDs, d.ID)
	}
	return s.sevenDaysReadingFromSingleMicroControllers(ctx, locationID, deviceIDs)
}

func (s *locationService) sevenDaysReadingFromSingleMicroControllers(
	ctx context.Context,
	locationID uint,
	deviceID []uint,
) ([]dto.ReadingsResponse, error) {
	if len(deviceID) == 0 {
		return []dto.ReadingsResponse{}, nil
	}

	g, gCtx := errgroup.WithContext(ctx)
	var mu sync.Mutex
	responsesPerDevice := map[uint][]dto.ReadingsResponse{}

	for _, dID := range deviceID {
		dID := dID
		g.Go(func() error {
			responses, err := s.sevenDaysReadingFromSingleMicroController(gCtx, locationID, dID)
			if err != nil {
				return err
			}
			mu.Lock()
			responsesPerDevice[dID] = responses
			mu.Unlock()
			return nil
		})
	}

	err := g.Wait()
	if err != nil {
		return []dto.ReadingsResponse{}, err
	}

	var finalResponses []dto.ReadingsResponse
	for _, res := range responsesPerDevice {
		finalResponses = append(finalResponses, res...)
	}
	if len(finalResponses) == 0 {
		return []dto.ReadingsResponse{}, nil
	}
	return finalResponses, nil
}
func (s *locationService) sevenDaysReadingFromSingleMicroController(
	ctx context.Context,
	locationID uint,
	deviceID uint,
) ([]dto.ReadingsResponse, error) {
	readings, err := s.rr.SevenDaysReadingsByLocation(ctx, locationID, deviceID)
	if err != nil {
		return nil, err
	}
	var responses []dto.ReadingsResponse
	for _, r := range readings {
		responses = append(responses, dto.ReadingsResponse{
			Voltage:   r.Voltage,
			Current:   r.Current,
			CreatedAt: r.Bucket,
		})
	}
	if len(responses) == 0 {
		return []dto.ReadingsResponse{}, nil
	}
	return responses, nil
}
