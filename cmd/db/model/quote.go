package model

import (
	"encoding/json"
)

type Quote struct {
	BaseModel
	Quote  string   `gorm:"not null" json:"quote"`
	Author string   `json:"author"`
	Tags   []string `gorm:"type:text[]"`
}

func (q Quote) MarshalBinary() ([]byte, error) {
	return json.Marshal(q)
}

func (Quote) TableName() string { return "quotes" }
