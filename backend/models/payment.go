package models

import (
	"time"
	"gorm.io/gorm"
)

type Payment struct {
	gorm.Model
	SubscriptionID   uint      `gorm:"not null"`
	Amount          float64   `gorm:"type:decimal(10,2);not null"`
	Currency        string    `gorm:"type:varchar(3);default:'USD'"`
	Status          string    `gorm:"type:enum('pending','completed','failed','refunded');default:'pending'"`
	StripePaymentID string    `gorm:"type:varchar(255);not null"`
	InvoiceURL      string    `gorm:"type:varchar(512)"`
	PaidAt          *time.Time
	
	Subscription Subscription `gorm:"foreignKey:SubscriptionID"`
}