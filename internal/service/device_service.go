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
	CreateDevice(
		ctx context.Context,
		userID uint,
		req *dto.CreateDeviceRequest,
	) (dto.DeviceView, error)
	UpdateDevice(ctx context.Context, id uint, req *dto.UpdateDeviceRequest) error
	FullUpdateDevice(ctx context.Context, id uint, req *dto.FullUpdateDeviceRequest, updatedBy uint) error
	DeleteDevice(ctx context.Context, id uint) error
	AddConnectedDevice(ctx context.Context, parentID, childID uint) error
	GetConnectedDevices(ctx context.Context, parentID uint) ([]dto.DeviceView, error)
	SearchDevices(ctx context.Context, query string) ([]dto.GenericDropdown, error)
	SearchMicrocontrollers(ctx context.Context, query string) ([]dto.GenericDropdown, error)
	SearchSensors(ctx context.Context, query string) ([]dto.GenericDropdown, error)
	ListAllSensors(ctx context.Context) ([]dto.DeviceView, error)
	GetSensorDevice(ctx context.Context, id uint) (*dto.DeviceView, error)
	CreateSensorDevice(ctx context.Context, userID uint, req *dto.CreateDeviceRequest) (dto.DeviceView, error)
}

type deviceService struct {
	repo             repository.DeviceRepository
	userRepo         repository.UserRepository
	auditService     AuditService
	stateMgmtService DeviceStateService
}

func NewDeviceService(
	repo repository.DeviceRepository,
	userRepo repository.UserRepository,
	stateMgmtService DeviceStateService,
	auditService AuditService,

) DeviceService {
	return &deviceService{
		repo:             repo,
		userRepo:         userRepo,
		stateMgmtService: stateMgmtService,
		auditService:     auditService,
	}
}

func (s *deviceService) mapDeviceToDeviceView(d model.Device) dto.DeviceView {
	return dto.DeviceView{
		ID:              d.ID,
		Name:            d.Name,
		Type:            d.DeviceType.Name,
		HardwareType:    d.DeviceType.HardwareType.String(),
		Status:          d.DeviceState.Name,
		IPAddress:       d.Details.IPAddress,
		MACAddress:      d.Details.MACAddress,
		FirmwareVersion: d.Version.Name,
		Address:         d.Address.Address,
		City:            d.Address.City,
	}
}

func (s *deviceService) ListDevices(ctx context.Context) ([]dto.DeviceView, error) {
	devices, err := s.repo.ListDevices(ctx)
	if err != nil || len(devices) == 0 {
		return []dto.DeviceView{}, err
	}
	var dtos []dto.DeviceView
	for _, device := range devices {
		dtos = append(dtos, s.mapToDeviceModelViewToView(device))
	}
	return dtos, nil
}
func (s *deviceService) mapToDeviceModelViewToView(d model.DeviceView) dto.DeviceView {
	return dto.DeviceView{
		ID:              d.ID,
		Name:            d.Name,
		Type:            d.Type,
		HardwareType:    d.HardwareType.String(),
		IPAddress:       d.IPAddress,
		MACAddress:      d.MACAddress,
		FirmwareVersion: d.FirmwareVersion,
		Address:         d.Address,
		City:            d.City,
		// Their State is the status
		Status: d.DeviceState,
	}
}

func (s *deviceService) ListDevicesByUser(ctx context.Context, userID uint) ([]dto.DeviceView, error) {
	devices, err := s.repo.ListDevicesByUser(ctx, userID)
	if err != nil || len(devices) == 0 {
		return []dto.DeviceView{}, err
	}
	return devices, nil
}

