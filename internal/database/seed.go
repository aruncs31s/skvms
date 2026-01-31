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

/* ---------------- Device Types ---------------- */

func seedDeviceTypes(db *gorm.DB) error {
	deviceTypes := []string{
		"esp8266",
		"esp32",
		"volt-current-meter",
		"smart-switch",
		"sensor-node",
		"temperature-sensor",
		"humidity-sensor",
		"motion-detector",
		"relay-module",
		"power-monitor",
		"energy-meter",
	}

	for _, name := range deviceTypes {
		if err := db.FirstOrCreate(
			&model.DeviceTypes{},
			model.DeviceTypes{Name: name},
		).Error; err != nil {
			return err
		}
	}
	return nil
}

/* ---------------- Admin User ---------------- */

func seedAdminUser(db *gorm.DB) error {
	var user model.User
	err := db.Where("username = ?", "admin").First(&user).Error
	if err == nil {
		return nil
	}
	if err != gorm.ErrRecordNotFound {
		return err
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

/* ---------------- Versions + Features ---------------- */

func seedVersions(db *gorm.DB) error {
	versions := []string{"1.0.0", "1.1.0", "2.0.0"}

	for _, v := range versions {
		if err := db.FirstOrCreate(
			&model.Version{},
			model.Version{Version: v},
		).Error; err != nil {
			return err
		}
	}

	var v2 model.Version
	if err := db.Where("version = ?", "2.0.0").First(&v2).Error; err != nil {
		return err
	}

	features := []model.Feature{
		{VersionID: v2.ID, FeatureName: "remote-control", Enabled: true},
		{VersionID: v2.ID, FeatureName: "energy-monitoring", Enabled: true},
		{VersionID: v2.ID, FeatureName: "alert-system", Enabled: false},
	}

	for _, f := range features {
		if err := db.FirstOrCreate(
			&model.Feature{},
			model.Feature{
				VersionID:   f.VersionID,
				FeatureName: f.FeatureName,
			},
		).Error; err != nil {
			return err
		}
	}

	return nil
}

/* ---------------- Devices ---------------- */

func seedDevices(db *gorm.DB) error {
	var count int64
	if err := db.Model(&model.Device{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	var deviceType model.DeviceTypes
	if err := db.Where("type_name = ?", "volt-current-meter").
		First(&deviceType).Error; err != nil {
		return err
	}

	var version model.Version
	if err := db.Where("version = ?", "1.0.0").
		First(&version).Error; err != nil {
		return err
	}

	devices := []struct {
		Name string
		MAC  string
		IP   string
	}{
		{"Main Panel Meter", "AA:BB:CC:DD:EE:FF", "192.168.1.100"},
		{"Workshop Feeder", "BB:CC:DD:EE:FF:AA", "192.168.1.101"},
		{"Solar Inverter Line", "CC:DD:EE:FF:AA:BB", "192.168.1.102"},
	}

	for _, d := range devices {
		err := db.Transaction(func(tx *gorm.DB) error {
			device := model.Device{
				Name:      d.Name,
				Type:      deviceType.ID,
				VersionID: version.ID,
				CreatedBy: 1,
				UpdatedBy: 1,
			}
			if err := tx.Create(&device).Error; err != nil {
				return err
			}

			now := time.Now()
			if err := tx.Create(&model.DeviceDetails{
				DeviceID:        device.ID,
				IPAddress:       d.IP,
				MACAddress:      d.MAC,
				FirmwareVersion: "1.0.0",
				LastSeenAt:      &now,
			}).Error; err != nil {
				return err
			}

			return tx.Create(&model.DeviceAddress{
				DeviceID: device.ID,
				Address:  "Smart Kerala Facility",
				City:     "Kochi",
			}).Error
		})
		if err != nil {
			return err
		}
	}

	return nil
}

/* ---------------- Readings ---------------- */

func seedReadings(db *gorm.DB) error {
	var count int64
	if err := db.Model(&model.Reading{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	var deviceIDs []uint
	if err := db.Model(&model.Device{}).Pluck("id", &deviceIDs).Error; err != nil {
		return err
	}

	now := time.Now()
	for _, deviceID := range deviceIDs {
		for i := 0; i < 60; i++ {
			reading := model.Reading{
				DeviceID:  deviceID,
				Voltage:   220 + float64((i%7)-3)*0.8,
				Current:   2.5 + float64((i%9)-4)*0.12,
				Timestamp: now.Add(-time.Duration(60-i) * time.Minute).Unix(),
			}
			if err := db.Create(&reading).Error; err != nil {
				return err
			}
		}
	}
	return nil
}
