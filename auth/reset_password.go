package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	resetPasswordAudience = "member-api:reset-password"
	resetPasswordTTL      = 1 * time.Hour
)

// ErrResetPasswordBadToken 重設密碼 token 無效或已過期
var ErrResetPasswordBadToken = errors.New("RESET_PASSWORD_BAD_TOKEN")

// ResetPasswordClaims 重設密碼專用 JWT claims（對應 fastapi-users reset token）
type ResetPasswordClaims struct {
	jwt.RegisteredClaims
}

// GenerateResetPasswordToken 產生重設密碼用的一次性 token
func GenerateResetPasswordToken(userID uuid.UUID) (string, error) {
	if userID == uuid.Nil {
		return "", errors.New("user id must be a valid uuid")
	}

	now := time.Now()
	claims := &ResetPasswordClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    jwtIssuer,
			Subject:   userID.String(),
			Audience:  jwt.ClaimStrings{resetPasswordAudience},
			ExpiresAt: jwt.NewNumericDate(now.Add(resetPasswordTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ValidateResetPasswordToken 驗證重設密碼 token 並回傳使用者 UUID
func ValidateResetPasswordToken(tokenString string) (uuid.UUID, error) {
	claims := &ResetPasswordClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, jwt.ErrSignatureInvalid
			}
			return jwtSecret, nil
		},
		jwt.WithIssuer(jwtIssuer),
		jwt.WithAudience(resetPasswordAudience),
	)
	if err != nil {
		return uuid.Nil, ErrResetPasswordBadToken
	}

	if token == nil || !token.Valid {
		return uuid.Nil, ErrResetPasswordBadToken
	}

	if claims.Subject == "" {
		return uuid.Nil, ErrResetPasswordBadToken
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, ErrResetPasswordBadToken
	}

	return userID, nil
}
