package models

import (
	

	"gorm.io/gorm"
)

type PremiumContent struct {
    gorm.Model
    ID          uint      `json:"id"`  
    CreatorID   uint   `gorm:"not null" json:"creatorId"`
    Creator     User   `gorm:"foreignKey:CreatorID;references:ID" json:"-"`
    Title       string `gorm:"not null" json:"title"`
    Description string `json:"description"`
    S3Key       string `gorm:"not null" json:"s3Key"`
    IsLocked    bool   `gorm:"default:true" json:"isLocked"`
    Price       float64 `gorm:"type:decimal(10,2);default:0.00" json:"price"`
    ContentType string `gorm:"type:enum('video','image','podcast','document')" json:"contentType"`
    Slug        string `gorm:"not null;unique" json:"slug"`
}




