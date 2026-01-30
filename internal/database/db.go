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

    if err := db.AutoMigrate(&model.User{}, &model.Device{}, &model.DeviceDetails{}, &model.DeviceAddress{}, &model.Reading{}); err != nil {
        return nil, err
    }

    return db, nil
}