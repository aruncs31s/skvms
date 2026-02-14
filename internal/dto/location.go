package dto

type CreateLocationRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Code        string `json:"code" binding:"required"`
}

type UpdateLocationRequest struct {
	ID          uint   `json:"id"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Code        string `json:"code" binding:"required"`
}

type LocationResponse struct {
	ID                    uint   `json:"id"`
	Name                  string `json:"name"`
	Description           string `json:"description"`
	Code                  string `json:"code"`
	City                  string `json:"city"`
	State                 string `json:"state"`
	PinCode               string `json:"pin_code"`
	ConnectedDevicesCount int    `json:"device_count"`
	UserCount             int    `json:"user_count"`
}
