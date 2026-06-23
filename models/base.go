package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Base 基礎模型結構，包含通用的審計欄位
type Base struct {
	ID                   uint       `gorm:"primaryKey"                                               json:"id"`
	UUID                 uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();uniqueIndex;not null" json:"uuid"`
	Sort                 int        `                                                                json:"sort"`
	CreationTime         time.Time  `gorm:"autoCreateTime"                                           json:"created_at"`
	CreatorId            uint       `                                                                json:"creator_id"`
	LastModificationTime *time.Time `gorm:"autoUpdateTime"                                           json:"last_modification_time"`
	LastModifierId       uint       `                                                                json:"last_modifier_id"`
	IsDeleted            bool       `gorm:"default:false"                                            json:"-"`
	DeletedAt            *time.Time `gorm:"index"                                                    json:"-"`
}

func (b *Base) BeforeCreate(_ *gorm.DB) error {
	if b.UUID == uuid.Nil {
		b.UUID = uuid.New()
	}
	return nil
}
