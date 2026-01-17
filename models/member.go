package models

import "time"

// Member represents a user stored in PostgreSQL and managed by GORM.
type Member struct {
	Name                string     `gorm:"size:255;not null" json:"name"`
	Email               string     `gorm:"size:255;uniqueIndex;not null" json:"email"`
	PasswordHash        string     `gorm:"size:255" json:"-"`
	IsLocked            bool       `gorm:"default:false" json:"-"`
	FailedLoginAttempts int        `gorm:"default:0" json:"-"`
	LockedUntil         *time.Time `json:"-"`
	Base
}
