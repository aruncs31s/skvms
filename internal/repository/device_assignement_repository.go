package repository

import (
	"context"

	"github.com/aruncs31s/skvms/internal/model"
	"gorm.io/gorm"
)

type DeviceAssignmentRepository interface {
	DeviceAssignmentWriter
}

type DeviceAssignmentReader interface {
}

type DeviceAssignmentWriter interface {
	Create(
		ctx context.Context,
		assignment *model.DeviceAssignment,
	) error
}
type deviceAssignmentRepository struct {
	db *gorm.DB
}

func NewDeviceAssignmentRepository(
	db *gorm.DB,
) DeviceAssignmentRepository {
	return &deviceAssignmentRepository{
		db: db,
	}
}
func (r *deviceAssignmentRepository) Create(
	ctx context.Context,
	assignment *model.DeviceAssignment,
) error {
	return r.db.WithContext(ctx).Create(&assignment).Error
}
