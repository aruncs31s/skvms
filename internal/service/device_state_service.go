package service

import (
	"context"
	"errors"

	"github.com/aruncs31s/skvms/internal/dto"
	"github.com/aruncs31s/skvms/internal/model"
	"github.com/aruncs31s/skvms/internal/repository"
)

type DeviceStateService interface {
	DeviceControl
	DeviceStateReader
	DeviceStateWriter
}
type DeviceControl interface {
	Actuate(
		ctx context.Context,
		id uint,
		action model.DeviceAction,
		userID uint,
	) (string, error)
}
type DeviceStateReader interface {
	ListDeviceStates(
		ctx context.Context,
	) ([]dto.DeviceStateView, error)
	GetByID(
		ctx context.Context,
		id uint,
	) (*dto.DeviceStateView, error)
}

type DeviceStateWriter interface {
	Create(ctx context.Context, req *dto.CreateDeviceStateRequest) error
	Update(ctx context.Context, id uint, req *dto.UpdateDeviceStateRequest) error
	Delete(ctx context.Context, id uint) error
}
type deviceStateService struct {
	repo                      repository.DeviceStateRepository
	deviceRepo                repository.DeviceRepository
	deviceStateHistoryService DeviceStateHistoryService
}

func NewDeviceStateService(repo repository.DeviceStateRepository,
	deviceRepo repository.DeviceRepository,
	deviceStateHistoryService DeviceStateHistoryService,
) DeviceStateService {
	return &deviceStateService{
		repo:                      repo,
		deviceRepo:                deviceRepo,
		deviceStateHistoryService: deviceStateHistoryService,
	}
}

func (s *deviceStateService) Actuate(
	ctx context.Context,
	id uint,
	action model.DeviceAction,
	userID uint,
) (string, error) {
	if !action.Validate() {
		return "", errors.New("invalid device action")
	}

	tx := s.repo.BeginTx()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Lock device row
	device, err := s.deviceRepo.GetDevice(ctx, id)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	currentState := device.CurrentState
	state, err := s.repo.GetByID(
		ctx,
		currentState,
	)
	if err != nil {
		tx.Rollback()
		return "", err
	}
	// Is action allowed?
	allowedActions := model.DeviceStateTransitions[currentState]
	if !containsAction(allowedActions, action) {
		tx.Rollback()
		return state.Name, errors.New("action not allowed in current state")
	}

	// Resolve next state
	nextState, ok := model.DeviceStateActionResult[currentState][action]
	if !ok {
		tx.Rollback()
		return state.Name, errors.New("no state transition defined")
	}

	// Update device
	if err := s.repo.UpdateDeviceState(ctx, tx, device.ID, nextState); err != nil {
		tx.Rollback()
		return state.Name, err
	}

	// Log state change history
	if err := s.deviceStateHistoryService.Log(
		ctx,
		device.ID,
		action,
		nextState,
		userID,
	); err != nil {
		// Log error but don't fail the operation
		// TODO: Consider using a logger here
	}

	nextStateObj, err := s.repo.GetByID(ctx, nextState)
	if err != nil {
		tx.Rollback()
		return state.Name, err
	}

	return nextStateObj.Name, tx.Commit().Error
}
func containsAction(actions []model.DeviceAction, a model.DeviceAction) bool {
	for _, v := range actions {
		if v == a {
			return true
		}
	}
	return false
}

func (s *deviceStateService) ListDeviceStates(
	ctx context.Context,
) ([]dto.DeviceStateView, error) {
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

func (s *deviceStateService) GetByID(ctx context.Context, id uint) (*dto.DeviceStateView, error) {
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

func (s *deviceStateService) Update(ctx context.Context, id uint, req *dto.UpdateDeviceStateRequest) error {
	deviceState := &model.DeviceState{
		ID:   id,
		Name: req.Name,
	}
	return s.repo.Update(ctx, deviceState)
}

func (s *deviceStateService) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
