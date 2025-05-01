package models

import "gorm.io/gorm"

type Tag struct {
	gorm.Model
	Name string `gorm:"type:varchar(50);unique;not null"`
	Slug string `gorm:"type:varchar(50);unique;not null"`
}

type ContentTag struct {
	ContentID uint `gorm:"primaryKey"`
	TagID     uint `gorm:"primaryKey"`
}