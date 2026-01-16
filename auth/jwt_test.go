package auth

import (
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestGenerateToken(t *testing.T) {
	// 設定測試用的 JWT secret
	os.Setenv("JWT_SECRET", "test-secret-key")
	jwtSecret = []byte("test-secret-key")

	tests := []struct {
		name    string
		userID  int64
		email   string
		wantErr bool
	}{
		{
			name:    "Valid token generation",
			userID:  1,
			email:   "test@example.com",
			wantErr: false,
		},
		{
			name:    "Valid token with different user",
			userID:  999,
			email:   "another@example.com",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GenerateToken(tt.userID, tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && token == "" {
				t.Error("GenerateToken() returned empty token")
			}
		})
	}
}

func TestValidateToken(t *testing.T) {
	// 設定測試用的 JWT secret
	os.Setenv("JWT_SECRET", "test-secret-key")
	jwtSecret = []byte("test-secret-key")

	userID := int64(1)
	email := "test@example.com"
	token, err := GenerateToken(userID, email)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	tests := []struct {
		name    string
		token   string
		wantErr bool
		checkID bool
	}{
		{
			name:    "Valid token",
			token:   token,
			wantErr: false,
			checkID: true,
		},
		{
			name:    "Invalid token",
			token:   "invalid.token.here",
			wantErr: true,
			checkID: false,
		},
		{
			name:    "Empty token",
			token:   "",
			wantErr: true,
			checkID: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := ValidateToken(tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if claims == nil {
					t.Error("ValidateToken() returned nil claims")
					return
				}
				if tt.checkID && claims.UserID != userID {
					t.Errorf("ValidateToken() UserID = %v, want %v", claims.UserID, userID)
				}
				if tt.checkID && claims.Email != email {
					t.Errorf("ValidateToken() Email = %v, want %v", claims.Email, email)
				}
			}
		})
	}
}

func TestValidateTokenExpired(t *testing.T) {
	// 設定測試用的 JWT secret
	os.Setenv("JWT_SECRET", "test-secret-key")
	jwtSecret = []byte("test-secret-key")

	// 建立一個過期的 token
	expirationTime := time.Now().Add(-1 * time.Hour) // 1 小時前過期
	claims := &Claims{
		UserID: 1,
		Email:  "test@example.com",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			Issuer:    "member-api",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		t.Fatalf("Failed to create expired token: %v", err)
	}

	_, err = ValidateToken(tokenString)
	if err == nil {
		t.Error("ValidateToken() should reject expired token")
	}
}

func TestValidateTokenWrongSecret(t *testing.T) {
	// 使用錯誤的 secret 生成 token
	wrongSecret := []byte("wrong-secret-key")
	claims := &Claims{
		UserID: 1,
		Email:  "test@example.com",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "member-api",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(wrongSecret)
	if err != nil {
		t.Fatalf("Failed to create token with wrong secret: %v", err)
	}

	// 設定正確的 secret
	os.Setenv("JWT_SECRET", "test-secret-key")
	jwtSecret = []byte("test-secret-key")

	// 嘗試驗證（應該失敗）
	_, err = ValidateToken(tokenString)
	if err == nil {
		t.Error("ValidateToken() should reject token signed with wrong secret")
	}
}

func TestTokenRoundTrip(t *testing.T) {
	// 設定測試用的 JWT secret
	os.Setenv("JWT_SECRET", "test-secret-key")
	jwtSecret = []byte("test-secret-key")

	userID := int64(42)
	email := "roundtrip@example.com"

	// 生成 token
	token, err := GenerateToken(userID, email)
	if err != nil {
		t.Fatalf("GenerateToken() failed: %v", err)
	}

	// 驗證 token
	claims, err := ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken() failed: %v", err)
	}

	// 檢查 claims
	if claims.UserID != userID {
		t.Errorf("UserID = %v, want %v", claims.UserID, userID)
	}
	if claims.Email != email {
		t.Errorf("Email = %v, want %v", claims.Email, email)
	}
	if claims.Issuer != "member-api" {
		t.Errorf("Issuer = %v, want %v", claims.Issuer, "member-api")
	}
}
