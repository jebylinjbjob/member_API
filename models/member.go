package models

import "github.com/google/uuid"

// Member represents a user stored in PostgreSQL and managed by GORM.
// Member represents a user account in the system with authentication and tenant association.
// It contains user credentials, identification information, and API access details.
type Member struct {
	Name         string    `gorm:"size:255" json:"name"`                // 使用者名稱
	Email        string    `gorm:"size:255;uniqueIndex" json:"email"`   // 電子郵件，唯一
	PasswordHash string    `gorm:"size:255" json:"-"`                   // 密碼哈希，不對外暴露
	TenantsID    uuid.UUID `gorm:"type:uuid" json:"tenants_id"`         // 所屬租戶ID
	APIKey       string    `gorm:"size:255;uniqueIndex" json:"api_key"` // API金鑰，唯一
	Base
}
