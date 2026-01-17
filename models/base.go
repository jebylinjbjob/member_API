package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Base 基礎模型結構，包含通用的審計欄位
type Base struct {
	ID                   uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`               // 主鍵ID
	Sort                 int        `json:"sort"`                                         // 排序欄位
	CreationTime         time.Time  `gorm:"autoCreateTime" json:"created_at"`             // 創建時間
	CreatorId            uuid.UUID  `json:"creator_id"`                                   // 創建者ID
	LastModificationTime *time.Time `gorm:"autoUpdateTime" json:"last_modification_time"` // 最後修改時間
	LastModifierId       uuid.UUID  `json:"last_modifier_id"`                             // 最後修改者ID
	IsDeleted            bool       `gorm:"default:false" json:"-"`                       // 是否刪除
	DeletedAt            *time.Time `gorm:"index" json:"-"`                               // 刪除時間
}

// BeforeCreate hook to auto-generate UUID
func (b *Base) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}
