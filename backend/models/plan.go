package models

import "gorm.io/gorm"

type Plan struct {
	gorm.Model
	CreatorID      uint    `gorm:"not null"`
	Name          string  `gorm:"type:varchar(100);not null"`
	Description   string  `gorm:"type:text"`
	BasePrice     float64 `gorm:"type:decimal(10,2);not null"`
	Price         float64 `gorm:"type:decimal(10,2);not null"`
	StripePriceID string  `gorm:"type:varchar(255);not null"`
	IsActive      bool    `gorm:"default:true"`
	
	Creator Creator `gorm:"foreignKey:CreatorID"`
}