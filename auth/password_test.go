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
			hash:     "$2a$10$FULstuepT7uhcOKCWrdlL.NSFcrNXOIP8bfgqWwQLZSamJ5S/daaK", // bcrypt hash of "password123"
			want:     true,
		},
		{
			name:     "invalid password",
			password: "wrongpassword",
			hash:     "$2a$10$FULstuepT7uhcOKCWrdlL.NSFcrNXOIP8bfgqWwQLZSamJ5S/daaK", // bcrypt hash of "password123"
			want:     false,
		},
		{
			name:     "empty password",
			password: "",
			hash:     "$2a$10$FULstuepT7uhcOKCWrdlL.NSFcrNXOIP8bfgqWwQLZSamJ5S/daaK",
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
