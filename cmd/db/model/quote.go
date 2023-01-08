package model

import (
	"encoding/json"

	"gorm.io/datatypes"
)

type Quote struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Quote     string         `gorm:"not null" json:"quote"`
	Author    string         `json:"author"`
	CreatedAt datatypes.Date `json:"created_at"`
	UpdatedAt datatypes.Date `json:"updated_at"`
}

func (q Quote) MarshalBinary() ([]byte, error) {
	return json.Marshal(q)
}

func (Quote) TableName() string { return "quotes" }
