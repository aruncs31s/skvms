package model

type User struct {
	ID       uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name     string `gorm:"column:name;not null" json:"name"`
	Username string `gorm:"column:username;unique;not null" json:"username"`
	Email    string `gorm:"column:email;unique;not null" json:"email"`
	Password string `gorm:"column:passsword;not null" json:"-"`
}

func (User) TableName() string {
	return "users"
}
