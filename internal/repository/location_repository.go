package repository

import (
	"context"

	"github.com/aruncs31s/skvms/internal/model"
	"gorm.io/gorm"
)

type LocationRepository interface {
	LocationWriter
	LocationReader
}
type LocationWriter interface {
	Create(
		ctx context.Context,
		location *model.Location,
	) error
	Update(
		ctx context.Context,
		location *model.Location,
	) error
	Delete(
		ctx context.Context,
		id uint,
	) error
}

type LocationReader interface {
	List(
		ctx context.Context,
	) ([]model.Location, error)
	GetByID(
		ctx context.Context,
		id uint,
	) (*model.Location, error)
	Search(
		ctx context.Context,
		query string,
	) ([]model.Location, error)
	GetByCode(
		ctx context.Context,
		code string,
	) (*model.Location, error)
	GetConnectedDevicesCount(
		ctx context.Context,
		locationID uint,
		hardwareTypes ...model.HardwareType,
	) (int, error)
	GetUserCount(
		ctx context.Context,
		locationID uint,
	) (int, error)
}

type locationRepository struct {
	db *gorm.DB
}

func NewLocationRepository(db *gorm.DB) LocationRepository {
	// Move this to seperate constructor
	return &locationRepository{
		db: db,
	}
}

func (r *locationRepository) List(ctx context.Context) ([]model.Location, error) {
	var locations []model.Location
	if err := r.
		db.
		WithContext(ctx).
		Find(&locations).
		Error; err != nil {
		return nil, err
	}

	return locations, nil
}

func (r *locationRepository) GetByID(ctx context.Context, id uint) (*model.Location, error) {
	var location model.Location
	if err := r.db.WithContext(ctx).First(&location, id).Error; err != nil {
		return nil, err
	}
	return &location, nil
}

func (r *locationRepository) GetByCode(ctx context.Context, code string) (*model.Location, error) {
	var location model.Location
	if err := r.db.WithContext(ctx).Where("location_code = ?", code).First(&location).Error; err != nil {
		return nil, err
	}
	return &location, nil
}

func (r *locationRepository) Create(ctx context.Context, location *model.Location) error {
	return r.db.WithContext(ctx).Create(location).Error
}

func (r *locationRepository) Update(ctx context.Context, location *model.Location) error {
	return r.db.WithContext(ctx).Save(location).Error
}

func (r *locationRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Location{}, id).Error
}

func (r *locationRepository) Search(
	ctx context.Context,
	query string,
) ([]model.Location, error) {
	var locations []model.Location
	if err := r.db.WithContext(ctx).
		Where(
			"name LIKE ? OR code LIKE ?",
			"%"+query+"%", "%"+query+"%",
		).
		Find(&locations).Error; err != nil {
		return nil, err
	}
	return locations, nil
}

func (r *locationRepository) GetConnectedDevicesCount(
	ctx context.Context,
	locationID uint,
	hardwareTypes ...model.HardwareType,
) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Table("device_assignment").
		Joins("JOIN devices d On device_assignment.device_id = d.id").
		Joins("JOIN device_types dt ON dt.id = d.device_type").
		Where("dt.hardware_type IN ?", hardwareTypes).
		Where("location_id = ?", locationID).
		Count(&count).Error
	return int(count), err
}

func (r *locationRepository) GetUserCount(ctx context.Context, locationID uint) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Table("users").
		Where("location_id = ?", locationID).
		Count(&count).Error
	return int(count), err
}
