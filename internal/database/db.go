package database

import (
	"fmt"

	"github.com/aruncs31s/skvms/internal/config"
	"github.com/aruncs31s/skvms/internal/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func New(cfg config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(
		&model.User{},
		&model.Device{},
		&model.DeviceDetails{},
		&model.DeviceAddress{},
		&model.Reading{},
		&model.AuditLog{},
		&model.DeviceTypes{},
		&model.Version{},
		&model.Feature{},
		&model.VersionFeature{},
		&model.ConnectedDevice{},
		&model.DeviceState{},
		&model.DeviceStateHistory{},
	); err != nil {
		return nil, err
	}

	// Fix existing devices with invalid device_state after migration
	if err := db.Exec("UPDATE devices SET current_state = 1 WHERE current_state = 0").Error; err != nil {
		return nil, fmt.Errorf("failed to update device states: %w", err)
	}

	return db, nil
}
