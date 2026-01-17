package models

import "github.com/google/uuid"

type Member struct {
	Name         string    `gorm:"size:255" json:"name"`
	Email        string    `gorm:"size:255;uniqueIndex" json:"email"`
	PasswordHash string    `gorm:"size:255" json:"-"`
	TenantsID    uuid.UUID `gorm:"type:uuid;index" json:"tenants_id"`
	APIKey       string    `gorm:"size:255;uniqueIndex" json:"api_key"`
	Base
}
