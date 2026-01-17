package models

import "github.com/google/uuid"

type Tenants struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name        string    `gorm:"size:255;not null" json:"name"`
	Description string    `gorm:"size:255" json:"description"`
	Base
}
