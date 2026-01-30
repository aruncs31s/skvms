package repository

import (
	"context"

	"github.com/aruncs31s/skvms/internal/model"
	"gorm.io/gorm"
)

type ReadingRepository interface {
	ListByDevice(ctx context.Context, deviceID uint, limit int) ([]model.Reading, error)
	ListByDeviceAndDateRange(ctx context.Context, deviceID uint, startTime, endTime int64) ([]model.Reading, error)
	GetStats(ctx context.Context, deviceID uint, startTime, endTime int64) (map[string]interface{}, error)
}

type readingRepository struct {
	db *gorm.DB
}

func NewReadingRepository(db *gorm.DB) ReadingRepository {
	return &readingRepository{db: db}
}

func (r *readingRepository) ListByDevice(ctx context.Context, deviceID uint, limit int) ([]model.Reading, error) {
	if limit <= 0 {
		limit = 50
	}
	if limit > 500 {
		limit = 500
	}

	var readings []model.Reading
	err := r.db.WithContext(ctx).
		Where("device_id = ?", deviceID).
		Order("timestamp ASC").
		Limit(limit).
		Find(&readings).Error
	if err != nil {
		return nil, err
	}
	return readings, nil
}

func (r *readingRepository) ListByDeviceAndDateRange(ctx context.Context, deviceID uint, startTime, endTime int64) ([]model.Reading, error) {
	var readings []model.Reading
	err := r.db.WithContext(ctx).
		Where("device_id = ? AND timestamp >= ? AND timestamp <= ?", deviceID, startTime, endTime).
		Order("timestamp ASC").
		Find(&readings).Error
	if err != nil {
		return nil, err
	}
	return readings, nil
}

func (r *readingRepository) GetStats(ctx context.Context, deviceID uint, startTime, endTime int64) (map[string]interface{}, error) {
	var stats struct {
		MaxVoltage     float64
		MinVoltage     float64
		MaxCurrent     float64
		MinCurrent     float64
		MaxVoltageTime int64
		MinVoltageTime int64
	}

	var maxVReading model.Reading
	err := r.db.WithContext(ctx).
		Where("device_id = ? AND timestamp >= ? AND timestamp <= ?", deviceID, startTime, endTime).
		Order("voltage DESC").
		First(&maxVReading).Error
	if err == nil {
		stats.MaxVoltage = maxVReading.Voltage
		stats.MaxVoltageTime = maxVReading.Timestamp
	}

	var minVReading model.Reading
	err = r.db.WithContext(ctx).
		Where("device_id = ? AND timestamp >= ? AND timestamp <= ?", deviceID, startTime, endTime).
		Order("voltage ASC").
		First(&minVReading).Error
	if err == nil {
		stats.MinVoltage = minVReading.Voltage
		stats.MinVoltageTime = minVReading.Timestamp
	}

	var maxCReading model.Reading
	err = r.db.WithContext(ctx).
		Where("device_id = ? AND timestamp >= ? AND timestamp <= ?", deviceID, startTime, endTime).
		Order("current DESC").
		First(&maxCReading).Error
	if err == nil {
		stats.MaxCurrent = maxCReading.Current
	}

	var minCReading model.Reading
	err = r.db.WithContext(ctx).
		Where("device_id = ? AND timestamp >= ? AND timestamp <= ?", deviceID, startTime, endTime).
		Order("current ASC").
		First(&minCReading).Error
	if err == nil {
		stats.MinCurrent = minCReading.Current
	}

	return map[string]interface{}{
		"max_voltage":      stats.MaxVoltage,
		"min_voltage":      stats.MinVoltage,
		"max_current":      stats.MaxCurrent,
		"min_current":      stats.MinCurrent,
		"max_voltage_time": stats.MaxVoltageTime,
		"min_voltage_time": stats.MinVoltageTime,
	}, nil
}
