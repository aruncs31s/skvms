package model

import "time"

type User struct {
	ID   uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `gorm:"column:name" json:"name"`

	Username string `gorm:"column:username;unique" json:"username"`
	Email    string `gorm:"column:email" json:"email"`
	Password string `gorm:"column:passsword;not null" json:"-"`
	Role     string `gorm:"column:role;default:'user'" json:"role"`

	// Timestamps
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	// For Admin
	CreatedByID *uint    `gorm:"column:created_by" json:"created_by"`
	UpdatedBy   uint     `gorm:"column:updated_by" json:"updated_by"`
	Devices     []Device `gorm:"foreignKey:CreatedBy;references:ID"`
}

func (User) TableName() string {
	return "users"
}

type UserDetail struct {
	ID     uint   `gorm:"primaryKey;autoIncrement"`
	UserID uint   `gorm:"column:user_id;not null;unique"`
	Phone  string `gorm:"column:phone"`

	User User `gorm:"foreignKey:UserID;references:ID"`
}

func (UserDetail) TableName() string {
	return "user_details"
}
