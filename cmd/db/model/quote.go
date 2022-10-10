package model

import "gorm.io/datatypes"

type Quote struct {
	ID        uint           `gorm:"primary_key" json:"id"`
	Quote     string         `json:"quote"`
	CreatedAt datatypes.Date `json:"created_at"`
	UpdatedAt datatypes.Date `json:"updated_at"`
}

func (Quote) TableName() string { return "quotes" }
