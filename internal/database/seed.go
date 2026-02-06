package database

import (
	"fmt"
	"math"
	"math/rand"
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
	// if err := seedDeviceStateHistory(db); err != nil {
	// 	return err
	// }
	// if err := seedAuditLogs(db); err != nil {
	// 	return err
	// }
	return nil
}

/* ---------------- Device Types ---------------- */

func seedDeviceTypes(db *gorm.DB) error {
	devTypes := []model.DeviceTypes{
		{
			Name:         "ESP8266 (NODEMCU) ",
			HardwareType: model.HardwareTypeMicroController,
			CreatedBy:    1, //ADMIN
		},
		{
			Name:         "ESP32 (NODEMCU-32)",
			HardwareType: model.HardwareTypeMicroController,
			CreatedBy:    1,
		},
		{
			Name:         "Voltage , Current Sensor",
			HardwareType: model.HardwareTypeSensor,
			CreatedBy:    1,
		},

		{
			Name:         "Temperature Sensor",
			HardwareType: model.HardwareTypeSensor,
			CreatedBy:    1,
		},
		{
			Name:         "Humidity Sensor",
			HardwareType: model.HardwareTypeSensor,
			CreatedBy:    1,
		},
		{
			Name:         "Relay Module",
			HardwareType: model.HardwareTypeActuator,
			CreatedBy:    1,
		},
		{
			Name:         "Solar Charge Controller",
			HardwareType: model.HardwareTypeSolar,
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
		{
			ID:   5,
			Name: "Initialized",
		},
	}

	for _, ds := range deviceStates {
		var existing model.DeviceState
		err := db.Where("id = ?", ds.ID).First(&existing).Error
		if err == gorm.ErrRecordNotFound {
			// Create new record with explicit ID
			if err := db.Exec("INSERT INTO device_states (id, name, created_at) VALUES (?, ?, ?)",
				ds.ID, ds.Name, time.Now()).Error; err != nil {
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
	versions := []string{
		"1.0.0",
		"1.0.1",
		"2.0.0",
	}

	for _, v := range versions {
		if err := db.FirstOrCreate(
			&model.Version{},
			model.Version{Name: v},
		).Error; err != nil {
			return err
		}
	}

	var v2 model.Version
	if err := db.Where("name = ?", "2.0.0").First(&v2).Error; err != nil {
		return err
	}

	features := []model.Feature{
		{FeatureName: "remote-control", Enabled: true},
		{FeatureName: "energy-monitoring", Enabled: true},
		{FeatureName: "alert-system", Enabled: false},
	}

	for _, f := range features {
		if err := db.FirstOrCreate(
			&model.Feature{},
			model.Feature{FeatureName: f.FeatureName},
			&f,
		).Error; err != nil {
			return err
		}
	}

	// Associate features with version 2.0.0
	var existingFeatures []model.Feature
	if err := db.Where("feature_name IN ?", []string{"remote-control", "energy-monitoring", "alert-system"}).Find(&existingFeatures).Error; err != nil {
		return err
	}

	if err := db.Model(&v2).Association("Features").Append(existingFeatures); err != nil {
		return err
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

	// Create 100 test devices
	for i := 0; i < 100; i++ {
		// Cycle through device types
		dtIndex := i % len(deviceTypes)
		dt := deviceTypes[dtIndex]

		// Cycle through versions
		versionIndex := i % len(versions)
		version := versions[versionIndex]

		// Generate unique device data
		deviceName := fmt.Sprintf("Test Device %d", i+1)
		macAddress := fmt.Sprintf("AA:BB:CC:%02X:%02X:%02X", (i/65536)%256, (i/256)%256, i%256)
		ipAddress := fmt.Sprintf("192.168.1.%d", 100+i)

		err := db.Transaction(func(tx *gorm.DB) error {
			device := model.Device{
				Name:         deviceName,
				DeviceTypeID: dt.ID,
				VersionID:    &version.ID,
				CurrentState: 1, // Active
				CreatedBy:    1,
				UpdatedBy:    1,
			}
			if err := tx.Create(&device).Error; err != nil {
				return err
			}

			now := time.Now()
			if err := tx.Create(&model.DeviceDetails{
				DeviceID:   device.ID,
				IPAddress:  ipAddress,
				MACAddress: macAddress,
				LastSeenAt: &now,
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

	rand.Seed(time.Now().UnixNano())

	startDate := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	daysInYear := 365
	readingsPerDay := 100
	totalReadings := float64(daysInYear * readingsPerDay)

	// Batch insert for better performance
	batchSize := 500
	readings := make([]model.Reading, 0, batchSize)

	for _, deviceID := range deviceIDs {
		for day := 0; day < daysInYear; day++ {
			for reading := 0; reading < readingsPerDay; reading++ {
				readingIndex := day*readingsPerDay + reading

				// Calculate timestamp: spread 10 readings evenly across 24 hours
				minutesInDay := 24 * 60
				minuteOffset := (reading * minutesInDay) / readingsPerDay
				timestamp := startDate.AddDate(0, 0, day).Add(time.Duration(minuteOffset) * time.Minute)

				// Generate voltage: 10.8 to 12.5, increasing then decreasing (sine wave)
				angle := 2 * math.Pi * float64(readingIndex) / totalReadings
				voltageBase := 11.65
				voltageAmplitude := 0.85
				voltage := voltageBase + voltageAmplitude*math.Sin(angle)
				// Add small random variation per device (±0.1V)
				voltage += (rand.Float64() - 0.5) * 2

				// Generate current: 1 to 5 A
				current := 1.0 + rand.Float64()*4.0

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
