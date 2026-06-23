package auth

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var testUserID = uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")

func setupTest(t *testing.T) {
	t.Helper()
	if err := os.Setenv("JWT_SECRET", "test-secret-key"); err != nil {
		t.Fatalf("Failed to set JWT_SECRET: %v", err)
	}
	jwtSecret = []byte("test-secret-key")
	jwtIssuer = defaultIssuer
	jwtAudience = defaultAudience
}

func createTestClaims(userID uuid.UUID, expiresAt, issuedAt time.Time) *Claims {
	return &Claims{
		Scope: defaultScope,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			Issuer:    jwtIssuer,
			Audience:  jwt.ClaimStrings{jwtAudience},
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			ID:        uuid.New().String(),
		},
	}
}

func TestGenerateToken(t *testing.T) {
	setupTest(t)

	tests := []struct {
		name      string
		userID    uuid.UUID
		scope     string
		wantErr   bool
		validator func(t *testing.T, token string, userID uuid.UUID)
	}{
		{
			name:    "成功生成標準 token",
			userID:  testUserID,
			scope:   "read write",
			wantErr: false,
			validator: func(t *testing.T, token string, userID uuid.UUID) {
				parts := strings.Split(token, ".")
				if len(parts) != 3 {
					t.Errorf("Token 應該有 3 個部分，實際: %d", len(parts))
				}

				claims, err := ValidateToken(token)
				if err != nil {
					t.Fatalf("生成的 token 無法驗證: %v", err)
				}

				subjectID, err := claims.SubjectID()
				if err != nil {
					t.Fatalf("SubjectID() failed: %v", err)
				}
				if subjectID != userID {
					t.Errorf("SubjectID = %v, want %v", subjectID, userID)
				}
				if claims.Scope != "read write" {
					t.Errorf("Scope = %v, want read write", claims.Scope)
				}
			},
		},
		{
			name:    "空 scope 使用預設值",
			userID:  uuid.New(),
			scope:   "",
			wantErr: false,
			validator: func(t *testing.T, token string, _ uuid.UUID) {
				claims, err := ValidateToken(token)
				if err != nil {
					t.Fatalf("ValidateToken() failed: %v", err)
				}
				if claims.Scope != defaultScope {
					t.Errorf("Scope = %v, want %v", claims.Scope, defaultScope)
				}
			},
		},
		{
			name:    "nil UUID 應失敗",
			userID:  uuid.Nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GenerateToken(tt.userID, tt.scope)

			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if token == "" {
					t.Error("GenerateToken() 返回空 token")
				}
				if tt.validator != nil {
					tt.validator(t, token, tt.userID)
				}
			}
		})
	}
}

func TestGenerateTokenExpiration(t *testing.T) {
	setupTest(t)

	token, err := GenerateToken(testUserID, defaultScope)
	if err != nil {
		t.Fatalf("GenerateToken() failed: %v", err)
	}

	claims, err := ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken() failed: %v", err)
	}

	expectedExpiry := time.Now().Add(tokenTTL)
	actualExpiry := claims.ExpiresAt.Time
	diff := actualExpiry.Sub(expectedExpiry)
	if diff < -5*time.Second || diff > 5*time.Second {
		t.Errorf("ExpiresAt = %v, want around %v (diff: %v)", actualExpiry, expectedExpiry, diff)
	}

	if claims.Issuer != defaultIssuer {
		t.Errorf("Issuer = %v, want %v", claims.Issuer, defaultIssuer)
	}

	if len(claims.Audience) != 1 || claims.Audience[0] != defaultAudience {
		t.Errorf("Audience = %v, want [%v]", claims.Audience, defaultAudience)
	}

	if claims.ID == "" {
		t.Error("jti should not be empty")
	}

	issuedAt := claims.IssuedAt.Time
	now := time.Now()
	diff = now.Sub(issuedAt)
	if diff < 0 || diff > 5*time.Second {
		t.Errorf("IssuedAt = %v, want around %v", issuedAt, now)
	}
}

//nolint:gosec // G101: 測試用假 token，非真實憑證
func TestValidateToken(t *testing.T) {
	setupTest(t)

	validToken, err := GenerateToken(testUserID, defaultScope)
	if err != nil {
		t.Fatalf("Failed to generate valid token: %v", err)
	}

	tests := []struct {
		name           string
		token          string
		wantErr        bool
		wantUserID     uuid.UUID
		validateClaims bool
	}{
		{
			name:           "驗證有效 token",
			token:          validToken,
			wantErr:        false,
			wantUserID:     testUserID,
			validateClaims: true,
		},
		{
			name:    "空字串 token",
			token:   "",
			wantErr: true,
		},
		{
			name:    "完全無效的字串",
			token:   "not-a-jwt-token",
			wantErr: true,
		},
		{
			name:    "只有一個部分",
			token:   "onlyonepart",
			wantErr: true,
		},
		{
			name:    "只有兩個部分",
			token:   "two.parts",
			wantErr: true,
		},
		{
			name:    "Base64 但不是有效的 JWT",
			token:   "eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.invalid",
			wantErr: true,
		},
		{
			name:    "格式正確但簽名無效",
			token:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.invalidsignature",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := ValidateToken(tt.token)

			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if claims != nil {
					t.Error("ValidateToken() 錯誤時應返回 nil claims")
				}
				return
			}

			if claims == nil {
				t.Fatal("ValidateToken() 成功時不應返回 nil claims")
			}

			if tt.validateClaims {
				subjectID, err := claims.SubjectID()
				if err != nil {
					t.Fatalf("SubjectID() failed: %v", err)
				}
				if subjectID != tt.wantUserID {
					t.Errorf("SubjectID = %v, want %v", subjectID, tt.wantUserID)
				}
			}
		})
	}
}

