package models

type User struct {
	Model
	Username string `gorm:"column:username"`
	Password string `gorm:"column:password"`
}

func (u *User) TableName() string {
	return "user"
}
