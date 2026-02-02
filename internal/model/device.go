package model

import (
	"time"
)

const (
	ActionCreate DeviceAction = 1 + iota
	ActionUpdate
	ActionDelete
	ActionTurnOn
	ActionTurnOff
	ActionConfigure
)

func (da DeviceAction) Validate() bool {
	_, exists := DeviceActionsMap[da]
	return exists
}

var DeviceActionsMap map[DeviceAction]string = map[DeviceAction]string{
	ActionCreate:    "create",
	ActionUpdate:    "update",
	ActionDelete:    "delete",
	ActionTurnOn:    "turn_on",
	ActionTurnOff:   "turn_off",
	ActionConfigure: "configure",
}

// State transitions: current state -> allowed actions
var DeviceStateTransitions = map[uint][]DeviceAction{
	1: { // Active
		ActionTurnOff,
		ActionConfigure,
	},
	2: { // Inactive
		ActionTurnOn,
		ActionConfigure,
	},
}

// State transition map: current state -> action -> next state
var DeviceStateActionResult = map[uint]map[DeviceAction]uint{
	1: { // Active
		ActionTurnOff:   2, // Active -> Inactive
		ActionConfigure: 1, // Stay active
	},
	2: { // Inactive
		ActionTurnOn:    1, // Inactive -> Active
		ActionConfigure: 2, // Stay inactive
	},
}

// Device Can be a Sensor, Actuator, Gateway, etc.
// I use device for esp32 and sensors
type Device struct {
	ID   uint   `gorm:"column:id;primaryKey;autoIncrement"`
	Name string `gorm:"column:name"`
	// 1 -  , 2 - Sensor
	// FK to DeviceTypes.ID
	DeviceTypeID uint `gorm:"column:device_type"`
	// FK to Version.ID
	VersionID uint `gorm:"column:version_id"`

	// 1= Active, 0 = Inactive , 2 = Maintenance, 3 = Decommissioned
	// Also FK to DeviceState.ID
	CurrentState uint `gorm:"column:device_state"`

	Details          DeviceDetails     `gorm:"foreignKey:DeviceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Address          DeviceAddress     `gorm:"foreignKey:DeviceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	DeviceType       DeviceTypes       `gorm:"foreignKey:DeviceTypeID;references:ID"`
	Version          Version           `gorm:"foreignKey:VersionID;references:ID"`
	Readings         []Reading         `gorm:"foreignKey:DeviceID;references:ID"`
	ConnectedDevices []ConnectedDevice `gorm:"foreignKey:ParentID;references:ID"`

	CreatedBy uint `gorm:"column:created_by"`
	UpdatedBy uint `gorm:"column:updated_by"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (Device) TableName() string {
	return "devices"
}

// Possible States for a device
// Different types of devices can have different states
// Every state must be a cause of some action and all the actions must be defined in DeviceActionsMap
type DeviceState struct {
	ID           uint      `gorm:"column:id;primaryKey"`
	Name         string    `gorm:"column:name"`
	DeviceTypeID uint      `gorm:"column:device_type_id"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (DeviceState) TableName() string {
	return "device_states"
}

type DeviceStateHistory struct {
	ID           uint         `gorm:"column:id;primaryKey;autoIncrement"`
	DeviceID     uint         `gorm:"column:device_id"`
	CausedAction DeviceAction `gorm:"column:caused_action"`
	StateID      uint         `gorm:"column:state_id"`
	CreatedBy    uint         `gorm:"column:created_by"`
	CreatedAt    time.Time    `gorm:"column:created_at;autoCreateTime"`
	DeviceState  DeviceState  `gorm:"foreignKey:StateID;references:ID"`
	Device       Device       `gorm:"foreignKey:DeviceID;references:ID"`
	User         User         `gorm:"foreignKey:CreatedBy;references:ID;constraint:-"`
}

func (DeviceStateHistory) TableName() string {
	return "device_state_history"
}

type DeviceStateHistoryReport struct {
	ActionCaused DeviceAction `gorm:"column:action"`
	StateName    string       `gorm:"column:state"`
	ChangedAt    string       `gorm:"column:changed_at"`
	ChangedBy    string       `gorm:"column:changed_by"`
}
