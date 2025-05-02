package models

import "gorm.io/gorm"

type User struct {
    gorm.Model
    Email        string `gorm:"unique;not null"`
    PasswordHash string `gorm:"column:password_hash;not null"`
    IsCreator    bool   `gorm:"default:false"`
    StripeCustomerID string `gorm:"column:stripe_customer_id"`
    AvatarURL    string `gorm:"column:avatar_url"`
    FCMToken     string `gorm:"column:fcm_token"`
    Password     string `gorm:"-" json:"password"`
    
    // Añade esta relación inversa
    PremiumContents []PremiumContent `gorm:"foreignKey:CreatorID"`
}