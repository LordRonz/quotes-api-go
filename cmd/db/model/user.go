package model

type User struct {
	BaseModel
	Username string `gorm:"uniqueIndex; not null" json:"username"`
	Password string `gorm:"not null" json:"password"`
	IsAdmin  bool   `json:"is_admin"`
}

func (User) TableName() string { return "users" }
