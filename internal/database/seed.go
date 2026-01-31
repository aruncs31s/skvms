package database

import (
	"fmt"
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
	devTypes := []model.DeviceTypes{
		{
			Name:         "esp8266",
			HardwareType: 1,
			CreatedBy:    1,
		},
		{
			Name:         "esp32",
			HardwareType: 1,
			CreatedBy:    1,
		},
		{
			Name:         "volt-current-meter",
			HardwareType: 2,
			CreatedBy:    1,
		},
		{
			Name:         "smart-switch",
			HardwareType: 3,
			CreatedBy:    1,
		},
		{
			Name:         "sensor-node",
			HardwareType: 4,
			CreatedBy:    1,
		},
		{
			Name:         "temperature-sensor",
			HardwareType: 4,
			CreatedBy:    1,
		},
		{
			Name:         "humidity-sensor",
			HardwareType: 4,
			CreatedBy:    1,
		},
		{
			Name:         "motion-detector",
			HardwareType: 4,
			CreatedBy:    1,
		},
		{
			Name:         "relay-module",
			HardwareType: 3,
			CreatedBy:    1,
		},
		{
			Name:         "power-monitor",
			HardwareType: 2,
			CreatedBy:    1,
		},
		{
			Name:         "energy-meter",
			HardwareType: 2,
			CreatedBy:    1,
		},
	}

	for _, dt := range devTypes {
		if err := db.FirstOrCreate(
			&model.DeviceTypes{},
			model.DeviceTypes{Name: dt.Name},
			dt,
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

	// Get all device types
	var deviceTypes []model.DeviceTypes
	if err := db.Find(&deviceTypes).Error; err != nil {
		return err
	}

	// Get versions
	var versions []model.Version
	if err := db.Order("id").Find(&versions).Error; err != nil {
		return err
	}

	// Device data for each type
	deviceData := map[string]struct {
		Name string
		MAC  string
		IP   string
	}{
		"volt-current-meter": {"Main Panel Meter", "AA:BB:CC:DD:EE:FF", "192.168.1.100"},
		"smart-switch":       {"Living Room Switch", "BB:CC:DD:EE:FF:AA", "192.168.1.101"},
		"sensor-node":        {"Temperature Sensor", "CC:DD:EE:FF:AA:BB", "192.168.1.102"},
		"temperature-sensor": {"Office Temp Sensor", "DD:EE:FF:AA:BB:CC", "192.168.1.103"},
		"humidity-sensor":    {"Warehouse Humidity", "EE:FF:AA:BB:CC:DD", "192.168.1.104"},
		"motion-detector":    {"Entrance Motion", "FF:AA:BB:CC:DD:EE", "192.168.1.105"},
		"relay-module":       {"Pump Relay", "AA:BB:CC:DD:EE:11", "192.168.1.106"},
		"power-monitor":      {"Server Room Monitor", "BB:CC:DD:EE:FF:22", "192.168.1.107"},
		"energy-meter":       {"Building Meter", "CC:DD:EE:FF:AA:33", "192.168.1.108"},
	}

	ipCounter := 100
	for _, dt := range deviceTypes {
		data, exists := deviceData[dt.Name]
		if !exists {
			data = struct {
				Name string
				MAC  string
				IP   string
			}{
				Name: dt.Name + " Device",
				MAC:  fmt.Sprintf("00:11:22:33:44:%02X", ipCounter),
				IP:   fmt.Sprintf("192.168.1.%d", ipCounter),
			}
		}

		// Assign version based on type
		versionIndex := (int(dt.ID) - 1) % len(versions)
		version := versions[versionIndex]

		err := db.Transaction(func(tx *gorm.DB) error {
			device := model.Device{
				Name:      data.Name,
				Type:      dt.ID,
				VersionID: version.ID,
				State:     1, // Active
				CreatedBy: 1,
				UpdatedBy: 1,
			}
			if err := tx.Create(&device).Error; err != nil {
				return err
			}

			now := time.Now()
			if err := tx.Create(&model.DeviceDetails{
				DeviceID:        device.ID,
				IPAddress:       data.IP,
				MACAddress:      data.MAC,
				FirmwareVersion: version.Version,
				LastSeenAt:      &now,
			}).Error; err != nil {
				return err
			}

			if err := tx.Create(&model.DeviceAddress{
				DeviceID: device.ID,
				Address:  "Smart Kerala Facility",
				City:     "Kochi",
			}).Error; err != nil {
				return err
			}

			return nil
		})
		if err != nil {
			return err
		}
		ipCounter++
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
