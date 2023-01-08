package model

import "gorm.io/datatypes"

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"uniqueIndex; not null" json:"username"`
	Password  string         `gorm:"not null" json:"password"`
	IsAdmin   bool           `json:"is_admin"`
	CreatedAt datatypes.Date `json:"created_at"`
	UpdatedAt datatypes.Date `json:"updated_at"`
}

func (User) TableName() string { return "users" }
