package service

import (
	"context"

	"github.com/aruncs31s/skvms/internal/dto"
	"github.com/aruncs31s/skvms/internal/model"
	"github.com/aruncs31s/skvms/internal/repository"
)

type DeviceStateService interface {
	ListDeviceStates(
		ctx context.Context,
		) ([]dto.DeviceStateView, error)
	GetByID(
		ctx context.Context,
		id int,
		) (*dto.DeviceStateView, error)
	Create(ctx context.Context, req *dto.CreateDeviceStateRequest) error
	Update(ctx context.Context, id int, req *dto.UpdateDeviceStateRequest) error
	Delete(ctx context.Context, id int) error
}

type deviceStateService struct {
	repo repository.DeviceStateRepository
}

func NewDeviceStateService(repo repository.DeviceStateRepository) DeviceStateService {
	return &deviceStateService{repo: repo}
}

func (s *deviceStateService) ListDeviceStates(ctx context.Context) ([]dto.DeviceStateView, error) {
	deviceStates, err := s.repo.ListDeviceStates(ctx)
	if err != nil {
		return nil, err
	}

	views := make([]dto.DeviceStateView, len(deviceStates))
	for i, deviceState := range deviceStates {
		views[i] = dto.DeviceStateView{
			ID:        deviceState.ID,
			Name:      deviceState.Name,
			CreatedAt: deviceState.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}
	return views, nil
}

func (s *deviceStateService) GetByID(ctx context.Context, id int) (*dto.DeviceStateView, error) {
	deviceState, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if deviceState == nil {
		return nil, nil
	}

	view := &dto.DeviceStateView{
		ID:        deviceState.ID,
		Name:      deviceState.Name,
		CreatedAt: deviceState.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	return view, nil
}

func (s *deviceStateService) Create(ctx context.Context, req *dto.CreateDeviceStateRequest) error {
	deviceState := &model.DeviceState{
		Name: req.Name,
	}
	return s.repo.Create(ctx, deviceState)
}

func (s *deviceStateService) Update(ctx context.Context, id int, req *dto.UpdateDeviceStateRequest) error {
	deviceState := &model.DeviceState{
		ID:   id,
		Name: req.Name,
	}
	return s.repo.Update(ctx, deviceState)
}

func (s *deviceStateService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
	return &deviceStateService{
		// Initialize dependencies
	}
}
func (s *)