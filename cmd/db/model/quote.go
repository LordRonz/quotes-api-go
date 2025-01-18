package model

import (
	"encoding/json"

	"github.com/lib/pq"
)

type Quote struct {
	BaseModel
	Quote  string   `gorm:"not null" json:"quote"`
	Author string   `json:"author"`
	Tags   pq.StringArray `gorm:"type:text[]" json:"tags"`
}

func (q Quote) MarshalBinary() ([]byte, error) {
	return json.Marshal(q)
}

func (Quote) TableName() string { return "quotes" }
