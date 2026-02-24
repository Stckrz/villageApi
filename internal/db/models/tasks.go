package models

import "time"

type Task struct {
	ID          uint       `gorm:"primaryKey"`
	Name        string     `gorm:"not null"`
	Description string     `gorm:"not null"`
	BuildingId  uint       `gorm:"index;not null"`
	Building    Building   `gorm:"constraint:OnDelete:CASCADE; json:"-""`
	IsCompleted bool       `gorm:"not null;default:false"`
	CompletedAt *time.Time `gorm:"index"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
