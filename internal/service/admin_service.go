package service

import (
	"context"

	"github.com/aruncs31s/skvms/internal/dto"
	"github.com/aruncs31s/skvms/internal/repository"
)

type AdminService interface {
	GetStats(ctx context.Context) (*dto.AdminAggregateResponse, error)
}

type adminService struct {
	userRepo     repository.UserRepository
	deviceRepo   repository.DeviceRepository
	readingRepo  repository.ReadingRepository
	auditRepo    repository.AuditRepository
}

func NewAdminService(
	userRepo repository.UserRepository,
	deviceRepo repository.DeviceRepository,
	readingRepo repository.ReadingRepository,
	auditRepo repository.AuditRepository,
) AdminService {
	return &adminService{
		userRepo:    userRepo,
		deviceRepo:  deviceRepo,
		readingRepo: readingRepo,
		auditRepo:   auditRepo,
	}
}

func (s *adminService) GetStats(ctx context.Context) (*dto.AdminAggregateResponse, error) {
	totalUsers, err := s.userRepo.Count(ctx)
	if err != nil {
		return nil, err
	}

	totalDevices, err := s.deviceRepo.Count(ctx)
	if err != nil {
		return nil, err
	}

	activeDevices, err := s.deviceRepo.CountActive(ctx)
	if err != nil {
		return nil, err
	}

	totalReadings, err := s.readingRepo.Count(ctx)
	if err != nil {
		return nil, err
	}

	totalAudits, err := s.auditRepo.Count(ctx)
	if err != nil {
		return nil, err
	}

	return &dto.AdminAggregateResponse{
		TotalUsers:      int(totalUsers),
		TotalDevices:    int(totalDevices),
		ActiveDevices:   int(activeDevices),
		InactiveDevices: int(totalDevices - activeDevices),
		TotalReadings:   int(totalReadings),
		TotalAudits:     int(totalAudits),
	}, nil
}