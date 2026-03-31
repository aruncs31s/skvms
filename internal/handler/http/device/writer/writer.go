package writer

import "github.com/aruncs31s/skvms/internal/service"

type DeviceWriter struct {
	deviceService service.DeviceService
	auditService  service.AuditService
}

func NewDeviceWriter(
	deviceService service.DeviceService,
	auditService service.AuditService,
) *DeviceWriter {
	return &DeviceWriter{
		deviceService: deviceService,
		auditService:  auditService,
	}
}
