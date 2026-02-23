package models

import "time"

type Building struct {
	ID            uint               `gorm:"primaryKey"`
	Name          string             `gorm:"not null"`
	Description   string             `gorm:"not null"`
	Categories    []BuildingCategory `gorm:"foreignKey:BuildingID;constraint:OnDelete:CASCADE;"`
	Tasks         []Task             `gorm:"foreignKey:BuildingID;constraint:OnDelete:CASCADE;"`
	ThumbnailPath string             `gorm:"not null"`
	ImagePath     string             `gorm:"not null"`
	CreatedAt     time.Time          
	UpdatedAt     time.Time         
}

type BuildingCategory struct {
	ID         uint   `gorm:"primaryKey"`
	BuildingID uint   `gorm:"index;not null"`
	Text       string `gorm:"not null"`
	CreatedAt     time.Time          
	UpdatedAt     time.Time         
}
