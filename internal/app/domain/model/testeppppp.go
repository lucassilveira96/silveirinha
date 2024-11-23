package model

import "time"

type Testeppppp struct {
	ID uint `gorm:"primaryKey;autoIncrement" json:"id"`
	Ssssssss *string ` json:"ssssssss"`
	CreatedAt time.Time `gorm:"autoCreateTime;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;not null" json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at"`
}

func (Testeppppp) TableName() string {
	return "testeppppp"
}
