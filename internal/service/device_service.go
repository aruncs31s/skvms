package service

import (
	"context"
	"fmt"

	"github.com/aruncs31s/skvms/internal/dto"
	"github.com/aruncs31s/skvms/internal/model"
	"github.com/aruncs31s/skvms/internal/repository"
)

type DeviceService interface {
	ListDevices(ctx context.Context) ([]dto.DeviceView, error)
	GetDevice(ctx context.Context, id uint) (*dto.DeviceView, error)
	ControlDevice(ctx context.Context, id uint, command string) (string, error)
	CreateDevice(ctx context.Context, req *dto.CreateDeviceRequest) error
	UpdateDevice(ctx context.Context, id uint, req *dto.UpdateDeviceRequest) error
	DeleteDevice(ctx context.Context, id uint) error
}

type deviceService struct {
	repo repository.DeviceRepository
}

func NewDeviceService(repo repository.DeviceRepository) DeviceService {
	return &deviceService{repo: repo}
}

func (s *deviceService) ListDevices(ctx context.Context) ([]dto.DeviceView, error) {
	return s.repo.ListDevices(ctx)
}

func (s *deviceService) GetDevice(ctx context.Context, id uint) (*dto.DeviceView, error) {
	return s.repo.GetDevice(ctx, id)
}

func (s *deviceService) ControlDevice(ctx context.Context, id uint, command string) (string, error) {
	device, err := s.repo.GetDevice(ctx, id)
	if err != nil {
		return "", err
	}
	if device == nil {
		return "", nil
	}
	if command == "" {
		command = "connect"
	}
	return fmt.Sprintf("Device %s (%d) command accepted: %s", device.Name, device.ID, command), nil
}

func (s *deviceService) CreateDevice(ctx context.Context, req *dto.CreateDeviceRequest) error {
	device := &model.Device{
		Name:      req.Name,
		Type:      req.Type,
		VersionID: req.FirmwareVersionID,
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
		ID:   id,
		Name: req.Name,
		Type: req.Type,
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