func (s *deviceService) GetDevice(ctx context.Context, id uint) (*dto.DeviceView, error) {
	device, err := s.repo.GetDevice(ctx, id)
	if device == nil || err != nil {
		return nil, err
	}

	dv := s.mapDeviceToDeviceView(*device)
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

func (s *deviceService) CreateDevice(
	ctx context.Context,
	userID uint,
	req *dto.CreateDeviceRequest,
) (dto.DeviceView, error) {

	if userID == 0 {
		return dto.DeviceView{}, fmt.Errorf("invalid user id")
	}
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return dto.DeviceView{}, err
	}
	if user == nil {
		return dto.DeviceView{}, fmt.Errorf("user not found")
	}

	var v *model.Version

	if req.FirmwareVersionID != 0 {
		var err error
		v, err = s.repo.FindVersionByID(ctx, req.FirmwareVersionID)
		if err != nil {
			return dto.DeviceView{}, err
		}
		if v == nil {
			return dto.DeviceView{}, fmt.Errorf("firmware version not found")
		}
	}

	if req.FirmwareVersionID == 0 {
		req.FirmwareVersionID = 1 // Default to V1.0.0 if not provided
	}

	device := &model.Device{
		Name:         req.Name,
		DeviceTypeID: req.Type,
		VersionID:    &req.FirmwareVersionID,
		CurrentState: 1, // Default to Active
		CreatedBy:    userID,
		UpdatedBy:    userID,
	}
	details := &model.DeviceDetails{
		IPAddress:  req.IPAddress,
		MACAddress: req.MACAddress,
	}

	address := &model.DeviceAddress{
		Address: req.Address,
		City:    req.City,
	}
	newDevice, err := s.repo.CreateDevice(ctx, device, details, address)
	if err != nil {
		return dto.DeviceView{}, err
	}

	loadedDevice, err := s.repo.GetDevice(ctx, newDevice.ID)
	if err != nil {
		return dto.DeviceView{}, err
	}
	if loadedDevice == nil {
		return dto.DeviceView{}, fmt.Errorf("failed to retrieve created device")
	}

	return s.mapDeviceToDeviceView(*loadedDevice), nil
}

func (s *deviceService) UpdateDevice(
	ctx context.Context,
	id uint,
	req *dto.UpdateDeviceRequest,
) error {
	existing, err := s.repo.GetDevice(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("device not found")
	}

	// Update device fields if provided
	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.Type != nil {
		existing.DeviceTypeID = *req.Type
	}

	// Update details fields if provided
	if req.IPAddress != nil {
		existing.Details.IPAddress = *req.IPAddress
	}
	if req.MACAddress != nil {
		existing.Details.MACAddress = *req.MACAddress
	}

	// Update firmware version if provided
	if req.FirmwareVersionID != nil {
		existing.VersionID = req.FirmwareVersionID
	}

	// Update address fields if provided
	if req.Address != nil {
		existing.Address.Address = *req.Address
	}
	if req.City != nil {
		existing.Address.City = *req.City
	}

	return s.repo.UpdateDevice(ctx, existing, &existing.Details, &existing.Address)
}

func (s *deviceService) FullUpdateDevice(
	ctx context.Context,
	id uint,
	req *dto.FullUpdateDeviceRequest,
	updatedBy uint,
) error {
	existing, err := s.repo.GetDevice(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("device not found")
	}

	// Update all device fields
	existing.Name = req.Name
	existing.DeviceTypeID = req.Type
	existing.CurrentState = req.CurrentState
	existing.UpdatedBy = updatedBy

	// Update details
	existing.Details.IPAddress = req.IPAddress
	existing.Details.MACAddress = req.MACAddress

	// Update firmware version

	existing.VersionID = &req.FirmwareVersionID

	// Update address
	existing.Address.Address = req.Address
	existing.Address.City = req.City

	return s.repo.UpdateDevice(ctx, existing, &existing.Details, &existing.Address)
}

func (s *deviceService) DeleteDevice(ctx context.Context, id uint) error {
	return s.repo.DeleteDevice(ctx, id)
}

