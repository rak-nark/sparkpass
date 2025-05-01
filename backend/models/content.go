package models

import "gorm.io/gorm"

type PremiumContent struct {
	gorm.Model
	CreatorID   uint    `gorm:"not null"`
	Title       string  `gorm:"not null"`
	Description string
	S3Key       string  `gorm:"not null"`
	IsLocked    bool    `gorm:"default:true"`
	Price       float64 `gorm:"type:decimal(10,2);default:0.00"`
	ContentType string  `gorm:"type:enum('video','image','podcast','document')"`
	Slug string `gorm:"not null;unique"`
}
