package models

import "gorm.io/gorm"

type Creator struct {
	UserID          uint    `gorm:"primaryKey"`
	Bio             string `gorm:"type:text;not null" validate:"required,max=500"`
	StripeAccountID string `gorm:"type:varchar(255);not null"`
}

// AfterCreate hook para manejar lógica post-creación
func (c *Creator) AfterCreate(tx *gorm.DB) error {
	// Aquí podrías agregar lógica para iniciar el proceso de Stripe Connect
	return nil
}