package models

import "github.com/google/uuid"

type MemberNotificationPreference struct {
	MemberID   uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:idx_member_provider" json:"member_id"`
	ProviderID uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:idx_member_provider" json:"provider_id"`
	IsDefault  bool      `gorm:"default:false" json:"is_default"`
	IsActive   bool      `gorm:"default:true" json:"is_active"`
	Base
}
