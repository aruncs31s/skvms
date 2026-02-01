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
	if err := seedDeviceStates(db); err != nil {
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
	if err := seedDeviceStateHistory(db); err != nil {
		return err
	}
	if err := seedAuditLogs(db); err != nil {
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

/* ---------------- Device States ---------------- */

func seedDeviceStates(db *gorm.DB) error {
	deviceStates := []model.DeviceState{
		{
			ID:   1,
			Name: "Active",
		},
		{
			ID:   2,
			Name: "Inactive",
		},
		{
			ID:   3,
			Name: "Maintenance",
		},
		{
			ID:   4,
			Name: "Decommissioned",
		},
	}

	for _, ds := range deviceStates {
		var existing model.DeviceState
		err := db.Where("id = ?", ds.ID).First(&existing).Error
		if err == gorm.ErrRecordNotFound {
			// Create new record with explicit ID
			if err := db.Exec("INSERT INTO device_states (id, name, device_type_id, created_at) VALUES (?, ?, ?, ?)",
				ds.ID, ds.Name, 0, time.Now()).Error; err != nil {
				return err
			}
		} else if err != nil {
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
				Name:         data.Name,
				DeviceTypeID: dt.ID,
				VersionID:    version.ID,
				CurrentState: 1, // Active
				CreatedBy:    1,
				UpdatedBy:    1,
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
	daysInYear := 365
	readingsPerDay := 1000

	// Batch insert for better performance
	batchSize := 500
	readings := make([]model.Reading, 0, batchSize)

	for _, deviceID := range deviceIDs {
		for day := 0; day < daysInYear; day++ {
			for reading := 0; reading < readingsPerDay; reading++ {
				// Calculate timestamp: spread 1000 readings evenly across 24 hours
				minutesInDay := 24 * 60
				minuteOffset := (reading * minutesInDay) / readingsPerDay
				timestamp := now.Add(-time.Duration(daysInYear-day)*24*time.Hour + time.Duration(minuteOffset)*time.Minute)

				// Generate realistic voltage and current with some variation
				voltageBase := 220.0
				voltageVariation := 5.0 * (float64(reading%100) / 100.0) // 0-5V variation
				voltage := voltageBase + voltageVariation - 2.5          // ±2.5V around base

				currentBase := 2.5
				currentVariation := 0.5 * (float64(reading%50) / 50.0) // 0-0.5A variation
				current := currentBase + currentVariation - 0.25       // ±0.25A around base

				readings = append(readings, model.Reading{
					DeviceID:  deviceID,
					Voltage:   voltage,
					Current:   current,
					CreatedAt: timestamp,
				})

				// Insert batch when reaching batch size
				if len(readings) >= batchSize {
					if err := db.Create(&readings).Error; err != nil {
						return err
					}
					readings = readings[:0] // Clear slice but keep capacity
				}
			}
		}

		// Insert remaining readings for this device
		if len(readings) > 0 {
			if err := db.Create(&readings).Error; err != nil {
				return err
			}
			readings = readings[:0]
		}
	}

	return nil
}

/* ---------------- Device State History ---------------- */

func seedDeviceStateHistory(db *gorm.DB) error {
	var count int64
	if err := db.Model(&model.DeviceStateHistory{}).Count(&count).Error; err != nil {
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
	// Create some state history for first few devices
	for i, deviceID := range deviceIDs {
		if i >= 5 { // Only for first 5 devices
			break
		}

		// Each device has some state transitions
		histories := []model.DeviceStateHistory{
			{
				DeviceID:     deviceID,
				CausedAction: model.ActionCreate,
				StateID:      1, // Active
				CreatedBy:    1, // Admin user
				CreatedAt:    now.Add(-72 * time.Hour),
			},
			{
				DeviceID:     deviceID,
				CausedAction: model.ActionTurnOff,
				StateID:      2, // Inactive
				CreatedBy:    1, // Admin user
				CreatedAt:    now.Add(-48 * time.Hour),
			},
			{
				DeviceID:     deviceID,
				CausedAction: model.ActionTurnOn,
				StateID:      1, // Active
				CreatedBy:    1, // Admin user
				CreatedAt:    now.Add(-24 * time.Hour),
			},
		}

		// Add extra history for some devices
		if i%2 == 0 {
			histories = append(histories, []model.DeviceStateHistory{
				{
					DeviceID:     deviceID,
					CausedAction: model.ActionConfigure,
					StateID:      1, // Stay Active
					CreatedBy:    1, // Admin user
					CreatedAt:    now.Add(-12 * time.Hour),
				},
				{
					DeviceID:     deviceID,
					CausedAction: model.ActionTurnOff,
					StateID:      2, // Inactive
					CreatedBy:    1, // Admin user
					CreatedAt:    now.Add(-6 * time.Hour),
				},
				{
					DeviceID:     deviceID,
					CausedAction: model.ActionTurnOn,
					StateID:      1, // Active
					CreatedBy:    1, // Admin user
					CreatedAt:    now.Add(-2 * time.Hour),
				},
			}...)
		}

		for _, history := range histories {
			if err := db.Create(&history).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

/* ---------------- Audit Logs ---------------- */

func seedAuditLogs(db *gorm.DB) error {
	var count int64
	if err := db.Model(&model.AuditLog{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	var deviceIDs []uint
	if err := db.Model(&model.Device{}).Limit(5).Pluck("id", &deviceIDs).Error; err != nil {
		return err
	}

	now := time.Now()

	// Admin login logs
	loginLogs := []model.AuditLog{
		{
			UserID:    1,
			Username:  "admin",
			Action:    model.ActionCreate,
			Details:   "User logged in",
			IPAddress: "192.168.1.50",
			CreatedAt: now.Add(-72 * time.Hour),
		},
		{
			UserID:    1,
			Username:  "admin",
			Action:    model.ActionCreate,
			Details:   "User logged in",
			IPAddress: "192.168.1.50",
			CreatedAt: now.Add(-48 * time.Hour),
		},
		{
			UserID:    1,
			Username:  "admin",
			Action:    model.ActionCreate,
			Details:   "User logged in",
			IPAddress: "192.168.1.50",
			CreatedAt: now.Add(-24 * time.Hour),
		},
		{
			UserID:    1,
			Username:  "admin",
			Action:    model.ActionCreate,
			Details:   "User logged in",
			IPAddress: "192.168.1.50",
			CreatedAt: now.Add(-2 * time.Hour),
		},
	}

	for _, log := range loginLogs {
		if err := db.Create(&log).Error; err != nil {
			return err
		}
	}

	// Device action logs
	actions := []struct {
		action  model.DeviceAction
		details string
		hours   int
	}{
		{model.ActionCreate, "Device created", 72},
		{model.ActionTurnOff, "Device turned off", 48},
		{model.ActionTurnOn, "Device turned on", 24},
		{model.ActionConfigure, "Device configured", 12},
		{model.ActionUpdate, "Device settings updated", 6},
	}

	for i, deviceID := range deviceIDs {
		for j, action := range actions {
			// Skip some actions for variety
			if (i+j)%3 == 0 {
				continue
			}

			devID := deviceID
			auditLog := model.AuditLog{
				UserID:    1,
				Username:  "admin",
				Action:    action.action,
				Details:   action.details,
				IPAddress: "192.168.1.50",
				DeviceID:  &devID,
				CreatedAt: now.Add(-time.Duration(action.hours) * time.Hour),
			}

			if err := db.Create(&auditLog).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
