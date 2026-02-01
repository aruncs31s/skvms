package service

import (
	"context"
	"time"

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
}

type readingService struct {
	repo repository.ReadingRepository
}

func NewReadingService(repo repository.ReadingRepository) ReadingService {
	return &readingService{repo: repo}
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
