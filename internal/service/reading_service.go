package service

import (
	"context"
	"time"

	"github.com/aruncs31s/skvms/internal/dto"
	"github.com/aruncs31s/skvms/internal/model"
	"github.com/aruncs31s/skvms/internal/repository"
)

type ReadingService interface {
	ListByDevice(ctx context.Context, deviceID uint, limit int) ([]model.Reading, *model.Reading, error)
	ListByDeviceAndDateRange(
		ctx context.Context,
		deviceID uint,
		startTime,
		endTime time.Time,
	) ([]model.Reading, error)
	ListByDeviceWithInterval(
		ctx context.Context,
		deviceID uint,
		startTime time.Time,
		endTime time.Time,
		interval time.Duration,
		count int,
	) ([]model.Reading, error)
	GetStats(ctx context.Context, deviceID uint, startTime, endTime time.Time) (map[string]interface{}, error)
	RecordEssentialReadings(
		ctx context.Context,
		deviceID uint,
		req *dto.EssentialReadingRequest,
	) (*model.Reading, error)
	// RecordReadings(
	// 	ctx context.Context,
	// 	deviceID uint,
	// 	req *dto.ReadingRequest,
	// ) (*model.Reading, error)
	GetReadingsOfConnectedDevice(
		ctx context.Context,
		parentDeviceID uint,
		childDeviceID uint,
		startTime time.Time,
		endTime time.Time,
	) ([]model.Reading, model.Reading, error)
}

type readingService struct {
	repo   repository.ReadingRepository
	device DeviceService
}

func NewReadingService(
	repo repository.ReadingRepository,
	device DeviceService,
) ReadingService {
	return &readingService{
		repo:   repo,
		device: device,
	}
}

func (s *readingService) ListByDevice(ctx context.Context, deviceID uint, limit int) ([]model.Reading, *model.Reading, error) {
	readings, err := s.repo.ListByDevice(ctx, deviceID, limit)
	if err != nil {
		return nil, nil, err
	}
	if len(readings) == 0 {
		return readings, nil, nil
	}
	latest := readings[len(readings)-1]
	return readings, &latest, nil
}

func (s *readingService) ListByDeviceAndDateRange(
	ctx context.Context,
	deviceID uint,
	startTime time.Time,
	endTime time.Time,
) ([]model.Reading, error) {
	return s.repo.ListByDeviceAndDateRange(
		ctx,
		deviceID,
		startTime,
		endTime)
}

func (s *readingService) ListByDeviceWithInterval(
	ctx context.Context,
	deviceID uint,
	startTime time.Time,
	endTime time.Time,
	interval time.Duration,
	count int,
) ([]model.Reading, error) {
	return s.repo.ListByDeviceWithInterval(
		ctx,
		deviceID,
		startTime,
		endTime,
		interval,
		count)
}

func (s *readingService) GetStats(ctx context.Context, deviceID uint, startTime, endTime time.Time) (map[string]interface{}, error) {
	return s.repo.GetStats(ctx, deviceID, startTime, endTime)
}

func (s *readingService) RecordEssentialReadings(
	ctx context.Context,
	deviceID uint,
	req *dto.EssentialReadingRequest,
) (*model.Reading, error) {
	reading := &model.Reading{
		DeviceID:  deviceID,
		Voltage:   req.Voltage,
		Current:   req.Current,
		CreatedAt: time.Now(),
	}
	return s.repo.Create(
		ctx,
		reading,
	)
}
func (s *readingService) GetReadingsOfConnectedDevice(
	ctx context.Context,
	parentDeviceID uint,
	childDeviceID uint,
	startTime time.Time,
	endTime time.Time,
) ([]model.Reading, model.Reading, error) {

	if isRelated, err := s.device.IsParent(
		ctx,
		parentDeviceID,
		childDeviceID,
	); err != nil || !isRelated {
		return nil, model.Reading{}, err
	}

	readings, latestReading, err := s.repo.GetReadingsOfConnectedDevice(
		ctx,
		childDeviceID,
		startTime,
		endTime,
	)
	if err != nil {
		return nil, model.Reading{}, err
	}
	return readings, latestReading, err
}
