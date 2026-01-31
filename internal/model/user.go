package model

import "time"

type User struct {
	ID   uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `gorm:"column:name;not null" json:"name"`

	Username string `gorm:"column:username;unique;not null" json:"username"`
	Email    string `gorm:"column:email;unique;not null" json:"email"`
	Password string `gorm:"column:passsword;not null" json:"-"`
	Role     string `gorm:"column:role;not null" json:"role"`

	// Timestamps
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	// For Admin
	CreatedBy uint `gorm:"column:created_by" json:"created_by"`
	UpdatedBy uint `gorm:"column:updated_by" json:"updated_by"`
}

func (User) TableName() string {
	return "users"
}
