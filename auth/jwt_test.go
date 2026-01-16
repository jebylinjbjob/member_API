package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestValidateToken(t *testing.T) {
	t.Run("valid token", func(t *testing.T) {
		userID := int64(123)
		email := "test@example.com"

		token, err := GenerateToken(userID, email)
		if err != nil {
			t.Fatalf("failed to generate token: %v", err)
		}

		claims, err := ValidateToken(token)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if claims.UserID != userID {
			t.Errorf("expected userID %d, got %d", userID, claims.UserID)
		}

		if claims.Email != email {
			t.Errorf("expected email %s, got %s", email, claims.Email)
		}
	})

	t.Run("invalid token format", func(t *testing.T) {
		_, err := ValidateToken("invalid-token")
		if err == nil {
			t.Error("expected error for invalid token format")
		}
	})

	t.Run("expired token", func(t *testing.T) {
		claims := &Claims{
			UserID: 123,
			Email:  "test@example.com",
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
				Issuer:    "member-api",
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString(jwtSecret)

		_, err := ValidateToken(tokenString)
		if err == nil {
			t.Error("expected error for expired token")
		}
	})

	t.Run("token with wrong signature", func(t *testing.T) {
		claims := &Claims{
			UserID: 123,
			Email:  "test@example.com",
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				Issuer:    "member-api",
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString([]byte("wrong-secret"))

		_, err := ValidateToken(tokenString)
		if err == nil {
			t.Error("expected error for token with wrong signature")
		}
	})

	t.Run("empty token string", func(t *testing.T) {
		_, err := ValidateToken("")
		if err == nil {
			t.Error("expected error for empty token string")
		}
	})
}
