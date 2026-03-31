package http

import (
	"github.com/aruncs31s/skvms/internal/handler/http/device/control"
	"github.com/aruncs31s/skvms/internal/handler/http/device/reader"
	"github.com/aruncs31s/skvms/internal/handler/http/device/writer"
	"github.com/aruncs31s/skvms/internal/service"
	"github.com/gin-gonic/gin"
)

type DeviceHandler struct {
	reader     reader.DeviceReader
	writer     *writer.DeviceWriter
	controller *control.DeviceControlHandler
}

func NewDeviceHandler(
	deviceService service.DeviceService,
	auditService service.AuditService,
) *DeviceHandler {
	return &DeviceHandler{
		reader: reader.NewDeviceReader(
			deviceService,
		),
		writer: writer.NewDeviceWriter(
			deviceService,
			auditService,
		),
		controller: control.NewDeviceControlHandler(deviceService, auditService),
	}
}

func (h *DeviceHandler) GetMyDevices(c *gin.Context) {
	h.reader.GetMyDevices(c)
}

func (h *DeviceHandler) ListDevices(c *gin.Context) {
	h.reader.ListDevices(c)
}

func (h *DeviceHandler) ListRecentDevices(c *gin.Context) {
	h.reader.ListRecentDevices(c)
}

func (h *DeviceHandler) GetDevice(c *gin.Context) {
	h.reader.GetDevice(c)
}

func (h *DeviceHandler) ControlDevice(c *gin.Context) {
	h.controller.ControlDevice(c)
}

func (h *DeviceHandler) CreateDevice(c *gin.Context) {
	h.writer.CreateDevice(c)
}

func (h *DeviceHandler) UpdateDevice(c *gin.Context) {
	h.writer.UpdateDevice(c)
}

func (h *DeviceHandler) FullUpdateDevice(c *gin.Context) {
	h.writer.FullUpdateDevice(c)
}

func (h *DeviceHandler) DeleteDevice(c *gin.Context) {
	h.writer.DeleteDevice(c)
}

func (h *DeviceHandler) AddConnectedDevice(c *gin.Context) {
	h.writer.AddConnectedDevice(c)
}

func (h *DeviceHandler) GetConnectedDevices(c *gin.Context) {
	h.reader.GetConnectedDevices(c)
}

func (h *DeviceHandler) SearchDevices(c *gin.Context) {
	h.reader.SearchDevices(c)
}

func (h *DeviceHandler) SearchMicrocontollerDevices(c *gin.Context) {
	h.reader.SearchMicrocontollerDevices(c)
}

func (h *DeviceHandler) SearchSensorDevices(c *gin.Context) {
	h.reader.SearchSensorDevices(c)
}

func (h *DeviceHandler) ListAllSensors(c *gin.Context) {
	h.reader.ListAllSensors(c)
}

func (h *DeviceHandler) GetSensorDevice(c *gin.Context) {
	h.reader.GetSensorDevice(c)
}

func (h *DeviceHandler) CreateSensorDevice(c *gin.Context) {
	h.writer.CreateSensorDevice(c)
}

func (h *DeviceHandler) CreateConnectedDevice(c *gin.Context) {
	h.writer.AddConnectedDevice(c)
}

func (h *DeviceHandler) RemoveConnectedDevice(c *gin.Context) {
	h.writer.RemoveConnectedDevice(c)
}

func (h *DeviceHandler) CreateConnectedDeviceWithDetails(c *gin.Context) {
	h.writer.CreateConnectedDeviceWithDetails(c)
}

func (h *DeviceHandler) GetRecentlyCreatedDevices(c *gin.Context) {
	h.reader.GetRecentlyCreatedDevices(c)
}

func (h *DeviceHandler) GetTotalCount(c *gin.Context) {
	h.reader.GetTotalCount(c)
}

func (h *DeviceHandler) GetOfflineDevices(c *gin.Context) {
	h.reader.GetOfflineDevices(c)
}

func (h *DeviceHandler) ListMicrocontrollerDevices(c *gin.Context) {
	h.reader.ListMicrocontrollerDevices(c)
}

func (h *DeviceHandler) GetMicrocontrollerStats(c *gin.Context) {
	h.reader.GetMicrocontrollerStats(c)
}
