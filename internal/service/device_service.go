package service

import (
	"context"
	"fmt"

	"github.com/aruncs31s/skvms/internal/dto"
	"github.com/aruncs31s/skvms/internal/model"
	"github.com/aruncs31s/skvms/internal/repository"
)

type DeviceService interface {
	ListDevices(
		ctx context.Context,
	) ([]dto.DeviceView, error)
	ListDevicesByUser(
		ctx context.Context,
		userID uint,
	) ([]dto.DeviceView, error)
	GetDevice(
		ctx context.Context,
		id uint,
	) (*dto.DeviceView, error)
	ControlDevice(
		ctx context.Context,
		id uint,
		action uint,
		userID uint,
	) (dto.DeviceControlResponse, error)
	CreateDevice(ctx context.Context, req *dto.CreateDeviceRequest) error
	UpdateDevice(ctx context.Context, id uint, req *dto.UpdateDeviceRequest) error
	DeleteDevice(ctx context.Context, id uint) error
}

type deviceService struct {
	repo             repository.DeviceRepository
	auditService     AuditService
	stateMgmtService DeviceStateService
}

func NewDeviceService(
	repo repository.DeviceRepository,
	stateMgmtService DeviceStateService,
	auditService AuditService,
) DeviceService {
	return &deviceService{
		repo:             repo,
		stateMgmtService: stateMgmtService,
		auditService:     auditService,
	}
}

func (s *deviceService) ListDevices(ctx context.Context) ([]dto.DeviceView, error) {
	return s.repo.ListDevices(ctx)
}

func (s *deviceService) ListDevicesByUser(ctx context.Context, userID uint) ([]dto.DeviceView, error) {
	return s.repo.ListDevicesByUser(ctx, userID)
}

func (s *deviceService) GetDevice(ctx context.Context, id uint) (*dto.DeviceView, error) {
	device, err := s.repo.GetDevice(ctx, id)
	if device == nil || err != nil {
		return nil, err
	}

	dv := dto.DeviceView{
		ID:              device.ID,
		Name:            device.Name,
		Type:            device.DeviceType.Name,
		FirmwareVersion: device.Details.FirmwareVersion,
		IPAddress:       device.Details.IPAddress,
		MACAddress:      device.Details.MACAddress,
		Address:         device.Address.Address,
		City:            device.Address.City,
		DeviceState:     device.CurrentState,
	}
	return &dv, nil
}

// ControlDevice means , changing a state like turning on/off a selected device.-
func (s *deviceService) ControlDevice(
	ctx context.Context,
	id uint,
	action uint,
	userID uint,
) (dto.DeviceControlResponse, error) {
	device, err := s.repo.GetDevice(ctx, id)
	if err != nil {
		return dto.DeviceControlResponse{}, err
	}
	if device == nil {
		return dto.DeviceControlResponse{}, nil
	}

	requestedAction := model.DeviceAction(action)

	if !requestedAction.Validate() {
		return dto.DeviceControlResponse{}, fmt.Errorf("invalid action")
	}

	// Check if the action is allowed for the current state
	allowedActions, stateExists := model.DeviceStateTransitions[device.CurrentState]
	if !stateExists {
		return dto.DeviceControlResponse{}, fmt.Errorf("unknown device state: %d", device.CurrentState)
	}

	// Check if the requested action is in the allowed actions for this state
	actionAllowed := false
	for _, allowedAction := range allowedActions {
		if allowedAction == requestedAction {
			actionAllowed = true
			break
		}
	}

	if !actionAllowed {
		actionName, _ := model.DeviceActionsMap[requestedAction]
		return dto.DeviceControlResponse{}, fmt.Errorf("action '%s' is not allowed for current state (state ID: %d)", actionName, device.CurrentState)
	}

	// Update the device state based on the action
	newState, err := s.stateMgmtService.Actuate(
		ctx,
		device.ID,
		requestedAction,
		userID,
	)
	if err != nil {
		return dto.DeviceControlResponse{}, err
	}

	return dto.DeviceControlResponse{
		Device: device.Name,
		State:  newState,
	}, nil
}

func (s *deviceService) CreateDevice(ctx context.Context, req *dto.CreateDeviceRequest) error {
	device := &model.Device{
		Name:         req.Name,
		DeviceTypeID: req.Type,
		VersionID:    req.FirmwareVersionID,
	}
	details := &model.DeviceDetails{
		IPAddress:       req.IPAddress,
		MACAddress:      req.MACAddress,
		FirmwareVersion: "",
	}
	// Set FirmwareVersion from version
	if req.FirmwareVersionID != 0 {
		version, err := s.repo.FindVersionByID(ctx, req.FirmwareVersionID)
		if err != nil {
			return err
		}
		details.FirmwareVersion = version.Version
	}
	address := &model.DeviceAddress{
		Address: req.Address,
		City:    req.City,
	}
	return s.repo.CreateDevice(ctx, device, details, address)
}

func (s *deviceService) UpdateDevice(ctx context.Context, id uint, req *dto.UpdateDeviceRequest) error {
	device := &model.Device{
		ID:           id,
		Name:         req.Name,
		DeviceTypeID: req.Type,
	}
	details := &model.DeviceDetails{
		DeviceID:        id,
		IPAddress:       req.IPAddress,
		MACAddress:      req.MACAddress,
		FirmwareVersion: "",
	}
	// Set VersionID and FirmwareVersion
	if req.FirmwareVersionID != 0 {
		version, err := s.repo.FindVersionByID(ctx, req.FirmwareVersionID)
		if err != nil {
			return err
		}
		device.VersionID = req.FirmwareVersionID
		details.FirmwareVersion = version.Version
	}
	address := &model.DeviceAddress{
		DeviceID: id,
		Address:  req.Address,
		City:     req.City,
	}
	return s.repo.UpdateDevice(ctx, device, details, address)
}

func (s *deviceService) DeleteDevice(ctx context.Context, id uint) error {
	return s.repo.DeleteDevice(ctx, id)
}
