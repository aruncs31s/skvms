package dto

import "github.com/aruncs31s/skvms/internal/model"

type ControlRequest struct {
	Action uint `json:"action"`
}

func (r *ControlRequest) Validate() error {
	actions := model.DeviceActionsMap

	if _, exists := actions[model.DeviceAction(r.Action)]; !exists {
		return NewValidationError(
			"action",
			"invalid action",
		)
	}
	// Add validation logic if needed
	return nil
}

func NewControlRequest(action uint, command string) *ControlRequest {
	return &ControlRequest{
		Action: action,
	}
}

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Field + ": " + e.Message
}

func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}

type DeviceControlResponse struct {
	Device string `json:"device"`
	State  string `json:"state"`
}
