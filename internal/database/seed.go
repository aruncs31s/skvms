package database

import (
	"time"

	"github.com/aruncs31s/skvms/internal/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) error {
	if err := seedDeviceTypes(db); err != nil {
		return err
	}

	if err := seedAdminUser(db); err != nil {
		return err
	}

	if err := seedVersions(db); err != nil {
		return err
	}

	if err := seedDevices(db); err != nil {
		return err
	}

	if err := seedReadings(db); err != nil {
		return err
	}

	return nil
}

func seedDeviceTypes(db *gorm.DB) error {
	var count int64
	if err := db.Model(&model.DeviceTypes{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	deviceTypes := []model.DeviceTypes{
		{Name: "volt-current-meter"},
		{Name: "smart-switch"},
		{Name: "sensor-node"},
	}

	for _, dt := range deviceTypes {
		if err := db.Create(&dt).Error; err != nil {
			return err
		}
	}

	return nil
}

func seedAdminUser(db *gorm.DB) error {
	var count int64
	if err := db.Model(&model.User{}).Where("username = ?", "admin").Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	admin := model.User{
		Name:     "System Admin",
		Username: "admin",
		Email:    "admin@example.com",
		Password: string(hashed),
	}
	return db.Create(&admin).Error
}

func seedDevices(db *gorm.DB) error {
	var count int64
	if err := db.Model(&model.Device{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	// Get the volt-current-meter device type ID
	var deviceType model.DeviceTypes
	if err := db.Where("type_name = ?", "volt-current-meter").First(&deviceType).Error; err != nil {
		return err
	}
	var version model.Version
	if err := db.Where("version = ?", "1.0.0").First(&version).Error; err != nil {
		return err
	}
	devices := []model.Device{
		{Name: "Main Panel Meter", Type: deviceType.ID, VersionID: version.ID, CreatedBy: 1, UpdatedBy: 1},
		{Name: "Workshop Feeder", Type: deviceType.ID, VersionID: version.ID, CreatedBy: 1, UpdatedBy: 1},
		{Name: "Solar Inverter Line", Type: deviceType.ID, VersionID: version.ID, CreatedBy: 1, UpdatedBy: 1},
	}

	macs := []string{"AA:BB:CC:DD:EE:FF", "BB:CC:DD:EE:FF:AA", "CC:DD:EE:FF:AA:BB"}
	ips := []string{"192.168.1.100", "192.168.1.101", "192.168.1.102"}

	for i, device := range devices {
		if err := db.Create(&device).Error; err != nil {
			return err
		}

		lastSeen := time.Now()
		details := model.DeviceDetails{
			DeviceID:        device.ID,
			IPAddress:       ips[i],
			MACAddress:      macs[i],
			FirmwareVersion: "1.0.0",
			LastSeenAt:      &lastSeen,
		}
		if err := db.Create(&details).Error; err != nil {
			return err
		}

		address := model.DeviceAddress{
			DeviceID: device.ID,
			Address:  "Smart Kerala Facility",
			City:     "Kochi",
		}
		if err := db.Create(&address).Error; err != nil {
			return err
		}
	}

	return nil
}

func seedReadings(db *gorm.DB) error {
	var deviceIDs []uint
	if err := db.Model(&model.Device{}).Pluck("id", &deviceIDs).Error; err != nil {
		return err
	}
	if len(deviceIDs) == 0 {
		return nil
	}

	// If readings already exist, don't reseed.
	var count int64
	if err := db.Model(&model.Reading{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	now := time.Now()
	for _, deviceID := range deviceIDs {
		// 60 sample points, one per minute
		for i := 0; i < 60; i++ {
			ts := now.Add(-time.Duration(60-i) * time.Minute).Unix()
			voltage := 220.0 + float64((i%7)-3)*0.8
			current := 2.5 + float64((i%9)-4)*0.12

			reading := model.Reading{
				DeviceID:  deviceID,
				Voltage:   voltage,
				Current:   current,
				Timestamp: ts,
			}
			if err := db.Create(&reading).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

func seedVersions(db *gorm.DB) error {
	var count int64
	if err := db.Model(&model.Version{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	versions := []model.Version{
		{Version: "1.0.0"},
		{Version: "1.1.0"},
		{Version: "2.0.0"},
	}

	for _, v := range versions {
		if err := db.Create(&v).Error; err != nil {
			return err
		}
	}

	// Seed features for version 2.0.0 (assuming it gets ID 3)
	features := []model.Feature{
		{VersionID: 3, FeatureName: "remote-control", Enabled: true},
		{VersionID: 3, FeatureName: "energy-monitoring", Enabled: true},
		{VersionID: 3, FeatureName: "alert-system", Enabled: false},
	}

	for _, f := range features {
		if err := db.Create(&f).Error; err != nil {
			return err
		}
	}

	return nil
}
