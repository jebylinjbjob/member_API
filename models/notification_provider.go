package models

import "github.com/google/uuid"

type ProviderType string

const (
	ProviderSMTP     ProviderType = "SMTP"
	ProviderSMS      ProviderType = "SMS"
	ProviderFirebase ProviderType = "FIREBASE"
)

type NotificationProvider struct {
	TenantsID uuid.UUID    `gorm:"type:uuid;not null;index" json:"tenants_id"`
	Name      string       `gorm:"size:255;not null" json:"name"`
	Type      ProviderType `gorm:"size:50;not null;index" json:"type"`
	Config    string       `gorm:"type:jsonb" json:"config"`
	IsActive  bool         `gorm:"default:true" json:"is_active"`
	Base
}

type SMTPConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	From     string `json:"from"`
	UseTLS   bool   `json:"use_tls"`
}

type SMSConfig struct {
	Provider  string `json:"provider"`
	AccountID string `json:"account_id"`
	AuthToken string `json:"auth_token"`
	FromPhone string `json:"from_phone"`
}

type FirebaseConfig struct {
	ProjectID      string `json:"project_id"`
	CredentialJSON string `json:"credential_json"`
	ServerKey      string `json:"server_key"`
}
