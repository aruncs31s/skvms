package dto

type UserView struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

// Username and Email Should Be unique , but only Username And Passowrd Are required
type CreateUserRequest struct {
	Name     string `json:"name"`
	Username string `json:"username" binding:"required"`
	Email    string `json:"email"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role"`
}

type UpdateUserRequest struct {
	Name     string `json:"name,omitempty"`
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	Phone    string `json:"phone,omitempty"`
	Password string `json:"password,omitempty"`
	Role     string `json:"role,omitempty"`
}

type UserProfile struct {
	User     UserView       `json:"user"`
	Devices  []DeviceView   `json:"devices"`
	Activity []AuditLogView `json:"activity"`
}
