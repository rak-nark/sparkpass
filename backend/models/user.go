package models

import "gorm.io/gorm"

type User struct {
    gorm.Model
    Email        string `gorm:"unique;not null"`
    PasswordHash string `gorm:"column:password_hash;not null"` // Mapea explícitamente a la columna
    IsCreator    bool   `gorm:"default:false"`
    StripeCustomerID string `gorm:"column:stripe_customer_id"`
    AvatarURL    string `gorm:"column:avatar_url"`
    FCMToken     string `gorm:"column:fcm_token"`
    
    // Este campo SOLO para binding (no se guarda en DB)
    Password     string `gorm:"-" json:"password"` // El "-" lo excluye de la DB
}