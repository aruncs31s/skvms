package repository

import (
	"context"
	"strings"
	"time"

	"github.com/aruncs31s/skvms/internal/model"
	"github.com/aruncs31s/skvms/utils"
	"gorm.io/gorm"
)

type ReadingRepository interface {
	ListByDevice(
		ctx context.Context,
		deviceID uint,
		limit int,
	) ([]model.Reading, error)
	ListByDeviceAndDateRange(
		ctx context.Context,
		deviceID uint,
		startTime time.Time,
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
	GetStats(
		ctx context.Context,
		deviceID uint,
		startTime time.Time,
		endTime time.Time,
	) (map[string]interface{}, error)
	Count(ctx context.Context) (int64, error)
	ReadingWriter
	GetLastReading(
		ctx context.Context,
		deviceID uint,
	) (*model.Reading, error)

	GetReadingsOfConnectedDevice(
		ctx context.Context,
		childDeviceID uint,
		startTime time.Time,
		endTime time.Time,
	) ([]model.Reading, model.Reading, error)
	ListByDeviceProgressive(ctx context.Context, device uint) ([]model.AvgCurentVoltageReading, error)
}
type ReadingWriter interface {
	Create(
		ctx context.Context,
		reading *model.Reading,
	) (*model.Reading, error)
	// Methods for writing readings can be added here
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
		Order("created_at DESC").
		Limit(limit).
		Find(&readings).Error
	if err != nil {
		return nil, err
	}
	return readings, nil
}

func (r *readingRepository) ListByDeviceWithInterval(ctx context.Context, deviceID uint, startTime, endTime time.Time, interval time.Duration, count int) ([]model.Reading, error) {
	if count <= 0 {
		count = 50
	}
	if count > 1000 {
		count = 1000
	}

	var readings []model.Reading

	// Calculate time points for sampling
	timePoints := make([]time.Time, 0, count)
	currentTime := startTime

	for i := 0; i < count && currentTime.Before(endTime); i++ {
		timePoints = append(timePoints, currentTime)
		currentTime = currentTime.Add(interval)
	}

	// If we have time points, query for readings closest to each time point
	if len(timePoints) > 0 {
		// Build query to get readings closest to each time point
		query := r.db.WithContext(ctx).Where("device_id = ?", deviceID)

		// Add conditions for each time point - find readings within a small window around each point
		conditions := make([]string, 0, len(timePoints))
		args := make([]interface{}, 0, len(timePoints)*3)

		for _, tp := range timePoints {
			// Look for readings within 30 seconds of each time point
			window := 30 * time.Second
			conditions = append(conditions, "(created_at >= ? AND created_at <= ?)")
			args = append(args, tp.Add(-window), tp.Add(window))
		}

		query = query.Where(strings.Join(conditions, " OR "), args...)

		err := query.Order("created_at DESC").Find(&readings).Error
		if err != nil {
			return nil, err
		}

		// Group readings by time windows and pick the closest one to each time point
		result := make([]model.Reading, 0, len(timePoints))
		for _, tp := range timePoints {
			var closest *model.Reading
			minDiff := time.Hour // Large initial difference

			for _, reading := range readings {
				diff := reading.CreatedAt.Sub(tp)
				if diff < 0 {
					diff = -diff
				}
				if diff < minDiff {
					minDiff = diff
					closest = &reading
				}
			}

			if closest != nil {
				result = append(result, *closest)
			}
		}

		return result, nil
	}

	return readings, nil
}

func (r *readingRepository) ListByDeviceAndDateRange(ctx context.Context, deviceID uint, startTime, endTime time.Time) ([]model.Reading, error) {
	var readings []model.Reading
	err := r.db.WithContext(ctx).
		Where("device_id = ? AND created_at >= ? AND created_at <= ?", deviceID, startTime, endTime).
		Order("created_at DESC").
		Find(&readings).Error
	if err != nil {
		return nil, err
	}
	return readings, nil
}

func (r *readingRepository) GetStats(ctx context.Context, deviceID uint, startTime, endTime time.Time) (map[string]interface{}, error) {
	var stats struct {
		MaxVoltage     float64
		MinVoltage     float64
		MaxCurrent     float64
		MinCurrent     float64
		MaxVoltageTime time.Time
		MinVoltageTime time.Time
	}

	var maxVReading model.Reading
	err := r.db.WithContext(ctx).
		Where("device_id = ? AND created_at >= ? AND created_at <= ?", deviceID, startTime, endTime).
		Order("voltage DESC").
		First(&maxVReading).Error
	if err == nil {
		stats.MaxVoltage = maxVReading.Voltage
		stats.MaxVoltageTime = maxVReading.CreatedAt
	}

	var minVReading model.Reading
	err = r.db.WithContext(ctx).
		Where("device_id = ? AND created_at >= ? AND created_at <= ?", deviceID, startTime, endTime).
		Order("voltage ASC").
		First(&minVReading).Error
	if err == nil {
		stats.MinVoltage = minVReading.Voltage
		stats.MinVoltageTime = minVReading.CreatedAt
	}

	var maxCReading model.Reading
	err = r.db.WithContext(ctx).
		Where("device_id = ? AND created_at >= ? AND created_at <= ?", deviceID, startTime, endTime).
		Order("current DESC").
		First(&maxCReading).Error
	if err == nil {
		stats.MaxCurrent = maxCReading.Current
	}

	var minCReading model.Reading
	err = r.db.WithContext(ctx).
		Where("device_id = ? AND created_at >= ? AND created_at <= ?", deviceID, startTime, endTime).
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
func (r *readingRepository) Create(ctx context.Context, reading *model.Reading) (*model.Reading, error) {
	err := r.db.WithContext(ctx).Create(reading).Error
	if err != nil {
		return nil, err
	}
	return reading, nil
}

func (r *readingRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Reading{}).Count(&count).Error
	return count, err
}
func (r *readingRepository) GetLastReading(ctx context.Context, deviceID uint) (*model.Reading, error) {
	var reading model.Reading
	err := r.db.WithContext(ctx).
		Where("device_id = ?", deviceID).
		Order("created_at DESC").
		First(&reading).Error
	if err != nil {
		return nil, err
	}
	return &reading, nil
}

func (r *readingRepository) GetReadingsOfConnectedDevice(
	ctx context.Context,
	childDeviceID uint,
	startTime time.Time,
	endTime time.Time,
) ([]model.Reading, model.Reading, error) {

	lastReadings, err := r.GetLastReading(
		ctx,
		childDeviceID,
	)
	if err != nil {
		return nil, model.Reading{}, err
	}
	readings, err := r.ListByDeviceAndDateRange(
		ctx,
		childDeviceID,
		startTime,
		endTime,
	)
	if err != nil {
		return nil, model.Reading{}, err
	}

	return readings, *lastReadings, nil
}
func (r *readingRepository) ListByDeviceProgressive(ctx context.Context, device uint) ([]model.AvgCurentVoltageReading, error) {

	var readings []model.AvgCurentVoltageReading

	today := utils.GetBeginningOfDay()
	night := utils.GetEndOfDay()
	query := `
		SELECT
		d.id,
		r.created_at,
		r.voltage,
		r.current,
		AVG(r.voltage) OVER (
			ORDER BY r.created_at
			ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW
		) AS avg_voltage,
		AVG(r.current) OVER (
			ORDER BY r.created_at
			ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW
		) AS avg_current
	FROM readings r
	JOIN devices d ON r.device_id = d.id
	WHERE r.created_at >= ? AND r.created_at <= ? AND d.id = ?
	ORDER BY r.created_at DESC
	`
	err := r.db.WithContext(ctx).Raw(query, today, night, device).Scan(&readings).Error
	if err != nil {
		return nil, err
	}
	return readings, nil
}
