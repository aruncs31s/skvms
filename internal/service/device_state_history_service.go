package service

import (
	"context"
	"fmt"

	"github.com/aruncs31s/skvms/internal/dto"
	"github.com/aruncs31s/skvms/internal/model"
	"github.com/aruncs31s/skvms/internal/repository"
)

type DeviceStateHistoryService interface {
	Log(
		ctx context.Context,
		deviceID uint,
		action model.DeviceAction,
		newStateID uint,
	) error
	GetDeviceStateHistory(
		ctx context.Context,
		deviceStateRequest dto.DeviceStateFilterRequest,
	) (dto.DeviceStateHistoryViewResponse, error)
}
type deviceStateHistoryService struct {
	repo repository.DeviceStateHistoryRepository
}

func NewDeviceStateHistoryService(
	repo repository.DeviceStateHistoryRepository,
) DeviceStateHistoryService {
	return &deviceStateHistoryService{}
}
func (s *deviceStateHistoryService) Log(
	ctx context.Context,
	deviceID uint,
	action model.DeviceAction,
	newStateID uint,
) error {
	history := &model.DeviceStateHistory{
		DeviceID:     deviceID,
		CausedAction: action,
		StateID:      newStateID,
	}
	return s.repo.Create(ctx, nil, history)
}
func (s *deviceStateHistoryService) GetDeviceStateHistory(
	ctx context.Context,
	deviceStateRequest dto.DeviceStateFilterRequest,
) (dto.DeviceStateHistoryViewResponse, error) {

	records, err := s.repo.GetDeviceStateHistory(
		ctx,
		deviceStateRequest.DeviceID,
		deviceStateRequest.States,
		deviceStateRequest.FromDate,
		deviceStateRequest.ToDate,
	)
	if err != nil {
		return dto.DeviceStateHistoryViewResponse{}, err
	}

	views, err := s.mapToDTO(ctx, records)
	if err != nil {
		return dto.DeviceStateHistoryViewResponse{}, err
	}
	return dto.DeviceStateHistoryViewResponse{
		TotalRecords: 0,
		History:      views,
	}, nil
}

func (s *deviceStateHistoryService) mapToDTO(
	ctx context.Context,
	histories []model.DeviceStateHistoryReport,
) ([]dto.DeviceStateHistoryView, error) {
	views := make([]dto.DeviceStateHistoryView, len(histories))
	for i, history := range histories {
		action, ok := model.DeviceActionsMap[history.ActionCaused]
		if !ok {
			return nil, fmt.Errorf("invalid action code: %d", history.ActionCaused)
		}
		views[i] = dto.DeviceStateHistoryView{
			StateName:    history.StateName,
			ActionCaused: action,
			ChangedAt:    history.ChangedAt,
			ChangedBy:    history.ChangedBy,
		}
	}
	return views, nil
}
