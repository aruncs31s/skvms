package repository

import (
	"context"
	"time"

	"github.com/aruncs31s/skvms/internal/model"
	"gorm.io/gorm"
)

type DeviceStateHistoryRepository interface {
	BeginTx() *gorm.DB
	Create(
		ctx context.Context,
		tx *gorm.DB,
		history *model.DeviceStateHistory,
	) error
	GetDeviceStateHistory(
		ctx context.Context,
		deviceID uint,
		states []uint,
		from time.Time,
		to time.Time,
	) ([]model.DeviceStateHistoryReport, error)
}
type deviceStateHistoryRepository struct {
	db *gorm.DB
}

func NewDeviceStateHistoryRepository(
	db *gorm.DB,
) DeviceStateHistoryRepository {
	return &deviceStateHistoryRepository{db: db}
}

func (r *deviceStateHistoryRepository) BeginTx() *gorm.DB {
	return r.db.Begin()
}

func (r *deviceStateHistoryRepository) Create(
	ctx context.Context,
	tx *gorm.DB,
	history *model.DeviceStateHistory,
) error {
	if tx == nil {
		tx = r.db
	}
	return tx.WithContext(ctx).Create(history).Error
}
func (r *deviceStateHistoryRepository) GetDeviceStateHistory(
	ctx context.Context,
	deviceID uint,
	states []uint,
	from time.Time,
	to time.Time,
) ([]model.DeviceStateHistoryReport, error) {
	var histories []model.DeviceStateHistoryReport
	query := r.db.WithContext(ctx).Table("device_state_history dsh")

	if deviceID != 0 {
		query = query.Where("device_id = ?", deviceID)
	}
	if len(states) > 0 {
		query = query.Where("state_id IN ?", states)
	}
	if !from.IsZero() {
		query = query.Where("created_at >= ?", from)
	}
	if !to.IsZero() {
		query = query.Where("created_at <= ?", to)
	}
	query = query.Joins(
		"JOIN device_state ds ON ds.id = dsh.state_id",
		"JOIN users u ON u.id = dsh.created_by",
	)
	query = query.Select(
		[]string{
			"dsh.caused_action AS action",
			"ds.name AS state",
			"TO_CHAR(dsh.created_at, 'YYYY-MM-DD HH24:MI:SS') AS changed_at",
			"u.username AS changed_by",
		},
	)
	err := query.Find(&histories).Error
	return histories, err
}
