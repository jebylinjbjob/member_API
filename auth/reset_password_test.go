package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestResetPasswordTokenRoundTrip(t *testing.T) {
	setupTest(t)

	userID := uuid.New()
	token, err := GenerateResetPasswordToken(userID)
	if err != nil {
		t.Fatalf("GenerateResetPasswordToken() failed: %v", err)
	}

	got, err := ValidateResetPasswordToken(token)
	if err != nil {
		t.Fatalf("ValidateResetPasswordToken() failed: %v", err)
	}
	if got != userID {
		t.Errorf("userID = %v, want %v", got, userID)
	}
}

func TestValidateResetPasswordTokenRejectsAccessToken(t *testing.T) {
	setupTest(t)

	accessToken, err := GenerateToken(uuid.New(), defaultScope)
	if err != nil {
		t.Fatalf("GenerateToken() failed: %v", err)
	}

	if _, err := ValidateResetPasswordToken(accessToken); err != ErrResetPasswordBadToken {
		t.Errorf("ValidateResetPasswordToken() error = %v, want %v", err, ErrResetPasswordBadToken)
	}
}

func TestValidateResetPasswordTokenExpired(t *testing.T) {
	setupTest(t)

	userID := uuid.New()
	claims := &ResetPasswordClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    jwtIssuer,
			Subject:   userID.String(),
			Audience:  jwt.ClaimStrings{resetPasswordAudience},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			ID:        uuid.New().String(),
		},
	}

	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtSecret)
	if err != nil {
		t.Fatalf("SignedString() failed: %v", err)
	}

	if _, err := ValidateResetPasswordToken(tokenString); err != ErrResetPasswordBadToken {
		t.Errorf("ValidateResetPasswordToken() error = %v, want %v", err, ErrResetPasswordBadToken)
	}
}
