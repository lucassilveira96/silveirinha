package model

import "time"

type Teste struct {
	ID uint `gorm:"primaryKey;autoIncrement" json:"id"`
	Teste *string ` json:"teste"`
	CreatedAt time.Time `gorm:"autoCreateTime;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;not null" json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at"`
}

func (Teste) TableName() string {
	return "teste"
}
