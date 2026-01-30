package database

import (
	"time"

	"github.com/aruncs31s/skvms/internal/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) error {
	if err := seedAdminUser(db); err != nil {
		return err
	}

	if err := seedDevices(db); err != nil {
		return err
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

	devices := []model.Device{
		{Name: "Main Panel Meter", Type: "volt-current-meter", CreatedBy: 1, UpdatedBy: 1},
		{Name: "Workshop Feeder", Type: "volt-current-meter", CreatedBy: 1, UpdatedBy: 1},
		{Name: "Solar Inverter Line", Type: "volt-current-meter", CreatedBy: 1, UpdatedBy: 1},
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
