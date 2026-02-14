package repository

import (
	"context"

	"github.com/aruncs31s/skvms/internal/model"
	"gorm.io/gorm"
)

type MicrocontrollersRepository interface {
	ListMicrocontrollerDevices(
		ctx context.Context,
		limit,
		offset int,
	) ([]model.MicrocontrollerDeviceView, error)
	GetTotalMicrocontrollerCount(
		ctx context.Context,
	) (int64, error)
	GetMicrocontrollerStats(
		ctx context.Context,
	) (model.MicrocontrollerStatsView, error)
	LatestVersion(
		ctx context.Context,
	) (string, error)
}
type microcontrollersRepository struct {
	db *gorm.DB
}

func NewMicrocontrollersRepository(db *gorm.DB) MicrocontrollersRepository {
	return &microcontrollersRepository{
		db: db,
	}
}
func (r *microcontrollersRepository) ListMicrocontrollerDevices(
	ctx context.Context,
	limit,
	offset int,
) ([]model.MicrocontrollerDeviceView, error) {
	var devices []model.MicrocontrollerDeviceView
	sql := `SELECT  
		d.id,
		connected_to.parent_id AS parent_id,
		d.name,
		dt.name AS type,
		dd.ip_address,
		dd.mac_address,
		v.name as firmware_version,
		ds.name as current_state,
		cdevice.name as used_by
	FROM devices d
	JOIN device_states ds  ON d.current_state  = ds.id
	JOIN device_types dt 
		ON dt.id = d.device_type
	LEFT JOIN versions v
	  ON v.device_id = d.id
	 AND v.name = (
		 SELECT MAX(v2.name) as version_name
		 FROM versions v2
		 WHERE v2.device_id = d.id
	 )
	LEFT JOIN device_details dd 
		ON dd.device_id = d.id
	LEFT JOIN connected_devices cd 
		ON cd.parent_id = d.id 
	LEFT JOIN connected_devices connected_to 
		ON connected_to.child_id = d.id
	LEFT JOIN devices cdevice ON cdevice.id = connected_to.parent_id
	WHERE dt.hardware_type = ?
	ORDER BY d.created_at DESC
	LIMIT ? OFFSET ?`
	err := r.db.WithContext(ctx).Raw(sql, model.HardwareTypeMicroController, limit, offset).Scan(&devices).Error
	if err != nil {
		return nil, err
	}
	return devices, nil
}
func (r *microcontrollersRepository) GetTotalMicrocontrollerCount(
	ctx context.Context,
) (int64, error) {
	var count int64
	sql := `SELECT  
		COUNT(*) 
	FROM devices d
	JOIN device_types dt 
		ON dt.id = d.device_type
	WHERE dt.hardware_type = ?`
	err := r.db.WithContext(ctx).Raw(sql, int(model.HardwareTypeMicroController)).Scan(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
func (r *microcontrollersRepository) GetMicrocontrollerStats(
	ctx context.Context,
) (model.MicrocontrollerStatsView, error) {

	var stats model.MicrocontrollerStatsView
	sql := `SELECT  
		COUNT(*) AS total_devices,
		SUM(CASE WHEN ds.name = 'online' THEN 1 ELSE 0 END) AS online_devices,
		SUM(CASE WHEN ds.name = 'offline' THEN 1 ELSE 0 END) AS offline_devices
	FROM devices d
	JOIN device_states ds  ON d.current_state  = ds.id
	JOIN device_types dt 
		ON dt.id = d.device_type
	WHERE dt.hardware_type = ?`
	err := r.db.WithContext(ctx).Raw(sql, int(model.HardwareTypeMicroController)).Scan(&stats).Error
	if err != nil {
		return model.MicrocontrollerStatsView{}, err
	}
	return stats, nil
}
func (r *microcontrollersRepository) LatestVersion(
	ctx context.Context,
) (string, error) {
	var version string
	sql := `
	SELECT  
	v.name as version_name
	FROM versions v
	ORDER BY v.name DESC
	LIMIT 1
	`
	err := r.db.WithContext(ctx).Raw(sql).Scan(&version).Error
	if err != nil {
		return "", err
	}
	return version, nil
}