func (s *deviceService) AddConnectedDevice(ctx context.Context, parentID, childID uint) error {
	// Check if parent device exists
	parent, err := s.repo.GetDevice(ctx, parentID)
	if err != nil {
		return err
	}
	if parent == nil {
		return fmt.Errorf("parent device not found")
	}

	// Check if child device exists
	child, err := s.repo.GetDevice(ctx, childID)
	if err != nil {
		return err
	}
	if child == nil {
		return fmt.Errorf("child device not found")
	}

	// Prevent self-connection
	if parentID == childID {
		return fmt.Errorf("cannot connect device to itself")
	}

	return s.repo.AddConnectedDevice(ctx, parentID, childID)
}

func (s *deviceService) GetConnectedDevices(ctx context.Context, parentID uint) ([]dto.DeviceView, error) {
	// Check if parent device exists
	parent, err := s.repo.GetDevice(ctx, parentID)
	if err != nil {
		return nil, err
	}
	if parent == nil {
		return nil, fmt.Errorf("parent device not found")
	}

	return s.repo.GetConnectedDevices(ctx, parentID)
}

func (s *deviceService) SearchDevices(ctx context.Context, query string) ([]dto.GenericDropdown, error) {
	return s.repo.SearchDevices(ctx, query)
}

func (s *deviceService) SearchMicrocontrollers(ctx context.Context, query string) ([]dto.GenericDropdown, error) {
	devices, err := s.repo.SearchDevicesByHardwareType(
		ctx,
		query,
		model.HardwareTypeMicroController,
	)
	if err != nil || len(devices) == 0 {
		return []dto.GenericDropdown{}, nil
	}
	return devices, nil
}
func (s *deviceService) SearchSensors(ctx context.Context, query string) ([]dto.GenericDropdown, error) {
	devices, err := s.repo.SearchDevicesByHardwareTypes(
		ctx,
		query,
		[]model.HardwareType{
			model.HardwareTypeSensor,
			model.HardwareTypeCurrentSensor,
			model.HardwareTypePowerMeter,
			model.HardwareTypeVoltageMeter,
		},
	)
	if err != nil || len(devices) == 0 {
		return []dto.GenericDropdown{}, nil
	}
	return devices, nil
}

func (s *deviceService) ListAllSensors(ctx context.Context) ([]dto.DeviceView, error) {
	devices, err := s.repo.ListDevicesByHardwareTypes(
		ctx,
		[]model.HardwareType{
			model.HardwareTypeSensor,
			model.HardwareTypeCurrentSensor,
			model.HardwareTypePowerMeter,
			model.HardwareTypeVoltageMeter,
		},
	)
	if err != nil || len(devices) == 0 {
		return []dto.DeviceView{}, nil
	}
	var dtos []dto.DeviceView
	for _, device := range devices {
		dtos = append(dtos, s.mapToDeviceModelViewToView(device))
	}
	return dtos, nil
}

func (s *deviceService) GetSensorDevice(ctx context.Context, id uint) (*dto.DeviceView, error) {
	device, err := s.repo.GetDevice(ctx, id)
	if err != nil {
		return nil, err
	}
	// Check if it's a sensor type
	sensorTypes := []model.HardwareType{
		model.HardwareTypeSensor,
		model.HardwareTypeCurrentSensor,
		model.HardwareTypePowerMeter,
		model.HardwareTypeVoltageMeter,
	}
	isSensor := false
	for _, t := range sensorTypes {
		if device.DeviceType.HardwareType == t {
			isSensor = true
			break
		}
	}
	if !isSensor {
		return nil, fmt.Errorf("device is not a sensor")
	}
	dto := s.mapDeviceToDeviceView(*device)
	return &dto, nil
}

func (s *deviceService) CreateSensorDevice(ctx context.Context, userID uint, req *dto.CreateDeviceRequest) (dto.DeviceView, error) {
	// For now, just create the device, assuming the type is sensor
	// In future, could validate the device type is sensor
	return s.CreateDevice(ctx, userID, req)
}
