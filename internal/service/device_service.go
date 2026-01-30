package service

import (
	"context"
	"fmt"

	"github.com/aruncs31s/skvms/internal/dto"
	"github.com/aruncs31s/skvms/internal/repository"
)

type DeviceService interface {
	ListDevices(ctx context.Context) ([]dto.DeviceView, error)
	GetDevice(ctx context.Context, id uint) (*dto.DeviceView, error)
	ControlDevice(ctx context.Context, id uint, command string) (string, error)
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
