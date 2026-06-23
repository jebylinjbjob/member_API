package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	defaultIssuer   = "member-api"
	defaultAudience = "member-api"
	defaultScope    = "read write"
	tokenTTL        = 24 * time.Hour
)

var (
	jwtSecret   []byte
	jwtIssuer   string
	jwtAudience string
)

func init() {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your-secret-key-change-in-production"
	}
	jwtSecret = []byte(secret)

	jwtIssuer = envOrDefault("JWT_ISSUER", defaultIssuer)
	jwtAudience = envOrDefault("JWT_AUDIENCE", defaultAudience)
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// Claims 定義 OAuth 2.0 JWT Access Token claims（RFC 9068）
type Claims struct {
	Scope string `json:"scope,omitempty"`
	jwt.RegisteredClaims
}

// SubjectID 從 sub claim 解析使用者 UUID
func (c *Claims) SubjectID() (uuid.UUID, error) {
	if c.Subject == "" {
		return uuid.Nil, errors.New("missing sub claim")
	}
	return uuid.Parse(c.Subject)
}

// GenerateToken 生成 OAuth 2.0 JWT access token
func GenerateToken(userID uuid.UUID, scope string) (string, error) {
	if userID == uuid.Nil {
		return "", errors.New("user id must be a valid uuid")
	}
	if scope == "" {
		scope = defaultScope
	}

	now := time.Now()
	claims := &Claims{
		Scope: scope,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    jwtIssuer,
			Subject:   userID.String(),
			Audience:  jwt.ClaimStrings{jwtAudience},
			ExpiresAt: jwt.NewNumericDate(now.Add(tokenTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ValidateToken 驗證 OAuth 2.0 JWT access token
func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
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
		jwt.WithAudience(jwtAudience),
	)
	if err != nil {
		return nil, err
	}

	if token == nil || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	if _, err := claims.SubjectID(); err != nil {
		return nil, err
	}

	return claims, nil
}
