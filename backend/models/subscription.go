package models

import (
	"time"
	"gorm.io/gorm"
)

type Subscription struct {
	gorm.Model
	UserID               uint      `gorm:"not null"`
	CreatorID            uint      `gorm:"not null"`
	PlanID               uint      `gorm:"not null"`
	StripeSubscriptionID string    `gorm:"type:varchar(255);not null"`
	Status               string    `gorm:"type:enum('active','canceled','paused','incomplete');default:'active'"`
	StartDate            time.Time `gorm:"autoCreateTime"`
	EndDate              *time.Time
	TrialEnd             *time.Time
	
	// Relaciones
	User    User    `gorm:"foreignKey:UserID"`
	Creator Creator `gorm:"foreignKey:CreatorID"`
	Plan    Plan    `gorm:"foreignKey:PlanID"`
}