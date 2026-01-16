package auth

import (
	"testing"
)

func TestCheckPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		hash     string
		want     bool
	}{
		{
			name:     "valid password",
			password: "password123",
			hash:     "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy", // bcrypt hash of "password123"
			want:     true,
		},
		{
			name:     "invalid password",
			password: "wrongpassword",
			hash:     "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy", // bcrypt hash of "password123"
			want:     false,
		},
		{
			name:     "empty password",
			password: "",
			hash:     "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy",
			want:     false,
		},
		{
			name:     "invalid hash format",
			password: "password123",
			hash:     "invalid_hash",
			want:     false,
		},
		{
			name:     "empty hash",
			password: "password123",
			hash:     "",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CheckPassword(tt.password, tt.hash)
			if got != tt.want {
				t.Errorf("CheckPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}
