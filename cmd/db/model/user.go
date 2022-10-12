package model

import "gorm.io/datatypes"

type User struct {
	ID        uint           `gorm:"primary_key" json:"id"`
	Username  string         `json:"quote"`
	Password  string         `json:"password"`
	IsAdmin   bool           `json:"is_admin"`
	CreatedAt datatypes.Date `json:"created_at"`
	UpdatedAt datatypes.Date `json:"updated_at"`
}

func (User) TableName() string { return "users" }
