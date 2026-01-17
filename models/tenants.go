package models

import "github.com/google/uuid"

// Member represents a user stored in PostgreSQL and managed by GORM.
type Tenants struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"` // 租戶唯一標識符
	Name        string    `gorm:"size:255;not null" json:"name"`  // 租戶名稱
	description string    `gorm:"size:255" json:"description"`    // 租戶描述
	Base
}
