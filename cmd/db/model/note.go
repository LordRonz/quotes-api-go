package model

import (
	"encoding/json"
)

type Note struct {
	BaseModel
	Note string `gorm:"type:text; not null"`
}

func (n Note) MarshalBinary() ([]byte, error) {
	return json.Marshal(n)
}

func (Note) TableName() string { return "quotes" }