func TestValidateTokenExpired(t *testing.T) {
	setupTest(t)

	tests := []struct {
		name       string
		expiryTime time.Time
		wantErr    bool
		errMsg     string
	}{
		{
			name:       "過期 1 小時",
			expiryTime: time.Now().Add(-1 * time.Hour),
			wantErr:    true,
			errMsg:     "expired",
		},
		{
			name:       "還沒過期（未來 1 小時）",
			expiryTime: time.Now().Add(1 * time.Hour),
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims := createTestClaims(
				testUserID,
				tt.expiryTime,
				time.Now().Add(-1*time.Hour),
			)
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			tokenString, err := token.SignedString(jwtSecret)
			if err != nil {
				t.Fatalf("Failed to create token: %v", err)
			}

			_, err = ValidateToken(tokenString)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errMsg != "" {
				if err == nil || !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("ValidateToken() error = %v, should contain %q", err, tt.errMsg)
				}
			}
		})
	}
}

func TestValidateTokenSignature(t *testing.T) {
	setupTest(t)

	tests := []struct {
		name       string
		setupToken func() string
		wantErr    bool
	}{
		{
			name: "使用錯誤的 secret 簽名",
			setupToken: func() string {
				claims := createTestClaims(testUserID, time.Now().Add(tokenTTL), time.Now())
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				tokenString, _ := token.SignedString([]byte("wrong-secret-key"))
				return tokenString
			},
			wantErr: true,
		},
		{
			name: "使用正確的 secret",
			setupToken: func() string {
				token, _ := GenerateToken(testUserID, defaultScope)
				return token
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := ValidateToken(tt.setupToken())
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateToken() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && claims != nil {
				t.Error("驗證失敗時應返回 nil claims")
			}
		})
	}
}

func TestValidateTokenIssuerAudience(t *testing.T) {
	setupTest(t)

	t.Run("錯誤 issuer", func(t *testing.T) {
		claims := createTestClaims(testUserID, time.Now().Add(tokenTTL), time.Now())
		claims.Issuer = "wrong-issuer"
		tokenString, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtSecret)

		if _, err := ValidateToken(tokenString); err == nil {
			t.Error("ValidateToken() should reject wrong issuer")
		}
	})

	t.Run("錯誤 audience", func(t *testing.T) {
		claims := createTestClaims(testUserID, time.Now().Add(tokenTTL), time.Now())
		claims.Audience = jwt.ClaimStrings{"wrong-audience"}
		tokenString, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtSecret)

		if _, err := ValidateToken(tokenString); err == nil {
			t.Error("ValidateToken() should reject wrong audience")
		}
	})

	t.Run("無效 sub", func(t *testing.T) {
		claims := createTestClaims(testUserID, time.Now().Add(tokenTTL), time.Now())
		claims.Subject = "not-a-uuid"
		tokenString, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtSecret)

		if _, err := ValidateToken(tokenString); err == nil {
			t.Error("ValidateToken() should reject invalid sub")
		}
	})
}

func TestValidateTokenAlgorithmSubstitution(t *testing.T) {
	setupTest(t)

	tests := []struct {
		name       string
		setupToken func() string
		wantErr    bool
	}{
		{
			name: "none 算法",
			setupToken: func() string {
				claims := createTestClaims(testUserID, time.Now().Add(tokenTTL), time.Now())
				token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
				tokenString, _ := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
				return tokenString
			},
			wantErr: true,
		},
		{
			name: "HS384 算法",
			setupToken: func() string {
				claims := createTestClaims(testUserID, time.Now().Add(tokenTTL), time.Now())
				token := jwt.NewWithClaims(jwt.SigningMethodHS384, claims)
				tokenString, _ := token.SignedString(jwtSecret)
				return tokenString
			},
			wantErr: true,
		},
		{
			name: "HS512 算法",
			setupToken: func() string {
				claims := createTestClaims(testUserID, time.Now().Add(tokenTTL), time.Now())
				token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
				tokenString, _ := token.SignedString(jwtSecret)
				return tokenString
			},
			wantErr: true,
		},
		{
			name: "正常 HS256 token",
			setupToken: func() string {
				token, _ := GenerateToken(testUserID, defaultScope)
				return token
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := ValidateToken(tt.setupToken())
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateToken() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && claims != nil {
				t.Error("驗證失敗時應返回 nil claims")
			}
		})
	}
}

func TestTokenRoundTrip(t *testing.T) {
	setupTest(t)

	userID := uuid.New()
	token, err := GenerateToken(userID, "read")
	if err != nil {
		t.Fatalf("GenerateToken() failed: %v", err)
	}

	claims, err := ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken() failed: %v", err)
	}

	subjectID, err := claims.SubjectID()
	if err != nil {
		t.Fatalf("SubjectID() failed: %v", err)
	}
	if subjectID != userID {
		t.Errorf("SubjectID = %v, want %v", subjectID, userID)
	}
	if claims.Scope != "read" {
		t.Errorf("Scope = %v, want read", claims.Scope)
	}
}
